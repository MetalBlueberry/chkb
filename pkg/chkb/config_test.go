/*
Copyright © 2020 Víctor Pérez @MetalBlueberry

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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

var _ = Describe("Config", func() {

	Describe("Load", func() {
		var (
			fileContent = `layers:
    base:
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
			book := Config{}
			err := book.Load(strings.NewReader(fileContent))
			Expect(err).To(BeNil())
			Expect(book.Layers).To(HaveKey("base"))
			Expect(book.Layers).To(HaveKey("swapAB"))

			Expect(book.Layers["base"].KeyMap).To(HaveKey(KeyCode(evdev.KEY_LEFTSHIFT)))
			Expect(book.Layers["base"].KeyMap[KeyCode(evdev.KEY_LEFTSHIFT)]).
				To(HaveKeyWithValue(KeyActionTap, []MapEvent{
					{
						Action:    KbActionPushLayer,
						LayerName: "swapAB",
					},
				}))
			Expect(book.Layers["swapAB"].KeyMap).To(HaveKey(KeyCode(evdev.KEY_A)))
			Expect(book.Layers["swapAB"].KeyMap[KeyCode(evdev.KEY_A)]).
				To(HaveKeyWithValue(KeyActionMap, []MapEvent{
					{
						KeyCode: evdev.KEY_B,
					},
				}))
		})
		It("Save simple", func() {
			book := Config{
				Layers: LayerBook{
					"base": {
						KeyMap: map[KeyCode]KeyMapActions{
							evdev.KEY_LEFTSHIFT: {KeyActionTap: {{Action: KbActionPushLayer, LayerName: "swapAB"}}},
						},
					},
					"swapAB": {
						KeyMap: map[KeyCode]KeyMapActions{
							evdev.KEY_LEFTSHIFT: {KeyActionTap: {{Action: KbActionPopLayer}}},
							evdev.KEY_A:         {KeyActionMap: {{KeyCode: evdev.KEY_B}}},
							evdev.KEY_B:         {KeyActionMap: {{KeyCode: evdev.KEY_A}}},
						},
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
