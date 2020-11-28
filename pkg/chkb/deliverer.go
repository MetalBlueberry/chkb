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

// KbAction are the actions that can assigned to a KeyEvent
type KbActions int

//go:generate stringer -type=KbActions -trimprefix KbAction

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

func (handler *Handler) AddDeliverer(deliverer Deliverer) {
	handler.Deliverers = append(handler.Deliverers, deliverer)
}

func (kb *Handler) Delivers(events []MapEvent) error {
	for _, event := range events {
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
	}
	return nil
}
