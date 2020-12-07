package chkb

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
	// log "github.com/sirupsen/logrus"
)

type Config struct {
	Layers LayerBook `yaml:"layers"`
}

func (b *Config) Save(w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	return encoder.Encode(b)
}

func (b *Config) Load(r io.Reader) error {
	encoder := yaml.NewDecoder(r)
	return encoder.Decode(b)
}

type LayerBook map[string]*Layer
type Layers []*Layer

type Layer struct {
	OnMiss []MapEvent `yaml:"onMiss,omitempty"`
	KeyMap KeyMap     `yaml:"keyMap,omitempty"`
}

type KeyMapActions map[KeyActions][]MapEvent

func (km KeyMapActions) hasSpecialMaps() bool {
	for k := range km {
		switch k {
		case KeyActionUp, KeyActionDown, KeyActionNil, KeyActionMap:
			continue
		default:
			return true
		}
	}
	return false
}

type KeyMap map[KeyCode]KeyMapActions

func (km *KeyMap) StringMap() map[string]map[string][]MapEvent {
	tmp := make(map[string]map[string][]MapEvent)
	for keyCode, actionMap := range *km {
		tmp[keyCode.String()] = make(map[string][]MapEvent)
		for action, mapEvent := range actionMap {
			tmp[keyCode.String()][action.String()] = mapEvent
		}
	}
	return tmp
}

func (km *KeyMap) FromStringMap(source map[string]map[string][]MapEvent) error {
	if (*km) == nil {
		(*km) = map[KeyCode]KeyMapActions{}
	}
	for keyCodeString, actionMap := range source {
		keyCode, err := ParseKeyCode(keyCodeString)
		if err != nil {
			return err
		}
		(*km)[keyCode] = make(map[KeyActions][]MapEvent)

		for actionString, mapEvent := range actionMap {
			action, err := ParseKeyAction(actionString)
			if err != nil {
				return err
			}
			(*km)[keyCode][action] = mapEvent
		}
	}
	return nil
}

func (km *KeyMap) UnmarshalYAML(value *yaml.Node) error {
	tmp := make(map[string]map[string][]MapEvent)
	err := value.Decode(&tmp)
	if err != nil {
		return err
	}
	return km.FromStringMap(tmp)
}

func (km KeyMap) MarshalYAML() (interface{}, error) {
	return km.StringMap(), nil
}

var KeyActionsMap map[string]KeyActions = map[string]KeyActions{
	KeyActionNil.String():       KeyActionNil,
	KeyActionMap.String():       KeyActionMap,
	KeyActionDown.String():      KeyActionDown,
	KeyActionUp.String():        KeyActionUp,
	KeyActionTap.String():       KeyActionTap,
	KeyActionDoubleTap.String(): KeyActionDoubleTap,
	KeyActionHold.String():      KeyActionHold,
}

func ParseKeyAction(value string) (KeyActions, error) {
	if value == "" {
		return KeyActionMap, nil
	}
	a, ok := KeyActionsMap[value]
	if !ok {
		return a, fmt.Errorf("Action %s not found", value)
	}
	return a, nil
}

func ParseKbAction(value string) (KbActions, error) {
	if value == "" {
		return KbActionMap, nil
	}
	a, ok := KbActionsMap[value]
	if !ok {
		return a, fmt.Errorf("Action %s not found", value)
	}
	return a, nil
}

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
