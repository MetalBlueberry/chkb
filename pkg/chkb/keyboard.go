package chkb

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Keyboard struct {
	*Captor
	*Handler

	Layers    []*Layer
	LayerBook Book
	downkeys  map[KeyCode]MapEvent
}

func NewKeyboard(book Book, initialLayer string) *Keyboard {
	kb := &Keyboard{
		Captor:  NewCaptor(),
		Handler: NewHandler(),

		LayerBook: book,
		Layers:    []*Layer{},
		downkeys:  map[KeyCode]MapEvent{},
	}
	kb.PushLayer(initialLayer)
	kb.AddDeliverer(kb)
	return kb
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
	KbActionNil.String():       KbActionNil,
	KbActionMap.String():       KbActionMap,
	KbActionDown.String():      KbActionDown,
	KbActionUp.String():        KbActionUp,
	KbActionTap.String():       KbActionTap,
	KbActionDoubleTap.String(): KbActionDoubleTap,
	KbActionHold.String():      KbActionHold,
	KbActionPushLayer.String(): KbActionPushLayer,
	KbActionPopLayer.String():  KbActionPopLayer,
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

func NewKeyEv(event InputEvent, action KeyActions) KeyEvent {
	return KeyEvent{
		// Time:    time.Unix(event.Time.Sec, event.Time.Usec*1000),
		KeyCode: event.KeyCode,
		Action:  action,
	}
}
func NewKeyInputEvent(event evdev.InputEvent) InputEvent {
	ie := InputEvent{
		Time:    time.Unix(event.Time.Sec, event.Time.Usec*1000),
		KeyCode: KeyCode(event.Code),
	}
	switch evdev.KeyEventState(event.Value) {
	case evdev.KeyDown:
		ie.Action = InputActionDown
	case evdev.KeyUp:
		ie.Action = InputActionUp
	case evdev.KeyHold:
		ie.Action = InputActionHold
	}
	return ie
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

		maps, err := kb.mapOne(event)
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

func (kb *Keyboard) mapOne(event KeyEvent) ([]MapEvent, error) {
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

func (kb *Keyboard) Deliver(event MapEvent) (bool, error) {
	switch event.Action {
	case KbActionPushLayer:
		err := kb.PushLayer(event.LayerName)
		return true, err
	case KbActionPopLayer:
		err := kb.PopLayer(event.LayerName)
		return true, err
	default:
		return false, nil
	}
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

func (kb *Keyboard) PopLayer(name string) error {
	log.Printf("Pop layer")
	if len(kb.Layers) == 1 {
		return errors.New("You cannot pop the last layer")
	}
	for i := range kb.Layers {
		if kb.Layers[i] == kb.LayerBook[name] {
			kb.Layers = append(kb.Layers[:i], kb.Layers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Layer %s not found", name)
}

func elapsed(from, to syscall.Timeval) time.Duration {
	return time.Unix(to.Sec, to.Usec*1000).Sub(time.Unix(from.Sec, from.Usec*1000))
}
