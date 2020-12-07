package vkb

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"

	log "github.com/sirupsen/logrus"

	"github.com/bendahl/uinput"
)

type Keyboard struct {
	uinput.Keyboard
}

func (kb *Keyboard) Deliver(event chkb.MapEvent) (bool, error) {
	log.
		WithField("Action", event.Action).
		WithField("Key", event.KeyCode).
		Debug("Key Event")
	switch event.Action {
	case chkb.KbActionDown:
		return true, kb.KeyDown(int(event.KeyCode))
	case chkb.KbActionUp:
		return true, kb.KeyUp(int(event.KeyCode))
	case chkb.KbActionTap:
		return true, kb.KeyPress(int(event.KeyCode))
	default:
		return false, nil
	}
}
