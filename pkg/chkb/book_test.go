package chkb

import (
	"bytes"
	"strings"

	evdev "github.com/gvalkov/golang-evdev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
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
			book := Book{}
			err := book.Load(strings.NewReader(fileContent))
			Expect(err).To(BeNil())
			Expect(book).To(HaveKey("base"))
			Expect(book).To(HaveKey("swapAB"))

			Expect(book["base"].KeyMap).To(HaveKey(KeyCode(evdev.KEY_LEFTSHIFT)))
			Expect(book["base"].KeyMap[KeyCode(evdev.KEY_LEFTSHIFT)]).
				To(HaveKeyWithValue(KeyActionTap, []MapEvent{
					{
						Action:    KbActionPushLayer,
						LayerName: "swapAB",
					},
				}))
			Expect(book["swapAB"].KeyMap).To(HaveKey(KeyCode(evdev.KEY_A)))
			Expect(book["swapAB"].KeyMap[KeyCode(evdev.KEY_A)]).
				To(HaveKeyWithValue(KeyActionMap, []MapEvent{
					{
						KeyCode: evdev.KEY_B,
					},
				}))
		})
		It("Save simple", func() {
			book := Book{
				"base": {
					KeyMap: map[KeyCode]map[KeyActions][]MapEvent{
						evdev.KEY_LEFTSHIFT: {KeyActionTap: {{Action: KbActionPushLayer, LayerName: "swapAB"}}},
					},
				},
				"swapAB": {
					KeyMap: map[KeyCode]map[KeyActions][]MapEvent{
						evdev.KEY_LEFTSHIFT: {KeyActionTap: {{Action: KbActionPopLayer}}},
						evdev.KEY_A:         {KeyActionMap: {{KeyCode: evdev.KEY_B}}},
						evdev.KEY_B:         {KeyActionMap: {{KeyCode: evdev.KEY_A}}},
					},
				},
			}
			buf := &bytes.Buffer{}
			book.Save(buf)

			Expect(buf.String()).To(Equal(fileContent))
		})
	})

	DescribeTable("findMap", func(layer *Layer, event KeyEvent, expected []MapEvent, foundExpected bool) {
		obtained, found := layer.findMap(event)
		assert.Equal(GinkgoT(), foundExpected, found)
		assert.Equal(GinkgoT(), expected, obtained)
	},
		Entry("Not found",
			&Layer{
				KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
			},
			KeyEvent{KeyCode: KEY_B, Action: KeyActionDown},
			nil,
			false,
		),
		Entry("MapA keyDown",
			&Layer{
				KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
			},
			KeyEvent{KeyCode: KEY_A, Action: KeyActionDown},
			[]MapEvent{{KeyCode: KEY_B, Action: KbActionDown}},
			true,
		),
		Entry("MapA keyUp",
			&Layer{
				KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
			},
			KeyEvent{KeyCode: KEY_A, Action: KeyActionUp},
			[]MapEvent{{KeyCode: KEY_B, Action: KbActionUp}},
			true,
		),
		Entry("ActionDown",
			&Layer{
				KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionDown: {{KeyCode: KEY_B}}}},
			},
			KeyEvent{KeyCode: KEY_A, Action: KeyActionDown},
			[]MapEvent{{KeyCode: KEY_B, Action: KbActionDown}},
			true,
		),
		Entry("ActionUp",
			&Layer{
				KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionUp: {{KeyCode: KEY_B}}}},
			},
			KeyEvent{KeyCode: KEY_A, Action: KeyActionUp},
			[]MapEvent{{KeyCode: KEY_B, Action: KbActionUp}},
			true,
		),
	)
})
