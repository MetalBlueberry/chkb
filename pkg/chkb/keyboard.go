package chkb

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Keyboard struct {
	*Captor

	Layers    []*Layer
	LayerBook Book
	downkeys  map[KeyCode]MapEvent

	vkb uinput.Keyboard
}

func NewKeyboard(book Book, initialLayer string, vkb uinput.Keyboard) *Keyboard {
	kb := &Keyboard{
		LayerBook: book,
		Layers:    []*Layer{},
		vkb:       vkb,
		Captor:    NewCaptor(),
		downkeys:  map[KeyCode]MapEvent{},
	}
	kb.PushLayer(initialLayer)
	return kb
}

//go:generate stringer -type=Actions -trimprefix Action
type Actions int

const (
	ActionNil Actions = iota - 1
	ActionMap         //Default Action
	ActionDown
	ActionUp
	ActionTap
	ActionDoubleTap
	ActionHold
	ActionPushLayer
	ActionPopLayer
)

var ActionsMap map[string]Actions = map[string]Actions{
	ActionNil.String():       ActionNil,
	ActionMap.String():       ActionMap,
	ActionDown.String():      ActionDown,
	ActionUp.String():        ActionUp,
	ActionTap.String():       ActionTap,
	ActionDoubleTap.String(): ActionDoubleTap,
	ActionHold.String():      ActionHold,
	ActionPushLayer.String(): ActionPushLayer,
	ActionPopLayer.String():  ActionPopLayer,
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

func (ev *Actions) UnmarshalYAML(value *yaml.Node) error {
	var actionString string
	err := value.Decode(&actionString)
	if err != nil {
		return err
	}
	action, err := ParseAction(actionString)
	if err != nil {
		return err
	}
	*ev = action
	return nil
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

func (ev *KeyCode) UnmarshalYAML(value *yaml.Node) error {
	var codeString string
	err := value.Decode(&codeString)
	if err != nil {
		return err
	}
	code, err := ParseKeyCode(codeString)
	if err != nil {
		return err
	}
	*ev = code
	return nil
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

func (ev MapEvent) String() string {
	switch ev.Action {
	case ActionUp, ActionDown:
		return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
	case ActionPushLayer, ActionPopLayer:
		return fmt.Sprintf("%v-%s", ev.Action, ev.LayerName)
	default:
		return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
	}
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
	mapped := make([]MapEvent, 0)
	for i := range events {
		switch events[i].Action {
		case ActionUp:
			event, ok := kb.downkeys[events[i].KeyCode]
			if ok {
				event.Action = ActionUp
				mapped = append(mapped, event)
				delete(kb.downkeys, events[i].KeyCode)
				continue
			}
		}

		m, err := kb.Map(events[i])
		if err != nil {
			log.
				WithField("event", events[i]).
				WithError(err).
				Debug("Ignored event")
			continue
		}

		switch m.Action {
		case ActionDown:
			kb.downkeys[events[i].KeyCode] = m
		}

		mapped = append(mapped, m)
	}
	return mapped, nil
}

func (kb *Keyboard) Map(event KeyEvent) (MapEvent, error) {
	for i := len(kb.Layers) - 1; i >= 0; i-- {
		kmap, ok := kb.findMap(kb.Layers[i], event)
		if !ok {
			continue
		}
		log.
			WithField("from", event).
			WithField("to", kmap).
			Info("Map Key")
		return kmap, nil
	}
	switch event.Action {
	case ActionUp, ActionDown:
		return MapEvent{
			Action:  event.Action,
			KeyCode: event.KeyCode,
		}, nil
	default:
		return MapEvent{}, errors.New("Ignore event")
	}
}

func (kb *Keyboard) Deliver(events []MapEvent) error {
	for _, event := range events {
		switch event.Action {
		case ActionDown, ActionUp, ActionTap:
			err := kb.SendKeyEvent(event)
			if err != nil {
				return err
			}
		case ActionPushLayer:
			err := kb.PushLayer(event.LayerName)
			if err != nil {
				return err
			}
		case ActionPopLayer:
			err := kb.PopLayer()
			if err != nil {
				return err
			}
		default:
			log.
				WithField("Action", event.Action).
				Debug("Ignored event")
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
	case ActionTap:
		return kb.vkb.KeyPress(int(event.KeyCode))
	default:
		return errors.New("unknown event")
	}
}

func elapsed(from, to syscall.Timeval) time.Duration {
	return time.Unix(to.Sec, to.Usec*1000).Sub(time.Unix(from.Sec, from.Usec*1000))
}
