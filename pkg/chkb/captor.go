package chkb

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
	log "github.com/sirupsen/logrus"
)

type Captor struct {
	Events []evdev.InputEvent
}

func NewCaptor() *Captor {
	return &Captor{
		Events: []evdev.InputEvent{},
	}
}

func (c *Captor) Capture(events []evdev.InputEvent) ([]KeyEvent, error) {
	captured := make([]KeyEvent, 0)
	for i := range events {
		c, err := c.CaptureOne(events[i])
		if err != nil {
			log.
				WithField("event", evstring(events[i])).
				WithError(err).
				Debug("Skip event")
			continue
		}
		captured = append(captured, c...)
	}
	return captured, nil
}

func evstring(event evdev.InputEvent) string {
	evmap := map[uint16]string{
		0x01: "EV_KEY",
		0x02: "EV_REL",
		0x03: "EV_ABS",
		0x04: "EV_MSC",
		0x05: "EV_SW",
		0x11: "EV_LED",
		0x12: "EV_SND",
		0x14: "EV_REP",
		0x15: "EV_FF",
		0x16: "EV_PWR",
		0x17: "EV_FF_STATUS",
	}
	t, ok := evmap[event.Type]
	if !ok {
		t = strconv.Itoa(int(event.Type))
	}
	return fmt.Sprintf("Type: %s, Key: %s, Value: %d", t, evdev.KEY[int(event.Code)], event.Value)
}

func (c *Captor) CaptureOne(event evdev.InputEvent) ([]KeyEvent, error) {
	if event.Type != evdev.EV_KEY {
		return nil, errors.New("Only EV_KEY type supported")
	}
	captured := make([]KeyEvent, 0)

	switch evdev.KeyEventState(event.Value) {
	case evdev.KeyDown:
		captured = append(captured, NewKeyEv(event, KeyActionDown))
	case evdev.KeyUp:
		captured = append(captured, NewKeyEv(event, KeyActionUp))
	case evdev.KeyHold:
		lastEvent := c.Events[len(c.Events)-1]
		if lastEvent.Code == event.Code &&
			evdev.KeyEventState(lastEvent.Value) == evdev.KeyDown {

			log.Printf("Hold %s", evdev.KEY[int(event.Code)])
			captured = append(captured, NewKeyEv(event, KeyActionHold))
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
			captured = append(captured, NewKeyEv(event, KeyActionTap))
			break
		}

		if previous.Code == event.Code &&
			evdev.KeyEventState(event.Value) == evdev.KeyDown &&
			evdev.KeyEventState(previous.Value) == evdev.KeyDown &&
			el < time.Millisecond*200 {

			log.Printf("DoubleTap %s", evdev.KEY[int(event.Code)])
			captured = append(captured, NewKeyEv(event, KeyActionDoubleTap))
			break
		}
	}

	c.Events = append(c.Events, event)

	return captured, nil
}
