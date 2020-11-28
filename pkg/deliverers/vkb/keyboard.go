package vkb

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"

	"github.com/bendahl/uinput"
)

type Keyboard struct {
	uinput.Keyboard
}

func (kb *Keyboard) Deliver(event chkb.MapDefinition) (bool, error) {
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
