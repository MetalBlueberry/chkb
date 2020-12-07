package chkb

import (
	"time"

	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
)

type Captor struct {
	Clock      clock.Clock
	IdleTimers map[KeyCode]*IdleTimer
}

func NewCaptor() *Captor {
	return NewCaptorWithClock(clock.New())
}

func NewCaptorWithClock(clock clock.Clock) *Captor {
	return &Captor{
		Clock:      clock,
		IdleTimers: make(map[KeyCode]*IdleTimer),
	}
}

type IdleTimer struct {
	Timeout *clock.Timer
	Time    time.Time
	Count   int
}

type InputEvent struct {
	Time    time.Time
	KeyCode KeyCode
	Action  InputActions
}

//go:generate stringer -type=InputActions -trimprefix InputAction
type InputActions int

const (
	InputActionNil InputActions = iota
	InputActionDown
	InputActionUp
)

func (c *Captor) Run(capture func() ([]InputEvent, error), send func([]KeyEvent) error) error {
	for {
		events, err := capture()
		if err != nil {
			log.WithError(err).Error("cannot capture more events")
			return err
		}

		captured, err := c.Capture(events)
		if err != nil {
			log.WithError(err).Error("Cannot capture events")
			return err
		}

		for _, event := range captured {
			switch event.Action {
			case KeyActionDown:
				downKey, ok := c.IdleTimers[event.KeyCode]
				if !ok {
					downKey = &IdleTimer{
						Timeout: c.Clock.AfterFunc(TapDelay, c.idle(event.KeyCode, send)),
					}
					downKey.Timeout.Stop()
					c.IdleTimers[event.KeyCode] = downKey
				}
				if downKey.Timeout.Stop() {
					downKey.Count++
				} else {
					downKey.Count = 0
				}
				downKey.Timeout.Reset(TapDelay)
				downKey.Time = event.Time
			case KeyActionUp:
				if downKey, ok := c.IdleTimers[event.KeyCode]; ok {
					if downKey.Timeout.Stop() {
						downKey.Count++
						downKey.Timeout.Reset(TapDelay)
						downKey.Time = event.Time
					}
				}
			}
		}

		send(captured)
	}
}

func (c *Captor) idle(keyCode KeyCode, send func([]KeyEvent) error) func() {
	return func() {
		event := c.IdleTimers[keyCode]
		delete(c.IdleTimers, keyCode)
		if event.Count%2 == 0 {
			send([]KeyEvent{
				NewKeyEv(event.Time.Add(TapDelay), keyCode, KeyActionHold),
			})
		} else {
			switch event.Count / 2 {
			case 0:
				send([]KeyEvent{
					NewKeyEv(event.Time.Add(TapDelay), keyCode, KeyActionTap),
				})
			case 1:
				send([]KeyEvent{
					NewKeyEv(event.Time.Add(TapDelay), keyCode, KeyActionDoubleTap),
				})
			default:
				send([]KeyEvent{
					NewKeyEv(event.Time.Add(TapDelay), keyCode, KeyActionNil),
				})
			}
		}
	}
}

func (c *Captor) Capture(events []InputEvent) ([]KeyEvent, error) {
	captured := make([]KeyEvent, 0)
	for i := range events {
		switch events[i].Action {
		case InputActionDown:
			captured = append(captured, NewKeyEv(c.Clock.Now(), events[i].KeyCode, KeyActionDown))
		case InputActionUp:
			captured = append(captured, NewKeyEv(c.Clock.Now(), events[i].KeyCode, KeyActionUp))
		default:
			log.
				WithField("event", events[i]).
				Debug("Skip event")
			continue
		}
	}
	return captured, nil
}

func NewKeyEv(time time.Time, keyCode KeyCode, action KeyActions) KeyEvent {
	return KeyEvent{
		Time:    time,
		KeyCode: keyCode,
		Action:  action,
	}
}
