package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

func main() {

	if len(os.Args) != 2 {
		panic("you must provide the device as the first param")
	}

	dev, err := evdev.Open(os.Args[1])
	if err != nil {
		fmt.Printf("unable to open input device: %s\n, %s", os.Args[1], err)
		os.Exit(1)
	}
	defer dev.File.Close()
	if err != nil {
		panic(err)
	}

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		return
	}
	defer keyboard.Close()

	defer dev.Release()
	err = dev.Grab()

	for {
		key, err := captureKey(dev)
		if err != nil {
			panic(err)
		}

		log.Printf("%s, %d", evdev.KEY[int(key.Code)], key.Code)

		sendEvnt(keyboard, key)
	}
}

func captureKey(dev *evdev.InputDevice) (*evdev.InputEvent, error) {
	for {
		event, err := dev.ReadOne()
		if err != nil {
			panic(err)
		}
		if event.Type != evdev.EV_KEY {
			continue
		}
		if evdev.KeyEventState(event.Value) == evdev.KeyHold && event.Code == evdev.KEY_ESC {
			return event, errors.New("done")
		}
		return event, nil
	}
}

func sendEvnt(keyboard uinput.Keyboard, event *evdev.InputEvent) {
	switch evdev.KeyEventState(event.Value) {
	case evdev.KeyDown:
		keyboard.KeyDown(int(event.Code))
	case evdev.KeyUp:
		keyboard.KeyUp(int(event.Code))
	}
}
