package main

import (
	"errors"
	"log"

	evdev "github.com/gvalkov/golang-evdev"
)

type Input struct {
	*evdev.InputDevice

	events chan []evdev.InputEvent
	grab   bool
	grabed bool
}

func NewInput(dev *evdev.InputDevice) *Input {
	i := &Input{
		InputDevice: dev,
		events:      make(chan []evdev.InputEvent),
	}
	go i.Capture()
	return i
}

func (input *Input) Capture() {
	for {
		if input.grab != input.grabed {
			if input.grab {
				err := input.Grab()
				if err != nil {
					panic(err)
				}
				log.Println("grabed")
			} else {
				err := input.Release()
				if err != nil {
					panic(err)
				}
				log.Println("released")
			}
			input.grabed = input.grab
		}

		evnts, err := input.Read()
		if err != nil {
			log.Println("cannot read device, %w", err)
			close(input.events)
			return
		}

		select {
		case input.events <- evnts:
		default:
			log.Println("skip key")
		}
	}
}

func (input *Input) Chan() <-chan []evdev.InputEvent {
	return input.events
}

func (input *Input) ReadOneKeyDown() (evdev.InputEvent, error) {
	input.grab = true
	defer func() { input.grab = false }()

	for events := range input.events {
		for _, event := range events {
			if event.Type != evdev.EV_KEY {
				continue
			}
			if evdev.KeyEventState(event.Value) == evdev.KeyDown {
				return event, nil
			}
		}
	}
	return evdev.InputEvent{}, errors.New("event channel is closed")
}
