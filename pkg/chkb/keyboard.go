package chkb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
	"gopkg.in/yaml.v3"
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

//go:generate stringer -type=Actions -trimprefix Action
type Actions int

const (
	ActionMap Actions = iota
	ActionDown
	ActionUp
	ActionTap
	ActionDoubleTap
	ActionHold
	ActionPush
	ActionPop
)

var ActionsMap map[string]Actions = map[string]Actions{
	ActionMap.String():       ActionMap,
	ActionDown.String():      ActionDown,
	ActionUp.String():        ActionUp,
	ActionTap.String():       ActionTap,
	ActionDoubleTap.String(): ActionDoubleTap,
	ActionHold.String():      ActionHold,
	ActionPush.String():      ActionPush,
	ActionPop.String():       ActionPop,
}

func ParseAction(value string) (Actions, error) {
	if value == "" {
		return Actions(0), nil
	}
	a, ok := ActionsMap[value]
	if !ok {
		return a, fmt.Errorf("Action %s not found", value)
	}
	return a, nil
}

func (action Actions) MarshalYAML() (interface{}, error) {
	return action.String(), nil
}

func (action Actions) MarshalJSON() ([]byte, error) {
	return json.Marshal(action.String())
}

type KeyCode uint16

func ParseKeyCode(value string) (KeyCode, error) {
	if value == "" {
		return KeyCode(0), nil
	}
	code, ok := ecodes[value]
	if !ok {
		return KeyCode(code), fmt.Errorf("Code %s not found", value)
	}
	return KeyCode(code), nil
}

func (keyCode KeyCode) MarshalYAML() (interface{}, error) {
	return keyCode.String(), nil
}

func (keyCode KeyCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(keyCode.String())
}

func (keyCode KeyCode) String() string {
	return evdev.KEY[int(keyCode)]
}

type KeyEvent struct {
	Action  Actions
	KeyCode KeyCode
}

func (ev KeyEvent) String() string {
	return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
}

type MapEvent struct {
	Action  Actions `yaml:"action,omitempty"`
	KeyCode KeyCode `yaml:"keyCode,omitempty"`

	LayerName string `yaml:"layerName,omitempty"`
}

func (ev *MapEvent) UnmarshalYAML(value *yaml.Node) error {
	type mapEventString struct {
		Action    string `yaml:"action,omitempty"`
		KeyCode   string `yaml:"keyCode,omitempty"`
		LayerName string `yaml:"layerName,omitempty"`
	}
	tmp := &mapEventString{}
	value.Decode(tmp)
	var err error
	ev.Action, err = ParseAction(tmp.Action)
	if err != nil {
		return err
	}
	ev.KeyCode, err = ParseKeyCode(tmp.KeyCode)
	if err != nil {
		return err
	}
	ev.LayerName = tmp.LayerName
	return nil
}

func (ev MapEvent) String() string {
	return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
}

func NewKeyEv(event evdev.InputEvent, action Actions) KeyEvent {
	return KeyEvent{
		// Time:   time.Unix(event.Time.Sec, event.Time.Usec),
		KeyCode: KeyCode(event.Code),
		Action:  action,
	}
}

func (kb *Keyboard) findMap(layer *Layer, event KeyEvent) (kmap MapEvent, ok bool) {
	keymap, ok := layer.KeyMap[event.KeyCode]
	if !ok {
		return MapEvent{}, false
	}
	kmap, ok = keymap[event.Action]
	if ok {
		return kmap, true
	}

	if event.Action == ActionUp || event.Action == ActionDown {
		kmap, ok = keymap[ActionMap]
		if ok {
			kmap.Action = event.Action
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
			log.Printf("Ignored event %s", event.Action)
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
