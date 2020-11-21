package chkb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"syscall"
	"time"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

type Keyboard struct {
	*Captor

	Layers    []*Layer
	LayerBook Book

	vkb uinput.Keyboard
}

func NewKeyboard(book Book, initialLayer string, vkb uinput.Keyboard) *Keyboard {
	kb := &Keyboard{
		LayerBook: book,
		Layers:    []*Layer{},
		vkb:       vkb,
		Captor:    NewCaptor(),
	}
	kb.PushLayer(initialLayer)
	return kb
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

type KeyEventString string
type KeyEvent struct {
	Action  Actions
	KeyCode uint16
}

func (ev KeyEvent) Key() KeyEventString {
	return KeyEventString(ev.String())
}

func (ev KeyEvent) String() string {
	return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
}

type MapEvent struct {
	Action  Actions
	KeyCode uint16

	LayerName string
}

func (ev MapEvent) String() string {
	return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
}

func NewKeyEv(event evdev.InputEvent, action Actions) KeyEvent {
	return KeyEvent{
		// Time:   time.Unix(event.Time.Sec, event.Time.Usec),
		KeyCode: event.Code,
		Action:  action,
	}
}

func (kb *Keyboard) findMap(layer *Layer, event KeyEvent) (kmap MapEvent, ok bool) {
	kmap, ok = layer.KeyMap[event.Key()]
	if ok {
		return kmap, true
	}

	if event.Action == ActionUp || event.Action == ActionDown {
		originalAction := event.Action
		event.Action = ActionForward
		kmap, ok = layer.KeyMap[event.Key()]
		if ok {
			kmap.Action = originalAction
			return kmap, true
		}
	}
	return MapEvent{}, false
}

func (kb *Keyboard) Maps(events []KeyEvent) ([]MapEvent, error) {
	mapped := make([]MapEvent, len(events))
	for i := range events {
		m, err := kb.Map(events[i])
		if err != nil {
			return nil, err
		}
		mapped[i] = m
	}
	return mapped, nil
}

func (kb *Keyboard) Map(event KeyEvent) (MapEvent, error) {
	for i := len(kb.Layers) - 1; i >= 0; i-- {
		kmap, ok := kb.findMap(kb.Layers[i], event)
		if !ok {
			continue
		}
		log.Printf("Map key %s - %s", event, kmap)
		return kmap, nil
	}
	return MapEvent{
		Action:  event.Action,
		KeyCode: event.KeyCode,
	}, nil
}

func (kb *Keyboard) Deliver(events []MapEvent) error {
	for _, event := range events {
		switch event.Action {
		case ActionDown, ActionUp:
			err := kb.SendKeyEvent(event)
			if err != nil {
				return err
			}
		case ActionPush:
			err := kb.PushLayer(event.LayerName)
			if err != nil {
				return err
			}
		case ActionPop:
			err := kb.PopLayer()
			if err != nil {
				return err
			}
		default:
			log.Print("Ignored event")
		}
	}
	return nil
}

func (kb *Keyboard) PushLayer(name string) error {
	log.Printf("Push layer %s", name)
	l, ok := kb.LayerBook[name]
	if !ok {
		return errors.New("Layer do not exist")
	}
	kb.Layers = append(kb.Layers, l)
	return nil
}

func (kb *Keyboard) PopLayer() error {
	log.Printf("Pop layer")
	if len(kb.Layers) == 1 {
		return errors.New("You cannot pop the last layer")
	}
	kb.Layers = kb.Layers[:len(kb.Layers)-1]
	return nil
}

func (kb *Keyboard) SendKeyEvent(event MapEvent) error {
	switch event.Action {
	case ActionDown:
		return kb.vkb.KeyDown(int(event.KeyCode))
	case ActionUp:
		return kb.vkb.KeyUp(int(event.KeyCode))
	default:
		return errors.New("unknown event")
	}
}

func elapsed(from, to syscall.Timeval) time.Duration {
	return time.Unix(to.Sec, to.Usec*1000).Sub(time.Unix(from.Sec, from.Usec*1000))
}

type Layer struct {
	KeyMap map[KeyEventString]MapEvent
}

type Book map[string]*Layer

func (b *Book) Save(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(b)
}
