package chkb

import (
	"io"

	"gopkg.in/yaml.v3"
	// log "github.com/sirupsen/logrus"
)

type Book map[string]*Layer

func (b *Book) Save(w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	return encoder.Encode(b)
}

func (b *Book) Load(r io.Reader) error {
	encoder := yaml.NewDecoder(r)
	return encoder.Decode(b)
}

type Layer struct {
	DefaultMap *MapEvent `yaml:"defaultMap,omitempty"`
	KeyMap     KeyMap    `yaml:"keyMap,omitempty"`
}

type KeyMap map[KeyCode]map[KeyActions][]MapEvent

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
		(*km) = map[KeyCode]map[KeyActions][]MapEvent{}
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

// func (layer *Layer) UnmarshalYAML(value *yaml.Node) error {
// 	return value.Decode(&layer.KeyMap)
// }

func (km KeyMap) MarshalYAML() (interface{}, error) {
	return km.StringMap(), nil
}
