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

import log "github.com/sirupsen/logrus"

type Deliverer interface {
	Deliver(event MapEvent) (handled bool, err error)
}
type Handler struct {
	Deliverers []Deliverer
}

func NewHandler() *Handler {
	return &Handler{
		Deliverers: []Deliverer{},
	}
}

type KeyActions string

const (
	KeyActionNil       KeyActions = "Nil"
	KeyActionMap       KeyActions = "Map"
	KeyActionDown      KeyActions = "Down"
	KeyActionUp        KeyActions = "Up"
	KeyActionTap       KeyActions = "Tap"
	KeyActionDoubleTap KeyActions = "DoubleTap"
	KeyActionHold      KeyActions = "Hold"
)

func (action KeyActions) String() string {
	return string(action)
}

type KbActions string

const (
	KbActionNil         KbActions = "Nil"
	KbActionMap         KbActions = "Map"
	KbActionDown        KbActions = "Down"
	KbActionUp          KbActions = "Up"
	KbActionTap         KbActions = "Tap"
	KbActionDoubleTap   KbActions = "DoubleTap"
	KbActionHold        KbActions = "Hold"
	KbActionPushLayer   KbActions = "PushLayer"
	KbActionPopLayer    KbActions = "PopLayer"
	KbActionChangeLayer KbActions = "ChangeLayer"
)

func (action KbActions) String() string {
	return string(action)
}

func (action KbActions) Is(compare KbActions) bool {
	if action == "" {
		action = KbActionMap
	}
	return action == compare
}

var KbActionsMap map[string]KbActions = map[string]KbActions{
	KbActionNil.String():         KbActionNil,
	KbActionMap.String():         KbActionMap,
	KbActionDown.String():        KbActionDown,
	KbActionUp.String():          KbActionUp,
	KbActionTap.String():         KbActionTap,
	KbActionDoubleTap.String():   KbActionDoubleTap,
	KbActionHold.String():        KbActionHold,
	KbActionPushLayer.String():   KbActionPushLayer,
	KbActionPopLayer.String():    KbActionPopLayer,
	KbActionChangeLayer.String(): KbActionChangeLayer,
}

func (handler *Handler) AddDeliverer(deliverer Deliverer) {
	handler.Deliverers = append(handler.Deliverers, deliverer)
}

func (kb *Handler) Handle(event MapEvent) error {
	handled := false
	for i := range kb.Deliverers {
		delivered, err := kb.Deliverers[i].Deliver(event)
		if err != nil {
			return err
		}
		if delivered {
			handled = true
		}
	}
	if !handled {
		log.
			WithField("Action", event.Action).
			Debug("Ignored event")
	}
	return nil
}
