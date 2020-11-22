package main

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"
	"log"
	"os"

	evdev "github.com/gvalkov/golang-evdev"
)

func main() {
	book := chkb.Book{
		"base": {
			KeyMap: map[chkb.KeyCode]map[chkb.Actions]chkb.MapEvent{
				evdev.KEY_LEFTSHIFT: {chkb.ActionTap: {Action: chkb.ActionPushLayer, LayerName: "swapAB"}},
			},
		},
		"swapAB": {
			KeyMap: map[chkb.KeyCode]map[chkb.Actions]chkb.MapEvent{
				evdev.KEY_LEFTSHIFT: {chkb.ActionTap: {Action: chkb.ActionPopLayer}},
				evdev.KEY_A:         {chkb.ActionMap: {KeyCode: evdev.KEY_B}},
				evdev.KEY_B:         {chkb.ActionMap: {KeyCode: evdev.KEY_A}},
			},
		},
	}

	err := book.Save(os.Stdout)
	if err != nil {
		log.Println(err)
	}
}
