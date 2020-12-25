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
	"container/list"
	"errors"
	"fmt"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
	log "github.com/sirupsen/logrus"
)

// Mapper converts KeyEvents into MapEvents based on the layer configuration.
type Mapper struct {
	Layers Layers

	downkeys        map[KeyCode][]MapEvent
	virtualDownKeys map[KeyCode]bool
	held            *list.List
	holding         KeyCode
}

func NewMapper() *Mapper {
	kb := &Mapper{
		Layers:          Layers{},
		downkeys:        map[KeyCode][]MapEvent{},
		virtualDownKeys: map[KeyCode]bool{},
		held:            list.New(),
	}
	return kb
}

func (mapper *Mapper) WithLayers(layers Layers) *Mapper {
	mapper.Layers = layers
	return mapper
}

type KeyCode uint16

func (keyCode KeyCode) String() string {
	return evdev.KEY[int(keyCode)]
}

type KeyEvent struct {
	Time    time.Time
	Action  KeyActions
	KeyCode KeyCode
}

func (ev KeyEvent) String() string {
	return fmt.Sprintf("%s-%v", evdev.KEY[int(ev.KeyCode)], ev.Action)
}

type MapEvent struct {
	// Time represents the time of the event that put is event in the queue and not the actual time.
	Time    time.Time `yaml:"-"`
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

func (layer *Layer) isKeyMapped(keyCode KeyCode) bool {
	_, ok := layer.KeyMap[keyCode]
	return ok
}

func (layer *Layer) findMap(event KeyEvent) ([]MapEvent, bool) {
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
				if m.Action.Is(KbActionMap) {
					m.Action = KbActions(event.Action)
				}
				m.Time = event.Time
				copymaps = append(copymaps, m)
			}
		}
	}

	kmaps, ok := keymap[event.Action]
	if ok {
		for i := range kmaps {
			m := kmaps[i]
			if m.Action.Is(KbActionMap) {
				m.Action = KbActions(event.Action)
			}
			m.Time = event.Time
			copymaps = append(copymaps, m)
		}
		return copymaps, true
	}
	return copymaps, len(copymaps) > 0
}

// HandleFunc is used to send MapEvents to handlers
type HandleFunc func(MapEvent) error

// Maps processes a slice of KeyEvents and calls handle func with the MapEvents
func (mapper *Mapper) Maps(events []KeyEvent, handle HandleFunc) error {
	for _, event := range events {
		mapper.held.PushBack(event)
	}

	for e := mapper.held.Front(); e != nil; {
		event := e.Value.(KeyEvent)

		if mapper.holding != KEY_RESERVED && mapper.holding != event.KeyCode {
			e = e.Next()
			continue
		}
		mapper.held.Remove(e)
		e = mapper.held.Front()

		switch event.Action {
		case KeyActionDown:
			keyMap, ok := mapper.Layers.findKeyMap(event.KeyCode)
			if ok && keyMap.hasSpecialMaps() {
				mapper.holding = event.KeyCode
			}
		default:
			mapper.holding = KEY_RESERVED
		case KeyActionUp:
		}

		switch event.Action {
		case KeyActionUp:
			// If key was down, raise virtual keys down
			downkeys, ok := mapper.downkeys[event.KeyCode]
			if ok {
				for _, downkey := range downkeys {
					mapper.virtualDownKeys[downkey.KeyCode] = false
					err := handle(MapEvent{
						Time:    event.Time,
						Action:  KbActionUp,
						KeyCode: downkey.KeyCode,
					})
					if err != nil {
						return err
					}
				}
				delete(mapper.downkeys, event.KeyCode)
			}
		}

		maps, err := mapper.mapOne(event)
		if err != nil {
			log.
				WithField("event", event).
				WithError(err).
				Trace("Ignored event")
			continue
		}

		for _, m := range maps {
			switch m.Action {
			case KbActionDown:
				mapper.downkeys[event.KeyCode] = append(mapper.downkeys[event.KeyCode], m)
				mapper.virtualDownKeys[m.KeyCode] = true
				err := handle(m)
				if err != nil {
					return err
				}
			case KbActionUp:
				isDown, ok := mapper.virtualDownKeys[m.KeyCode]
				if ok && isDown {
					mapper.virtualDownKeys[m.KeyCode] = false
					err := handle(MapEvent{
						Time:    m.Time,
						Action:  KbActionUp,
						KeyCode: m.KeyCode,
					})
					if err != nil {
						return err
					}
					delete(mapper.downkeys, m.KeyCode)
				}
			default:
				err := handle(m)
				if err != nil {
					return err
				}
			case KbActionNil:
			}
		}

	}
	return nil
}

func (layers Layers) findKeyMap(keyCode KeyCode) (KeyMapActions, bool) {
	for i := len(layers) - 1; i >= 0; i-- {
		keymap, ok := layers[i].KeyMap[keyCode]
		if ok {
			return keymap, true
		}
	}
	return nil, false
}

func (mapper *Mapper) mapOne(event KeyEvent) ([]MapEvent, error) {
	mapped := make([]MapEvent, 0)
	handled := false
	for i := len(mapper.Layers) - 1; i >= 0; i-- {
		kmaps, ok := mapper.Layers[i].findMap(event)
		if ok {
			mapped = append(mapped, kmaps...)
			handled = true
			break
		}
		if (!mapper.Layers[i].isKeyMapped(event.KeyCode)) && (event.Action == KeyActionDown) && len(mapper.Layers[i].OnMiss) > 0 {
			transparent := false
			for j := range mapper.Layers[i].OnMiss {
				if mapper.Layers[i].OnMiss[j].Action == KbActionMap {
					transparent = true
				} else {
					mapped = append(mapped, mapper.Layers[i].OnMiss[j])
				}
			}
			if !transparent {
				handled = true
				break
			}
		}
	}

	// No map detected, forward
	if !handled {
		switch event.Action {
		case KeyActionUp, KeyActionDown:
			mapped = append(mapped,
				MapEvent{
					Time:    event.Time,
					Action:  KbActions(event.Action),
					KeyCode: event.KeyCode,
				},
			)
		}
	}

	if len(mapped) > 0 {
		return mapped, nil
	}

	return nil, errors.New("Ignore event")
}

func (mapper *Mapper) addLayer(layer *Layer) {
	mapper.Layers = append(mapper.Layers, layer)
}

func (mapper *Mapper) removeLayer(layer *Layer) error {
	if len(mapper.Layers) == 1 {
		return fmt.Errorf("Cannot remove last layer")
	}
	for i := range mapper.Layers {
		if mapper.Layers[i] == layer {
			mapper.Layers = append(mapper.Layers[:i], mapper.Layers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Layer not previously applied")
}

func (mapper *Mapper) setLayer(layer *Layer) error {
	mapper.Layers[0] = layer
	return nil
}
