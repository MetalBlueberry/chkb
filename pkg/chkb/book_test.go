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
    keyMap:
        KEY_LEFTSHIFT:
            Tap:
              - action: PushLayer
                layerName: swapAB
swapAB:
    keyMap:
        KEY_A:
            Map:
              - keyCode: KEY_B
        KEY_B:
            Map:
              - keyCode: KEY_A
        KEY_LEFTSHIFT:
            Tap:
              - action: PopLayer
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
				To(HaveKeyWithValue(chkb.KeyActionTap, []chkb.MapEvent{
					{
						Action:    chkb.KbActionPushLayer,
						LayerName: "swapAB",
					},
				}))
			Expect(book["swapAB"].KeyMap).To(HaveKey(chkb.KeyCode(evdev.KEY_A)))
			Expect(book["swapAB"].KeyMap[chkb.KeyCode(evdev.KEY_A)]).
				To(HaveKeyWithValue(chkb.KeyActionMap, []chkb.MapEvent{
					{
						KeyCode: evdev.KEY_B,
					},
				}))
		})
		It("Save simple", func() {
			book := chkb.Book{
				"base": {
					KeyMap: map[chkb.KeyCode]map[chkb.KeyActions][]chkb.MapEvent{
						evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPushLayer, LayerName: "swapAB"}}},
					},
				},
				"swapAB": {
					KeyMap: map[chkb.KeyCode]map[chkb.KeyActions][]chkb.MapEvent{
						evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPopLayer}}},
						evdev.KEY_A:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_B}}},
						evdev.KEY_B:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_A}}},
					},
				},
			}
			buf := &bytes.Buffer{}
			book.Save(buf)

			Expect(buf.String()).To(Equal(fileContent))
		})
	})

})
