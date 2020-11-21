package chkb

import (
	"errors"
	"log"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
)

type Captor struct {
	Events []evdev.InputEvent
}

func NewCaptor() *Captor {
	return &Captor{
		Events: []evdev.InputEvent{},
	}
}

func (c *Captor) Capture(events []*evdev.InputEvent) ([]KeyEv, error) {
	captured := make([]KeyEv, 0)
	for i := range events {
		c, err := c.CaptureOne(*events[i])
		if err != nil {
			log.Printf("Skip event %s", events[i])
			continue
		}
		captured = append(captured, c...)
	}
	return captured, nil
}

func (c *Captor) CaptureOne(event evdev.InputEvent) ([]KeyEv, error) {
	if event.Type != evdev.EV_KEY {
		return nil, errors.New("Only EV_KEY type supported")
	}
	captured := make([]KeyEv, 0)

	switch evdev.KeyEventState(event.Value) {
	case evdev.KeyDown:
		captured = append(captured, NewKeyEv(event, ActionDown))
	case evdev.KeyUp:
		captured = append(captured, NewKeyEv(event, ActionUp))
	case evdev.KeyHold:
		lastEvent := c.Events[len(c.Events)-1]
		if lastEvent.Code == event.Code &&
			evdev.KeyEventState(lastEvent.Value) == evdev.KeyDown {

			log.Printf("Hold %s", evdev.KEY[int(event.Code)])
			captured = append(captured, NewKeyEv(event, ActionHold))
		}
	}

	for i := len(c.Events) - 1; i >= 0; i-- {
		previous := c.Events[i]

		if event.Code != previous.Code {
			break
		}
		el := elapsed(previous.Time, event.Time)
		if previous.Code == event.Code &&
			evdev.KeyEventState(event.Value) == evdev.KeyUp &&
			evdev.KeyEventState(previous.Value) == evdev.KeyDown &&
			el < time.Millisecond*200 {

			log.Printf("Tap %s", evdev.KEY[int(event.Code)])
			captured = append(captured, NewKeyEv(event, ActionTap))
			break
		}

		if previous.Code == event.Code &&
			evdev.KeyEventState(event.Value) == evdev.KeyDown &&
			evdev.KeyEventState(previous.Value) == evdev.KeyDown &&
			el < time.Millisecond*200 {

			log.Printf("DoubleTap %s", evdev.KEY[int(event.Code)])
			captured = append(captured, NewKeyEv(event, ActionDoubleTap))
			break
		}
	}

	c.Events = append(c.Events, event)

	return captured, nil
}
