/*
Copyright © 2020 Víctor Pérez @MetalBlueberry

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package chkb

import (
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v3"
	// log "github.com/sirupsen/logrus"
)

const (
	// DefaultTapDelay is the default setting for TapDelay
	DefaultTapDelay = 200 * time.Millisecond
)

// Config is the configuration file structure
type Config struct {
	TapDelay int       `yaml:"tapDelay,omitempty"`
	Layers   LayerBook `yaml:"layers"`
}

// Save writes the configuration content to a io.Writer
func (b *Config) Save(w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	return encoder.Encode(b)
}

// Load reads the configuration from a io.Reader
func (b *Config) Load(r io.Reader) error {
	encoder := yaml.NewDecoder(r)
	return encoder.Decode(b)
}

// GetTapDelay returns the configured delay or the default if it is empty
func (b *Config) GetTapDelay() time.Duration {
	if b.TapDelay != 0 {
		return time.Duration(b.TapDelay) * time.Millisecond
	}
	return DefaultTapDelay
}

// LayerBook contains all the layer configurations
type LayerBook map[string]*Layer

// Layers is a list of Layer
type Layers []*Layer

// Layer represents a behaviour of the keyboard that can be applied on demand
type Layer struct {
	OnMiss []MapEvent `yaml:"onMiss,omitempty"`
	KeyMap KeyMap     `yaml:"keyMap,omitempty"`
}

// KeyMapActions contains the events mapped to an action
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

// KeyMap contains the keys with KeyMapActions
type KeyMap map[KeyCode]KeyMapActions

// StringMap generates the equivalent map with types removed, used for serialization
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

// FromStringMap reverses the action from StringMap
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

// UnmarshalYAML reads yaml configuration
func (km *KeyMap) UnmarshalYAML(value *yaml.Node) error {
	tmp := make(map[string]map[string][]MapEvent)
	err := value.Decode(&tmp)
	if err != nil {
		return err
	}
	return km.FromStringMap(tmp)
}

// MarshalYAML writes yaml configuration
func (km KeyMap) MarshalYAML() (interface{}, error) {
	return km.StringMap(), nil
}

// keyActionsMap is used to validate if a certain action exists
var keyActionsMap map[string]KeyActions = map[string]KeyActions{
	KeyActionNil.String():       KeyActionNil,
	KeyActionMap.String():       KeyActionMap,
	KeyActionDown.String():      KeyActionDown,
	KeyActionUp.String():        KeyActionUp,
	KeyActionTap.String():       KeyActionTap,
	KeyActionDoubleTap.String(): KeyActionDoubleTap,
	KeyActionHold.String():      KeyActionHold,
}

// ParseKeyAction returns the KeyAction from a string, or an error if it is invalid
func ParseKeyAction(value string) (KeyActions, error) {
	if value == "" {
		return KeyActionMap, nil
	}
	a, ok := keyActionsMap[value]
	if !ok {
		return a, fmt.Errorf("Action %s not found", value)
	}
	return a, nil
}

// ParseKbAction returns the KbAction from a string, or an error if it is invalid
func ParseKbAction(value string) (KbActions, error) {
	if value == "" {
		return KbActionMap, nil
	}
	a, ok := kbActionsMap[value]
	if !ok {
		return a, fmt.Errorf("Action %s not found", value)
	}
	return a, nil
}

// ParseKeyCode returns the KeyCode from a string, or an error if it is invalid
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

// UnmarshalYAML reads yaml configuration
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

// UnmarshalYAML reads yaml configuration
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

// MarshalYAML writes yaml configuration
func (action KeyActions) MarshalYAML() (interface{}, error) {
	return action.String(), nil
}

// MarshalYAML writes yaml configuration
func (action KbActions) MarshalYAML() (interface{}, error) {
	return action.String(), nil
}

// UnmarshalYAML reads yaml configuration
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

// MarshalYAML writes yaml configuration
func (keyCode KeyCode) MarshalYAML() (interface{}, error) {
	return keyCode.String(), nil
}
