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
	KbActionNil       KbActions = "Nil"
	KbActionMap       KbActions = "Map"
	KbActionDown      KbActions = "Down"
	KbActionUp        KbActions = "Up"
	KbActionTap       KbActions = "Tap"
	KbActionDoubleTap KbActions = "DoubleTap"
	KbActionHold      KbActions = "Hold"
	KbActionPushLayer KbActions = "PushLayer"
	KbActionPopLayer  KbActions = "PopLayer"
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
