package chkb

import (
	"errors"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

type Keyboard struct {
	// Layout   Layout
	Layers   []*Layer
	LayerIds map[uint16]*Layer

	// KeyMaps map[KeyEv]KeyEv

	vkb uinput.Keyboard

	Events []*evdev.InputEvent
}

func NewKeyboard(LayerIds map[uint16]*Layer, initialLayer uint16, vkb uinput.Keyboard) *Keyboard {
	kb := &Keyboard{
		LayerIds: LayerIds,
		Layers:   []*Layer{},
		Events:   []*evdev.InputEvent{},
		vkb:      vkb,
	}
	kb.PushLayer(initialLayer)
	return kb
}

func (kb *Keyboard) Capture(event *evdev.InputEvent) ([]KeyEv, error) {
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
		log.Printf("Hold %s", evdev.KEY[int(event.Code)])
		captured = append(captured, NewKeyEv(event, ActionHold))
	}

	for i := len(kb.Events) - 1; i >= 0; i-- {
		previous := kb.Events[i]

		if event.Code != previous.Code {
			break
		}

		if previous.Code == event.Code &&
			evdev.KeyEventState(event.Value) == evdev.KeyUp &&
			evdev.KeyEventState(previous.Value) == evdev.KeyDown &&
			elapsed(previous.Time, event.Time) < time.Millisecond*200 {

			log.Printf("Tap %s", evdev.KEY[int(event.Code)])
			captured = append(captured, NewKeyEv(event, ActionTap))
			break
		}

		if previous.Code == event.Code &&
			evdev.KeyEventState(event.Value) == evdev.KeyDown &&
			evdev.KeyEventState(previous.Value) == evdev.KeyDown &&
			elapsed(previous.Time, event.Time) < time.Millisecond*200 {

			log.Printf("DoubleTap %s", evdev.KEY[int(event.Code)])
			captured = append(captured, NewKeyEv(event, ActionDoubleTap))
			break
		}
	}

	kb.Events = append(kb.Events, event)

	return captured, nil
}

//go:generate stringer -type=Actions
type Actions int

const (
	ActionForward Actions = iota
	ActionDown
	ActionUp
	ActionTap
	ActionDoubleTap
	ActionHold
	ActionPush
	ActionPop
)

type KeyEv struct {
	Code   uint16
	Action Actions
}

func (ev KeyEv) String() string {
	return fmt.Sprintf("%s - %v", evdev.KEY[int(ev.Code)], ev.Action)
}

func NewKeyEv(event *evdev.InputEvent, action Actions) KeyEv {
	return KeyEv{
		// Time:   time.Unix(event.Time.Sec, event.Time.Usec),
		Code:   event.Code,
		Action: action,
	}
}

func (kb *Keyboard) findMap(layer *Layer, event KeyEv) (kmap KeyEv, ok bool) {
	kmap, ok = layer.KeyMap[event]
	if ok {
		return kmap, true
	}

	if event.Action == ActionUp || event.Action == ActionDown {
		originalAction := event.Action
		event.Action = ActionForward
		kmap, ok = layer.KeyMap[event]
		if ok {
			kmap.Action = originalAction
			return kmap, true
		}
	}
	return KeyEv{}, false
}

func (kb *Keyboard) Maps(events []KeyEv) ([]KeyEv, error) {
	mapped := make([]KeyEv, len(events))
	for i := range events {
		m, err := kb.Map(events[i])
		if err != nil {
			return nil, err
		}
		mapped[i] = m
	}
	return mapped, nil
}

func (kb *Keyboard) Map(event KeyEv) (KeyEv, error) {
	for i := len(kb.Layers) - 1; i >= 0; i-- {
		kmap, ok := kb.findMap(kb.Layers[i], event)
		if !ok {
			continue
		}
		log.Printf("Map key %s - %s", event, kmap)
		return kmap, nil
	}
	return event, nil
}

func (kb *Keyboard) Deliver(events []KeyEv) error {
	for _, event := range events {
		switch event.Action {
		case ActionDown, ActionUp:
			err := kb.SendKeyEvent(event)
			if err != nil {
				return err
			}
		case ActionPush:
			err := kb.PushLayer(event.Code)
			if err != nil {
				return err
			}
		case ActionPop:
			err := kb.PopLayer(event.Code)
			if err != nil {
				return err
			}
		default:
			log.Print("Ignored event")
		}
	}
	return nil
}

func (kb *Keyboard) PushLayer(id uint16) error {
	log.Printf("Push layer %d", id)
	l, ok := kb.LayerIds[id]
	if !ok {
		return errors.New("Layer do not exist")
	}
	kb.Layers = append(kb.Layers, l)
	return nil
}

func (kb *Keyboard) PopLayer(id uint16) error {
	log.Printf("Pop layer %d", id)
	if len(kb.Layers) == 1 {
		return errors.New("You cannot pop the last layer")
	}
	kb.Layers = kb.Layers[:len(kb.Layers)-1]
	return nil
}

func (kb *Keyboard) SendKeyEvent(event KeyEv) error {
	switch event.Action {
	case ActionDown:
		return kb.vkb.KeyDown(int(event.Code))
	case ActionUp:
		return kb.vkb.KeyUp(int(event.Code))
	// case ActionForward:
	// 	switch evdev.KeyEventState(event.value) {
	// 	case evdev.KeyDown:
	// 		return kb.vkb.KeyDown(int(event.Code))
	// 	case evdev.KeyUp:
	// 		return kb.vkb.KeyUp(int(event.Code))
	// 	default:
	// 		return errors.New("unknown event value")
	// 	}
	default:
		return errors.New("unknown event")
	}
}

type Event struct {
	Type string
}

type Layout struct {
	Keys map[uint16]*Key
}

type Key struct {
	Code   uint16
	Status int
}

type KeyMap struct {
	Type  string
	Value int
}

type Layer struct {
	KeyMap map[KeyEv]KeyEv
}

type KeyEvent struct {
	ID     uint32
	Status uint32
}

type Script struct {
	ID   uint32
	Code string
}

type LayerPush struct {
	ID     uint32
	Target string
}

type LayerPop struct {
	ID uint32
}

func elapsed(from, to syscall.Timeval) time.Duration {
	return time.Unix(from.Sec, from.Usec).Sub(time.Unix(to.Sec, to.Usec))
}
