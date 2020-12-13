package main

import (
	"log"
	"os"

	"github.com/MetalBlueberry/chkb/pkg/chkb"

	evdev "github.com/gvalkov/golang-evdev"
)

func main() {
	config := chkb.Config{
		Layers: map[string]*chkb.Layer{
			"base": {
				KeyMap: chkb.KeyMap{
					evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPushLayer, LayerName: "swapAB"}}},
				},
			},
			"swapAB": {
				KeyMap: chkb.KeyMap{
					evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPopLayer}}},
					evdev.KEY_A:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_B}}},
					evdev.KEY_B:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_A}}},
				},
			},
		},
	}

	err := config.Save(os.Stdout)
	if err != nil {
		log.Println(err)
	}
}
