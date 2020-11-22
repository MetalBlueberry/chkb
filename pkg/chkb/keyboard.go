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

//go:generate stringer -type=KeyActions -trimprefix KeyAction
type KeyActions int

const (
	KeyActionNil KeyActions = iota - 1
	KeyActionMap            //Default Action
	KeyActionDown
	KeyActionUp
	KeyActionTap
	KeyActionDoubleTap
	KeyActionHold
)

var KeyActionsMap map[string]KeyActions = map[string]KeyActions{
	KeyActionNil.String():       KeyActionNil,
	KeyActionMap.String():       KeyActionMap,
	KeyActionDown.String():      KeyActionDown,
	KeyActionUp.String():        KeyActionUp,
	KeyActionTap.String():       KeyActionTap,
	KeyActionDoubleTap.String(): KeyActionDoubleTap,
	KeyActionHold.String():      KeyActionHold,
}

//go:generate stringer -type=KbActions -trimprefix KbAction
type KbActions int

const (
	KbActionNil KbActions = iota - 1
	KbActionMap           //Default Action
	KbActionDown
	KbActionUp
	KbActionTap
	KbActionDoubleTap
	KbActionHold
	KbActionPushLayer
	KbActionPopLayer
)

var KbActionsMap map[string]KbActions = map[string]KbActions{
	KeyActionNil.String():       KbActionNil,
	KeyActionMap.String():       KbActionMap,
	KeyActionDown.String():      KbActionDown,
	KeyActionUp.String():        KbActionUp,
	KeyActionTap.String():       KbActionTap,
	KeyActionDoubleTap.String(): KbActionDoubleTap,
	KeyActionHold.String():      KbActionHold,
	KbActionPushLayer.String():  KbActionPushLayer,
	KbActionPopLayer.String():   KbActionPopLayer,
}

func ParseKeyAction(value string) (KeyActions, error) {
	if value == "" {
		return KeyActions(0), nil
	}
	a, ok := KeyActionsMap[value]
	if !ok {
		return a, fmt.Errorf("Action %s not found", value)
	}
	return a, nil
}

func ParseKbAction(value string) (KbActions, error) {
	if value == "" {
		return KbActions(0), nil
	}
	a, ok := KbActionsMap[value]
	if !ok {
		return a, fmt.Errorf("Action %s not found", value)
	}
	return a, nil
}

func (ev *KbActions) UnmarshalYAML(value *yaml.Node) error {
	var actionString string
	err := value.Decode(&actionString)
	if err != nil {
		return err
	}
	action, err := ParseKbAction(actionString)
	if err != nil {
		return err
	}
	*ev = action
	return nil
}

func (ev *KeyActions) UnmarshalYAML(value *yaml.Node) error {
	var actionString string
	err := value.Decode(&actionString)
	if err != nil {
		return err
	}
	action, err := ParseKeyAction(actionString)
	if err != nil {
		return err
	}
	*ev = action
	return nil
}

func (action KeyActions) MarshalYAML() (interface{}, error) {
	return action.String(), nil
}

func (action KbActions) MarshalYAML() (interface{}, error) {
	return action.String(), nil
}

func (action KeyActions) MarshalJSON() ([]byte, error) {
	return json.Marshal(action.String())
}

func (action KbActions) MarshalJSON() ([]byte, error) {
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
	Action  KeyActions
	KeyCode KeyCode
}

func (ev KeyEvent) String() string {
	return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
}

type MapEvent struct {
	Action  KbActions `yaml:"action,omitempty"`
	KeyCode KeyCode   `yaml:"keyCode,omitempty"`

	LayerName string `yaml:"layerName,omitempty"`
}

func (ev MapEvent) String() string {
	switch ev.Action {
	case KbActionUp, KbActionDown:
		return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
	case KbActionPushLayer, KbActionPopLayer:
		return fmt.Sprintf("%v-%s", ev.Action, ev.LayerName)
	default:
		return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
	}
}

func NewKeyEv(event evdev.InputEvent, action KeyActions) KeyEvent {
	return KeyEvent{
		// Time:   time.Unix(event.Time.Sec, event.Time.Usec),
		KeyCode: KeyCode(event.Code),
		Action:  action,
	}
}

func (kb *Keyboard) findMap(layer *Layer, event KeyEvent) ([]MapEvent, bool) {
	keymap, ok := layer.KeyMap[event.KeyCode]
	if !ok {
		return nil, false
	}

	copymaps := make([]MapEvent, 0)
	//ActionMap
	if event.Action == KeyActionUp || event.Action == KeyActionDown {
		kmaps, ok := keymap[KeyActionMap]
		if ok {
			for i := range kmaps {
				m := kmaps[i]
				if m.Action == KbActionMap {
					m.Action = KbActions(event.Action)
				}
				copymaps = append(copymaps, m)
			}
		}
	}

	kmaps, ok := keymap[event.Action]
	if ok {
		copymaps = append(copymaps, kmaps...)
		return copymaps, true
	}
	return copymaps, len(copymaps) > 0
}

func (kb *Keyboard) Maps(events []KeyEvent) ([]MapEvent, error) {
	mapped := make([]MapEvent, 0)
	for _, event := range events {
		switch event.Action {
		case KeyActionUp:
			downkey, ok := kb.downkeys[event.KeyCode]
			if ok {
				downkey.Action = KbActionUp
				mapped = append(mapped, downkey)
				delete(kb.downkeys, event.KeyCode)
			}
		}

		maps, err := kb.Map(event)
		if err != nil {
			log.
				WithField("event", event).
				WithError(err).
				Debug("Ignored event")
			continue
		}

		for _, m := range maps {
			switch m.Action {
			case KbActionDown:
				kb.downkeys[event.KeyCode] = m
				mapped = append(mapped, m)
			case KbActionUp:
				_, ok := kb.downkeys[m.KeyCode]
				if ok {
					mapped = append(mapped, m)
					delete(kb.downkeys, m.KeyCode)
				}
			default:
				mapped = append(mapped, m)
			}
		}

	}
	return mapped, nil
}

func (kb *Keyboard) Map(event KeyEvent) ([]MapEvent, error) {
	for i := len(kb.Layers) - 1; i >= 0; i-- {
		kmaps, ok := kb.findMap(kb.Layers[i], event)
		if !ok {
			continue
		}
		log.
			WithField("from", event).
			WithField("to", kmaps).
			Info("Map Key")
		return kmaps, nil
	}
	// No map detected, forward
	switch event.Action {
	case KeyActionUp, KeyActionDown:
		return []MapEvent{
			{
				Action:  KbActions(event.Action),
				KeyCode: event.KeyCode,
			},
		}, nil
	default:
		return nil, errors.New("Ignore event")
	}
}

func (kb *Keyboard) Deliver(events []MapEvent) error {
	for _, event := range events {
		switch event.Action {
		case KbActionDown, KbActionUp, KbActionTap:
			err := kb.SendKeyEvent(event)
			if err != nil {
				return err
			}
		case KbActionPushLayer:
			err := kb.PushLayer(event.LayerName)
			if err != nil {
				return err
			}
		case KbActionPopLayer:
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
	case KbActionDown:
		return kb.vkb.KeyDown(int(event.KeyCode))
	case KbActionUp:
		return kb.vkb.KeyUp(int(event.KeyCode))
	case KbActionTap:
		return kb.vkb.KeyPress(int(event.KeyCode))
	default:
		return errors.New("unknown event")
	}
}

func elapsed(from, to syscall.Timeval) time.Duration {
	return time.Unix(to.Sec, to.Usec*1000).Sub(time.Unix(from.Sec, from.Usec*1000))
}
