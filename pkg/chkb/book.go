package chkb

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"
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
	KeyMap map[KeyCode]map[Actions]MapEvent
}

func (layer *Layer) StringMap() map[string]map[string]MapEvent {
	tmp := make(map[string]map[string]MapEvent)
	for keyCode, actionMap := range layer.KeyMap {
		tmp[keyCode.String()] = make(map[string]MapEvent)
		for action, mapEvent := range actionMap {
			tmp[keyCode.String()][action.String()] = mapEvent
		}
	}
	return tmp
}

func (layer *Layer) FromStringMap(source map[string]map[string]MapEvent) error {
	if layer.KeyMap == nil {
		layer.KeyMap = map[KeyCode]map[Actions]MapEvent{}
	}
	for keyCodeString, actionMap := range source {
		keyCode, err := ParseKeyCode(keyCodeString)
		if err != nil {
			return err
		}
		layer.KeyMap[keyCode] = make(map[Actions]MapEvent)

		for actionString, mapEvent := range actionMap {
			action, err := ParseAction(actionString)
			if err != nil {
				return err
			}
			layer.KeyMap[keyCode][action] = mapEvent
		}
	}
	return nil
}

func (layer *Layer) UnmarshalYAML(value *yaml.Node) error {
	tmp := make(map[string]map[string]MapEvent)
	err := value.Decode(tmp)
	if err != nil {
		return err
	}
	return layer.FromStringMap(tmp)
}

func (layer *Layer) MarshalYAML() (interface{}, error) {
	return layer.StringMap(), nil
}

func (layer *Layer) MarshalJSON() ([]byte, error) {
	return json.Marshal(layer.StringMap())
}
