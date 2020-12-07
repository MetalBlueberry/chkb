package chkb

import (
	"time"

	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
)

type Captor struct {
	Clock    clock.Clock
	DownKeys map[KeyCode]TapTimer
}

func NewCaptor() *Captor {
	return NewCaptorWithClock(clock.New())
}

func NewCaptorWithClock(clock clock.Clock) *Captor {
	return &Captor{
		Clock:    clock,
		DownKeys: make(map[KeyCode]TapTimer),
	}
}

type TapTimer struct {
	InputEvent
	Timeout *clock.Timer
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

		for _, event := range events {
			switch event.Action {
			case InputActionDown:
				c.DownKeys[event.KeyCode] = TapTimer{
					InputEvent: event,
					Timeout:    c.Clock.AfterFunc(TapDelay, c.idle(event.KeyCode, send)),
				}
			case InputActionUp:
				if downKey, ok := c.DownKeys[event.KeyCode]; ok {
					if downKey.Timeout.Stop() {
						delete(c.DownKeys, downKey.KeyCode)
						captured = append(captured, NewKeyEv(c.Clock.Now(), event.KeyCode, KeyActionTap))
					}
				}
			}
		}

		send(captured)
	}
}

func (c *Captor) idle(keyCode KeyCode, send func([]KeyEvent) error) func() {
	return func() {
		event := c.DownKeys[keyCode]
		delete(c.DownKeys, event.KeyCode)
		send([]KeyEvent{
			NewKeyEv(event.Time.Add(TapDelay), event.KeyCode, KeyActionHold),
		})
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
