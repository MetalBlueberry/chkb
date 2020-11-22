package chkb_test

import (
	"bytes"
	"strings"

	evdev "github.com/gvalkov/golang-evdev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"MetalBlueberry/cheap-keyboard/pkg/chkb"
)

var _ = Describe("Book", func() {

	Describe("Load", func() {
		var (
			fileContent = `base:
    KEY_LEFTSHIFT:
        Tap:
            action: Push
            layerName: swapAB
swapAB:
    KEY_A:
        Map:
            keyCode: KEY_B
    KEY_B:
        Map:
            keyCode: KEY_A
    KEY_LEFTSHIFT:
        Tap:
            action: Pop
`
		)
		It("Load simple", func() {
			book := chkb.Book{}
			err := book.Load(strings.NewReader(fileContent))
			Expect(err).To(BeNil())
			Expect(book).To(HaveKey("base"))
			Expect(book).To(HaveKey("swapAB"))

			Expect(book["base"].KeyMap).To(HaveKey(chkb.KeyCode(evdev.KEY_LEFTSHIFT)))
			Expect(book["base"].KeyMap[chkb.KeyCode(evdev.KEY_LEFTSHIFT)]).
				To(HaveKeyWithValue(chkb.ActionTap, chkb.MapEvent{
					Action:    chkb.ActionPushLayer,
					LayerName: "swapAB",
				}))
			Expect(book["swapAB"].KeyMap).To(HaveKey(chkb.KeyCode(evdev.KEY_A)))
			Expect(book["swapAB"].KeyMap[chkb.KeyCode(evdev.KEY_A)]).
				To(HaveKeyWithValue(chkb.ActionMap, chkb.MapEvent{
					KeyCode: evdev.KEY_B,
				}))
		})
		It("Save simple", func() {
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
			buf := &bytes.Buffer{}
			book.Save(buf)

			Expect(buf.String()).To(Equal(fileContent))
		})
	})

})
