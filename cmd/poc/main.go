package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("testmouse"))
	if err != nil {
		return
	}
	// always do this after the initialization in order to guarantee that the device will be properly closed
	defer mouse.Close()

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		return
	}
	defer keyboard.Close()

	// log.Println("Starting...")
	// time.Sleep(2 * time.Second)

	defer dev.Release()
	err = dev.Grab()

	log.Println("Ready")
	for {
		process(dev, mouse, keyboard)
	}
}

func process(dev *evdev.InputDevice, mouse uinput.Mouse, keyboard uinput.Keyboard) {

	events, err := dev.Read()
	if err != nil {
		panic(err)
	}

	for _, event := range events {
		// log.Printf("event: %s", event.String())
		// log.Println(evdev.ByEventType[int(event.Type)][int(event.Code)])

		if event.Type == evdev.EV_KEY {
			handleKey(dev, mouse, keyboard, event)
		}

	}
	log.Println("End step")
}

var active bool
var lastKey time.Time

func handleKey(dev *evdev.InputDevice, mouse uinput.Mouse, keyboard uinput.Keyboard, event evdev.InputEvent) {
	if !active {
		switch event.Code {
		case evdev.KEY_LEFTALT:
			if event.Value == 1 {
				active = true
			}
		default:
			switch evdev.KeyEventState(event.Value) {
			case evdev.KeyDown:
				log.Print("key down")
				keyboard.KeyDown(int(event.Code))
			case evdev.KeyHold:
				log.Print("key repeat")
				keyboard.KeyPress(int(event.Code))
			case evdev.KeyUp:
				log.Print("key up")
				keyboard.KeyUp(int(event.Code))
			}
		}
		return
	}

	quick := 175 * time.Millisecond
	lastKeyElapsed := time.Since(lastKey)

	log.Print("elapsed %s", lastKeyElapsed)
	speed := int32(25)
	switch event.Code {
	case evdev.KEY_LEFTALT:
		switch evdev.KeyEventState(event.Value) {
		case evdev.KeyUp:
			active = false
		}

	case evdev.KEY_J:
		switch evdev.KeyEventState(event.Value) {
		case evdev.KeyDown:
			if lastKeyElapsed < quick {
				speed *= 5
			}
			mouse.MoveDown(speed)
			lastKey = time.Now()
		}
	case evdev.KEY_K:
		switch evdev.KeyEventState(event.Value) {
		case evdev.KeyDown:
			if lastKeyElapsed < quick {
				speed *= 5
			}
			mouse.MoveUp(speed)
			lastKey = time.Now()
		}
	case evdev.KEY_L:
		switch evdev.KeyEventState(event.Value) {
		case evdev.KeyDown:
			if lastKeyElapsed < quick {
				speed *= 5
			}
			mouse.MoveRight(speed)
			lastKey = time.Now()
		}
	case evdev.KEY_H:
		switch evdev.KeyEventState(event.Value) {
		case evdev.KeyDown:
			if lastKeyElapsed < quick {
				speed *= 5
			}
			mouse.MoveLeft(speed)
			lastKey = time.Now()
		}

	case evdev.KEY_W:
		switch evdev.KeyEventState(event.Value) {
		case evdev.KeyDown:
			mouse.LeftPress()
		case evdev.KeyHold:
			mouse.LeftClick()
		case evdev.KeyUp:
			mouse.LeftRelease()
		}
	case evdev.KEY_Q:
		switch evdev.KeyEventState(event.Value) {
		case evdev.KeyDown:
			mouse.RightPress()
		case evdev.KeyHold:
			mouse.RightClick()
		case evdev.KeyUp:
			mouse.RightRelease()
		}
	case evdev.KEY_ESC:
		panic("exit")
		// default:
		// 	switch evdev.KeyEventState(event.Value) {
		// 	case evdev.KeyDown:
		// 		keyboard.KeyDown(int(event.Code))
		// 	case evdev.KeyHold:
		// 		keyboard.KeyPress(int(event.Code))
		// 	case evdev.KeyUp:
		// 		keyboard.KeyUp(int(event.Code))
		// 	}
	}
}
