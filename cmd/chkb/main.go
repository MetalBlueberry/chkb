package main

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"
	"MetalBlueberry/cheap-keyboard/pkg/deliverers/layerFile"
	"MetalBlueberry/cheap-keyboard/pkg/deliverers/vkb"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

var (
	Version = "development"
)

func main() {
	var (
		verbose bool
	)
	flag.BoolVar(&verbose, "v", false, "print debug information")
	flag.Parse()
	if verbose {
		log.SetLevel(log.DebugLevel)
		log.
			WithField("Version", Version).
			Debug("Set debug level")
	}

	// This is required to ensure that the enter key is not stuck down
	// when the device is grab.
	log.Info("You have 200 ms to release all keys")
	time.Sleep(200 * time.Millisecond)

	devs := make([]*evdev.InputDevice, 0)
	for _, arg := range flag.Args() {
		dev, err := evdev.Open(arg)
		if err != nil {
			fmt.Printf("unable to open input device: %s\n, %s", os.Args[1], err)
			os.Exit(1)
		}
		defer dev.File.Close()
		if err != nil {
			log.Fatal(err)
		}

		defer dev.Release()
		err = dev.Grab()
		if err != nil {
			log.Fatal(err)
		}

		devs = append(devs, dev)
	}

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	book := chkb.Config{}

	keys, err := os.Open("keys.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = book.Load(keys)
	keys.Close()
	if err != nil {
		log.Fatal(err)
	}
	kb := chkb.NewKeyboard(book, "base")

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	lf, err := layerFile.NewLayerFile(afero.NewOsFs(), kb, filepath.Join(usr.HomeDir, ".chkb_layout"))
	if err != nil {
		log.Fatal(err)
	}
	kb.AddDeliverer(lf)

	kb.AddDeliverer(&vkb.Keyboard{keyboard})

	events := make(chan []chkb.InputEvent)
	for _, dev := range devs {
		go capture(dev, events)
	}

	kb.Run(func() ([]chkb.InputEvent, error) {
		inputEvents, ok := <-events
		if !ok {
			return nil, errors.New("Closed chan")
		}
		return inputEvents, nil
	})
}

var ErrInvalidEvent = errors.New("Invalid event")

func NewKeyInputEvent(event evdev.InputEvent) (chkb.InputEvent, error) {
	ie := chkb.InputEvent{
		Time:    time.Unix(event.Time.Sec, event.Time.Usec*1000),
		KeyCode: chkb.KeyCode(event.Code),
	}
	switch evdev.KeyEventState(event.Value) {
	case evdev.KeyDown:
		ie.Action = chkb.InputActionDown
	case evdev.KeyUp:
		ie.Action = chkb.InputActionUp
	default:
		return chkb.InputEvent{}, ErrInvalidEvent
	}
	return ie, nil
}

func capture(dev *evdev.InputDevice, evs chan []chkb.InputEvent) {
	for {
		events, err := dev.Read()
		if err != nil {
			log.WithError(err).Error("Device closed")
			close(evs)
			return
		}
		for _, event := range events {
			if evdev.KeyEventState(event.Value) == evdev.KeyHold && event.Code == evdev.KEY_ESC {
				panic(err)
			}
		}

		inputEvents := make([]chkb.InputEvent, 0)
		for i := range events {
			if events[i].Type != evdev.EV_KEY {
				continue
			}
			ev, err := NewKeyInputEvent(events[i])
			if err != nil {
				continue
			}
			log.WithField("event", ev).Debug("Captured")
			inputEvents = append(inputEvents, ev)
		}
		evs <- inputEvents
	}
}
