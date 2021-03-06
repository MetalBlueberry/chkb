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

// Package chkb turns a regular keyboard intro a fully programmable keyboard.
// The Keyboard struct is the high level entry point that you can use to
// access the keyboard functionality
package chkb

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Keyboard is a cheap programmable keyboard
// It coordinates the communication between the different modules
type Keyboard struct {
	Captor *Captor
	Mapper *Mapper
	*Handler
	Config
}

// NewKeyboard creates a new chkb keyboard
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

// Run starts a capture loop,
func (kb *Keyboard) Run(event CaptureFunc) error {
	return kb.Captor.Run(event, func(captured []KeyEvent) error {
		err := kb.Mapper.Maps(captured, kb.Handler.Handle)
		if err != nil {
			return err
		}
		return nil
	})
}

// PushLayer adds a new layer to the keyboard
func (kb *Keyboard) PushLayer(name string) (err error) {
	log.
		WithField("name", name).
		Info("Push layer")
	layer, ok := kb.Config.Layers[name]
	if !ok {
		return fmt.Errorf("Layer %s not found", name)
	}
	kb.Mapper.addLayer(layer)
	return nil
}

// PopLayer removes a layer from the keyboard
func (kb *Keyboard) PopLayer(name string) (err error) {
	log.
		WithField("name", name).
		Info("Pop layer")
	layer, ok := kb.Config.Layers[name]
	if !ok {
		return fmt.Errorf("Layer %s not found", name)
	}
	err = kb.Mapper.removeLayer(layer)
	if err != nil {
		return err
	}
	return nil
}

// ChangeLayer changes the base layer for the keyboard
func (kb *Keyboard) ChangeLayer(name string) (err error) {
	log.
		WithField("name", name).
		Info("Change layer")
	layer, ok := kb.Config.Layers[name]
	if !ok {
		return fmt.Errorf("Layer %s not found", name)
	}
	err = kb.Mapper.setLayer(layer)
	if err != nil {
		return err
	}
	return nil
}

// Deliver implements Deliverer interface to handle Layer events
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
