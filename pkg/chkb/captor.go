package chkb

import (
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
)

type Captor struct {
	Clock      clock.Clock
	IdleTimers map[KeyCode]*IdleTimer
	LastKey    *IdleTimer
	m          sync.Mutex
}

func NewCaptor() *Captor {
	return NewCaptorWithClock(clock.New())
}

func NewCaptorWithClock(clock clock.Clock) *Captor {
	return &Captor{
		Clock:      clock,
		IdleTimers: make(map[KeyCode]*IdleTimer),
		m:          sync.Mutex{},
	}
}

type IdleTimer struct {
	Timeout *clock.Timer
	Time    time.Time
	Count   int
	KeyCode KeyCode
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

func (c *Captor) loop(capture func() ([]InputEvent, error), send func([]KeyEvent) error) error {
	events, err := capture()
	if err != nil {
		log.WithError(err).Error("cannot capture more events")
		return err
	}
	c.m.Lock()
	defer c.m.Unlock()

	captured, err := c.Capture(events)
	if err != nil {
		log.WithError(err).Error("Cannot capture events")
		return err
	}

	for _, event := range captured {
		idleTimer, ok := c.IdleTimers[event.KeyCode]
		switch event.Action {
		case KeyActionDown:
			if !ok {
				idleTimer = &IdleTimer{
					Timeout: c.Clock.AfterFunc(TapDelay, c.idle(event.KeyCode, send)),
					KeyCode: event.KeyCode,
				}
				idleTimer.Timeout.Stop()
				c.IdleTimers[event.KeyCode] = idleTimer
			}
			if idleTimer.Timeout.Stop() {
				idleTimer.Count++
			} else {
				idleTimer.Count = 0
			}
			idleTimer.Timeout.Reset(TapDelay)
			idleTimer.Time = event.Time

			if c.LastKey != nil && c.LastKey != idleTimer {
				// Resolve event due to another keypress
				if c.LastKey.Timeout.Stop() {
					if c.LastKey.Count%2 == 0 {
						// Hold due to another keypress
						send([]KeyEvent{
							NewKeyEv(c.Clock.Now(), c.LastKey.KeyCode, KeyActionHold),
						})
					} else {
						// Tap due to another keypress
						switch c.LastKey.Count / 2 {
						case 0:
							send([]KeyEvent{
								NewKeyEv(c.Clock.Now(), c.LastKey.KeyCode, KeyActionTap),
							})
						case 1:
							send([]KeyEvent{
								NewKeyEv(c.Clock.Now(), c.LastKey.KeyCode, KeyActionDoubleTap),
							})
						default:
							send([]KeyEvent{
								NewKeyEv(c.Clock.Now(), c.LastKey.KeyCode, KeyActionNil),
							})
						}
					}
				}
			}
		case KeyActionUp:
			if ok {
				if idleTimer.Timeout.Stop() {
					idleTimer.Count++
					idleTimer.Timeout.Reset(TapDelay)
					idleTimer.Time = event.Time
				}
			}
		}

		if idleTimer != nil {
			c.LastKey = idleTimer
		}
	}

	return send(captured)
}

func (c *Captor) Run(capture func() ([]InputEvent, error), send func([]KeyEvent) error) error {
	for {
		err := c.loop(capture, send)
		if err != nil {
			return err
		}
	}
}

func (c *Captor) idle(keyCode KeyCode, send func([]KeyEvent) error) func() {
	return func() {
		c.m.Lock()
		defer c.m.Unlock()

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
