package main

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"
	"MetalBlueberry/cheap-keyboard/pkg/deliverers/layerFile"
	"MetalBlueberry/cheap-keyboard/pkg/deliverers/vkb"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
	"github.com/spf13/afero"
)

func main() {

	dev, err := evdev.Open(os.Args[1])
	if err != nil {
		fmt.Printf("unable to open input device: %s\n, %s", os.Args[1], err)
		os.Exit(1)
	}
	defer dev.File.Close()
	if err != nil {
		log.Fatal(err)
	}

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	book := chkb.Book{}

	keys, err := os.Open("keys.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = book.Load(keys)
	keys.Close()
	if err != nil {
		log.Fatal(err)
	}

	kb := chkb.NewMapper(
		book,
		"base",
	)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	lf, err := layerFile.NewLayerFile(afero.NewOsFs(), kb, filepath.Join(usr.HomeDir, ".chkb_layout"))
	if err != nil {
		log.Fatal(err)
	}

	defer dev.Release()
	err = dev.Grab()
	if err != nil {
		log.Fatal(err)
	}

	c := chkb.NewCaptor()
	h := chkb.NewHandler()

	h.AddDeliverer(&vkb.Keyboard{keyboard})
	h.AddDeliverer(kb)
	h.AddDeliverer(lf)

	for {
		events, err := dev.Read()
		if err != nil {
			panic(err)
		}
		for _, event := range events {
			if evdev.KeyEventState(event.Value) == evdev.KeyHold && event.Code == evdev.KEY_ESC {
				panic("exit")
			}
		}

		inputEvents := make([]chkb.InputEvent, 0)
		for i := range events {
			if events[i].Type != evdev.EV_KEY {
				continue
			}
			inputEvents = append(inputEvents, NewKeyInputEvent(events[i]))
		}

		captured, err := c.Capture(inputEvents)
		if err != nil {
			panic(err)
		}
		maps, err := kb.Maps(captured)
		if err != nil {
			panic(err)
		}
		err = h.Delivers(maps)
		if err != nil {
			panic(err)
		}
	}
}

func NewKeyInputEvent(event evdev.InputEvent) chkb.InputEvent {
	ie := chkb.InputEvent{
		Time:    time.Unix(event.Time.Sec, event.Time.Usec*1000),
		KeyCode: chkb.KeyCode(event.Code),
	}
	switch evdev.KeyEventState(event.Value) {
	case evdev.KeyDown:
		ie.Action = chkb.InputActionDown
	case evdev.KeyUp:
		ie.Action = chkb.InputActionUp
	case evdev.KeyHold:
		ie.Action = chkb.InputActionHold
	}
	return ie
}
