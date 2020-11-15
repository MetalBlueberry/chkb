package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	evdev "github.com/gvalkov/golang-evdev"
	"github.com/rivo/tview"
)

type Layout struct {
	Keys [][]uint16
}

func NewLayout() *Layout {
	return &Layout{
		Keys: make([][]uint16, 0),
	}
}

func (layout *Layout) AddRow() {
	layout.Keys = append(layout.Keys, make([]uint16, 0))
}

func (layout *Layout) AddKey(key uint16) {
	layout.Keys[layout.lastRow()] = append(layout.Keys[layout.lastRow()], key)
}

func (layout *Layout) lastRow() int {
	return len(layout.Keys) - 1
}

func (layout *Layout) Save(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	return encoder.Encode(&layout.Keys)
}

func (layout *Layout) Load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	layout.Keys = make([][]uint16, 0)
	encoder := json.NewDecoder(f)
	return encoder.Decode(&layout.Keys)
}

type TablePosition struct {
	Row, Col int
}

func (layout *Layout) Map() map[uint16]TablePosition {
	m := make(map[uint16]TablePosition)
	for irow := range layout.Keys {
		for icol := range layout.Keys[irow] {
			code := layout.Keys[irow][icol]
			m[code] = TablePosition{irow, icol}
		}
	}
	return m
}

func (layout *Layout) Record(dev *evdev.InputDevice, file string) error {
	var firstKey uint16

	readKey := func() (*evdev.InputEvent, error) {
		for {
			event, err := dev.ReadOne()
			if err != nil {
				panic(err)
			}
			if event.Type != evdev.EV_KEY {
				continue
			}

			// Panic sequence
			if evdev.KeyEventState(event.Value) == evdev.KeyHold && event.Code == evdev.KEY_F12 {
				panic("Abort program")
			}

			if evdev.KeyEventState(event.Value) == evdev.KeyHold && event.Code == firstKey {
				return event, errors.New("done")
			}
			if evdev.KeyEventState(event.Value) == evdev.KeyUp && event.Code == firstKey {
				log.Print("next row")
				layout.AddRow()
				continue
			}
			if event.Code == firstKey {
				continue
			}
			if evdev.KeyEventState(event.Value) == evdev.KeyUp && firstKey == 0 {
				log.Println("Register next row key as ", kname(event))
				firstKey = event.Code
			}

			if evdev.KeyEventState(event.Value) == evdev.KeyDown {
				return event, nil
			}
		}
	}

	defer dev.Release()
	err := dev.Grab()
	if err != nil {
		return err
	}

	layout.AddRow()

	for {
		event, err := readKey()
		if err != nil {
			log.Println("done")
			break
		}
		log.Printf("event: %s", event.String())
		log.Println(evdev.ByEventType[int(event.Type)][int(event.Code)])

		layout.AddKey(event.Code)
	}

	return layout.Save(file)
}

func (layout *Layout) String() string {
	builder := strings.Builder{}

	for irow := range layout.Keys {
		for icol := range layout.Keys[irow] {
			code := layout.Keys[irow][icol]
			builder.WriteString(evdev.KEY[int(code)])
			builder.WriteRune(' ')
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (layout *Layout) Test(dev *evdev.InputDevice) error {
	app := tview.NewApplication()

	table := tview.NewTable().SetBorders(true)
	table.SetBorder(true)
	table.SetTitle("Layout")
	table.SetSelectable(true, true)

	log := tview.NewTextView()
	log.SetBorder(true)
	log.SetTitle("log")

	hint := tview.NewTextView()
	hint.SetText("Hold F12 to exit")
	hint.SetTextAlign(tview.AlignCenter)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(hint, 1, 10, false)
	flex.AddItem(table, 0, 80, true)
	flex.AddItem(log, 5, 10, false)

	for irow := range layout.Keys {
		for icol := range layout.Keys[irow] {
			code := layout.Keys[irow][icol]
			table.SetCell(irow, icol, tview.NewTableCell(evdev.KEY[int(code)]))
		}
	}

	defer dev.Release()
	err := dev.Grab()
	if err != nil {
		return fmt.Errorf("Cannot graph device, %w", err)
	}

	keyMap := layout.Map()

	go func() {
		for {
			app.Draw()
			event, err := dev.ReadOne()
			if err != nil {
				panic(err)
			}

			if event.Type != evdev.EV_KEY {
				continue
			}

			if evdev.KeyEventState(event.Value) == evdev.KeyHold && event.Code == evdev.KEY_F12 {
				app.Stop()
				return
			}

			key, ok := keyMap[event.Code]
			if !ok {
				log.SetText("key not found " + evdev.KEY[int(event.Code)])
				continue
			}
			switch evdev.KeyEventState(event.Value) {
			case evdev.KeyDown:
				log.SetText("key down " + evdev.KEY[int(event.Code)])
				c := table.GetCell(key.Row, key.Col).SetBackgroundColor(tcell.ColorRed)
				table.SetCell(key.Row, key.Col, c)
				table.Select(key.Row, key.Col)
			case evdev.KeyUp:
				log.SetText("key up " + evdev.KEY[int(event.Code)])
				c := table.GetCell(key.Row, key.Col).SetBackgroundColor(tcell.ColorGrey)
				table.SetCell(key.Row, key.Col, c)
				table.Select(key.Row, key.Col)
			}

		}
	}()

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		return fmt.Errorf("app finished unexpectedly, %w", err)
	}
	return nil
}
