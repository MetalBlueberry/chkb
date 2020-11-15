package main

import (
	"encoding/json"
	"os"

	evdev "github.com/gvalkov/golang-evdev"
)

type KeyCapture interface {
	KeyEvent(evdev.InputEvent)
}

type Keyboard interface {
	// KeyPress will cause the key to be pressed and immediately released.
	KeyPress(key int) error

	// KeyDown will send a keypress event to an existing keyboard device.
	// The key can be any of the predefined keycodes from keycodes.go.
	// Note that the key will be "held down" until "KeyUp" is called.
	KeyDown(key int) error

	// KeyUp will send a keyrelease event to an existing keyboard device.
	// The key can be any of the predefined keycodes from keycodes.go.
	KeyUp(key int) error
}

type Mouse interface {
	// MoveLeft will move the mouse cursor left by the given number of pixel.
	MoveLeft(pixel int32) error

	// MoveRight will move the mouse cursor right by the given number of pixel.
	MoveRight(pixel int32) error

	// MoveUp will move the mouse cursor up by the given number of pixel.
	MoveUp(pixel int32) error

	// MoveDown will move the mouse cursor down by the given number of pixel.
	MoveDown(pixel int32) error

	// LeftClick will issue a single left click.
	LeftClick() error

	// RightClick will issue a right click.
	RightClick() error

	// LeftPress will simulate a press of the left mouse button. Note that the button will not be released until
	// LeftRelease is invoked.
	LeftPress() error

	// LeftRelease will simulate the release of the left mouse button.
	LeftRelease() error

	// RightPress will simulate the press of the right mouse button. Note that the button will not be released until
	// RightRelease is invoked.
	RightPress() error

	// RightRelease will simulate the release of the right mouse button.
	RightRelease() error

	// Wheel will simulate a wheel movement.
	Wheel(horizontal bool, delta int32) error
}

type Layer struct {
	Keys map[uint16]uint16
}

func NewLayer() *Layer {
	return &Layer{
		Keys: make(map[uint16]uint16),
	}
}

func NewLayerFrom(a, b Layout) *Layer {
	layer := NewLayer()

	for irow := range a.Keys {
		for icol := range a.Keys[irow] {
			layer.Add(a.Keys[irow][icol], b.Keys[irow][icol])
		}
	}
	return layer
}

func (layer *Layer) Add(key uint16, to uint16) error {
	layer.Keys[key] = to
	return nil
}

func (layer *Layer) Apply(event evdev.InputEvent) (*evdev.InputEvent, error) {
	code, ok := layer.Keys[event.Code]
	if ok {
		event.Code = code
	}
	return &event, nil
}

func (layer *Layer) Save(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	return encoder.Encode(&layer.Keys)
}

func (layer *Layer) Load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	layer.Keys = make(map[uint16]uint16)
	encoder := json.NewDecoder(f)
	return encoder.Decode(&layer.Keys)
}
