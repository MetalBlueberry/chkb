package chkb

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	TapDelay = 200 * time.Millisecond
)

type Keyboard struct {
	*Captor
	*Mapper
	*Handler
	Config
}

func NewKeyboard(book Config, initialLayer string) *Keyboard {
	kb := &Keyboard{
		Captor:  NewCaptor(),
		Mapper:  NewMapper(),
		Handler: NewHandler(),
		Config:  book,
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
