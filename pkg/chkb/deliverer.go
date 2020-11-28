package chkb

import log "github.com/sirupsen/logrus"

type Handler struct {
	Deliverers []Deliverer
}

func NewHandler() *Handler {
	return &Handler{
		Deliverers: []Deliverer{},
	}
}

type Deliverer interface {
	Deliver(event MapEvent) (handled bool, err error)
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
