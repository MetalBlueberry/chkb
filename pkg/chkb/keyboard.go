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

	log "github.com/sirupsen/logrus"
)

type Keyboard struct {
	*Captor
	*Mapper
	*Handler
	Config
}

func NewKeyboard(config Config, initialLayer string) *Keyboard {
	kb := &Keyboard{
		Captor:  NewCaptor(&config),
		Mapper:  NewMapper(),
		Handler: NewHandler(),
		Config:  config,
	}
	kb.AddDeliverer(kb)
	kb.PushLayer(initialLayer)
	return kb
}

func (kb *Keyboard) Run(event func() ([]InputEvent, error)) error {
	return kb.Captor.Run(event, func(captured []KeyEvent) error {
		err := kb.Maps(captured, kb.Handler.Handle)
		if err != nil {
			return err
		}
		return nil
	})
}

func (kb *Keyboard) Deliver(event MapEvent) (bool, error) {
	switch event.Action {
	case KbActionPushLayer:
		err := kb.PushLayer(event.LayerName)
		return true, err
	case KbActionPopLayer:
		err := kb.PopLayer(event.LayerName)
		return true, err
	case KbActionChangeLayer:
		err := kb.ChangeLayer(event.LayerName)
		return true, err
	default:
		return false, nil
	}
}

func (kb *Keyboard) PushLayer(name string) (err error) {
	log.
		WithField("name", name).
		Info("Push layer")
	layer, ok := kb.Config.Layers[name]
	if !ok {
		return fmt.Errorf("Layer %s not found", name)
	}
	kb.addLayer(layer)
	return nil
}

func (kb *Keyboard) PopLayer(name string) (err error) {
	log.
		WithField("name", name).
		Info("Pop layer")
	layer, ok := kb.Config.Layers[name]
	if !ok {
		return fmt.Errorf("Layer %s not found", name)
	}
	err = kb.removeLayer(layer)
	if err != nil {
		return err
	}
	return nil
}

func (kb *Keyboard) ChangeLayer(name string) (err error) {
	log.
		WithField("name", name).
		Info("Change layer")
	layer, ok := kb.Config.Layers[name]
	if !ok {
		return fmt.Errorf("Layer %s not found", name)
	}
	err = kb.setLayer(layer)
	if err != nil {
		return err
	}
	return nil
}
