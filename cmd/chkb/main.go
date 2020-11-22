package main

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"
	"fmt"
	"os"

	"github.com/bendahl/uinput"
	"github.com/eiannone/keyboard"
	evdev "github.com/gvalkov/golang-evdev"
)

func main() {
	dev, err := evdev.Open(os.Args[1])
	if err != nil {
		fmt.Printf("unable to open input device: %s\n, %s", os.Args[1], err)
		os.Exit(1)
	}
	defer dev.File.Close()
	if err != nil {
		panic(err)
	}

	vkb, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		return
	}
	defer keyboard.Close()

	book := chkb.Book{}

	keys, err := os.Open("keys.yaml")
	if err != nil {
		panic(err)
	}
	err = book.Load(keys)
	keys.Close()
	if err != nil {
		panic(err)
	}

	kb := chkb.NewKeyboard(
		book,
		"base",
		vkb,
	)

	defer dev.Release()
	err = dev.Grab()
	if err != nil {
		panic(err)
	}

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

		captured, err := kb.Capture(events)
		if err != nil {
			panic(err)
		}
		maps, err := kb.Maps(captured)
		if err != nil {
			panic(err)
		}
		err = kb.Deliver(maps)
		if err != nil {
			panic(err)
		}
	}
}
