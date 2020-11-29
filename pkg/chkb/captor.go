package chkb

import (
	"container/list"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
	log "github.com/sirupsen/logrus"
)

type Captor struct {
	Events list.List
}

func NewCaptor() *Captor {
	return &Captor{
		Events: list.List{},
	}
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
	InputActionHold
)

func (c *Captor) PushHistory(event InputEvent) {
	c.Events.PushFront(event)
	if c.Events.Len() > 20 {
		c.Events.Remove(c.Events.Back())
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

	switch event.Action {
	case InputActionDown:
		captured = append(captured, NewKeyEv(event, KeyActionDown))
	case InputActionUp:
		captured = append(captured, NewKeyEv(event, KeyActionUp))
	case InputActionHold:
		lastEvent := c.Events.Front().Value.(InputEvent)
		if lastEvent.KeyCode == event.KeyCode &&
			lastEvent.Action == InputActionDown {

			log.Printf("Hold %s", evdev.KEY[int(event.KeyCode)])
			captured = append(captured, NewKeyEv(event, KeyActionHold))
		}
	}

	for i := c.Events.Front(); i != nil; i = i.Next() {

		previous := i.Value.(InputEvent)

		el := event.Time.Sub(previous.Time)
		if previous.KeyCode == event.KeyCode &&
			event.Action == InputActionUp &&
			previous.Action == InputActionDown &&
			el < time.Millisecond*200 {

			log.Printf("Tap %s", evdev.KEY[int(event.KeyCode)])
			captured = append(captured, NewKeyEv(event, KeyActionTap))
			break
		}

		if previous.KeyCode == event.KeyCode &&
			event.Action == InputActionDown &&
			previous.Action == InputActionDown &&
			el < time.Millisecond*200 {

			log.Printf("DoubleTap %s", evdev.KEY[int(event.KeyCode)])
			captured = append(captured, NewKeyEv(event, KeyActionDoubleTap))
			break
		}
	}

	c.PushHistory(event)

	return captured, nil
}

func NewKeyEv(event InputEvent, action KeyActions) KeyEvent {
	return KeyEvent{
		// Time:    time.Unix(event.Time.Sec, event.Time.Usec*1000),
		KeyCode: event.KeyCode,
		Action:  action,
	}
}
