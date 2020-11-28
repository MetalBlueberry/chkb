package main

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"
	"MetalBlueberry/cheap-keyboard/pkg/deliverers/vkb"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/bendahl/uinput"
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

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
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
	)
	kb.AddDeliverer(&vkb.Keyboard{keyboard})

	go FileUpdate(kb)

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

		inputEvents := make([]chkb.InputEvent, 0)
		for i := range events {
			if events[i].Type != evdev.EV_KEY {
				continue
			}
			inputEvents = append(inputEvents, chkb.NewKeyInputEvent(events[i]))
		}

		captured, err := kb.Capture(inputEvents)
		if err != nil {
			panic(err)
		}
		maps, err := kb.Maps(captured)
		if err != nil {
			panic(err)
		}
		err = kb.Delivers(maps)
		if err != nil {
			panic(err)
		}
	}
}

func LocalServer(kb *chkb.Keyboard) {
	http.Handle("/status", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, layerString(kb))
	}))
	log.Print("status server started")
	http.ListenAndServe(":9989", nil)
}

func FileUpdate(kb *chkb.Keyboard) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fileName := ".chkb_layout"
	file, err := os.Create(filepath.Join(usr.HomeDir, fileName))
	if err != nil {
		log.Printf("error file update, %s", err)
		return
	}
	defer file.Close()
	ticker := time.NewTicker(time.Millisecond * 100)

	previousString := ""
	for {
		<-ticker.C
		str := layerString(kb)
		if previousString == str {
			continue
		}
		previousString = str
		fmt.Fprint(file, str)
	}
}

func layerString(kb *chkb.Keyboard) string {
	builder := strings.Builder{}
	for i := range kb.Layers {
		for name := range kb.LayerBook {
			if kb.Layers[i] == kb.LayerBook[name] {
				builder.WriteString(name)
				builder.WriteString(" > ")
			}
		}
	}
	builder.WriteString("\n")
	return builder.String()
}
