package main

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

func LocalServer(kb *chkb.Keyboard) {
	http.Handle("/status", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, layerString(kb))
	}))
	log.Print("status server started")
	http.ListenAndServe(":9989", nil)
}

func FileUpdate(kb *chkb.Keyboard) {
	os.MkdirAll("/tmp/chkb", 0777)
	file, err := os.Create("/tmp/chkb/layout.tmp")
	if err != nil {
		log.Printf("error file update, %s", err)
		return
	}
	defer file.Close()
	defer os.Remove("/tmp/chkb/layout.tmp")
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
