/*
Emulate the KBC Poker 3 keyboard layout on a standard keyboard by remapping the
Caps_Lock as the Poker 3 Fn key. This needs to be run as root and it's probably
a good idea to disable the regular function of the Caps_Lock key as well.
License: WTFPL
*/

package main

import (
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
	defer dev.Release()
	err = dev.Grab()
	if err != nil {
		panic(err)
	}

	log.Println(dev.Capabilities)

	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("testmouse"))
	if err != nil {
		return
	}
	// always do this after the initialization in order to guarantee that the device will be properly closed
	defer mouse.Close()

	for {

		events, err := dev.Read()
		if err != nil {
			panic(err)
		}

		for _, event := range events {
			code := event.Code

			if event.Type == evdev.EV_REL {
				log.Println(evdev.NewRelEvent(&event).String())
				log.Println(evdev.REL[int(event.Code)])
				switch event.Code {
				case evdev.REL_X:
					mouse.MoveRight(event.Value)
				case evdev.REL_Y:
					mouse.MoveDown(event.Value)
				}
			} else {
				log.Printf("event: %s", event.String())
				log.Println(evdev.ByEventType[int(event.Type)][int(event.Code)])
			}

			if code == 273 {
				log.Print("done")
				return
			}
		}
		log.Println("End step")
	}
}
