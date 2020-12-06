package chkb

import (
	"container/list"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
)

type Captor struct {
	Events list.List

	TapTimeout chan InputEvent

	Clock clock.Clock
	// Keeps track of the downkeys pending jobs
	wg       sync.WaitGroup
	DownKeys map[KeyCode]TapTimer
}

func NewCaptor() *Captor {
	return NewCaptorWithClock(clock.New())
}

func NewCaptorWithClock(clock clock.Clock) *Captor {
	return &Captor{
		Events:   list.List{},
		Clock:    clock,
		DownKeys: make(map[KeyCode]TapTimer),
		wg:       sync.WaitGroup{},
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

func (c *Captor) PushHistory(event InputEvent) {
	c.Events.PushFront(event)
	if c.Events.Len() > 20 {
		c.Events.Remove(c.Events.Back())
	}
}

func (c *Captor) Run(capture func() ([]InputEvent, error), send func([]KeyEvent) error) error {
	for {
		events, err := capture()
		if err != nil {
			log.WithError(err).Error("cannot capture more events")
			c.wg.Wait()
			return err
		}

		captured, err := c.Capture(events)
		if err != nil {
			log.WithError(err).Error("Cannot capture events")
			return err
		}

		for _, event := range events {
			if event.Action == InputActionDown {
				timer := TapTimer{
					InputEvent: event,
					Timeout: c.Clock.AfterFunc(TapDelay, func() {
						defer c.wg.Done()
						delete(c.DownKeys, event.KeyCode)
						now := event.Time.Add(TapDelay)
						send([]KeyEvent{
							NewKeyEv(now, event, KeyActionHold),
						})
					}),
				}
				c.wg.Add(1)
				c.DownKeys[event.KeyCode] = timer
			}
			if event.Action == InputActionUp {
				if downKey, ok := c.DownKeys[event.KeyCode]; ok {
					delete(c.DownKeys, downKey.KeyCode)
					if downKey.Timeout.Stop() {
						c.wg.Done()
						now := c.Clock.Now()
						captured = append(captured, NewKeyEv(now, event, KeyActionTap))
					}
				}
			}
		}

		send(captured)
	}
}

func (c *Captor) Capture(events []InputEvent) ([]KeyEvent, error) {
	captured := make([]KeyEvent, 0)
	for i := range events {
		c, err := c.CaptureOne(events[i])
		if err != nil {
			log.
				WithField("event", events[i]).
				WithError(err).
				Debug("Skip event")
			continue
		}
		captured = append(captured, c...)
	}
	return captured, nil
}

func (c *Captor) CaptureOne(event InputEvent) ([]KeyEvent, error) {
	captured := make([]KeyEvent, 0)

	now := c.Clock.Now()
	switch event.Action {
	case InputActionDown:
		captured = append(captured, NewKeyEv(now, event, KeyActionDown))
	case InputActionUp:
		captured = append(captured, NewKeyEv(now, event, KeyActionUp))
	}

	return captured, nil
}

func NewKeyEv(time time.Time, event InputEvent, action KeyActions) KeyEvent {
	return KeyEvent{
		Time:    time,
		KeyCode: event.KeyCode,
		Action:  action,
	}
}
