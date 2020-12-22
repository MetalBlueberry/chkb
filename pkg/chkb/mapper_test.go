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
package chkb_test

import (
	"github.com/MetalBlueberry/chkb/pkg/chkb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	log "github.com/sirupsen/logrus"

	// . "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Mapper", func() {
	DescribeTable("Map",
		func(m *chkb.Mapper, events []chkb.KeyEvent, expected []chkb.MapEvent, errExpected bool) {
			obtained := []chkb.MapEvent{}
			err := m.Maps(events, func(event chkb.MapEvent) error {
				obtained = append(obtained, event)
				return nil
			})
			for i := range obtained {
				if len(expected) > i {
					log.Printf("%s - %s", obtained[i], expected[i])
				} else {
					log.Printf("%s", obtained[i])
				}
			}
			assert.Equal(GinkgoT(), expected, obtained)
			if errExpected {
				assert.Error(GinkgoT(), err)
			} else {
				assert.NoError(GinkgoT(), err)
			}
		},
		Entry("Map key",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_B: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{Action: chkb.KbActionMap, KeyCode: chkb.KEY_A},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionTap},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Forward real key if not mapped",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionTap},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Forward if keymap only contains map mapped",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_B: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{Action: chkb.KbActionMap, KeyCode: chkb.KEY_C},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionTap},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Forward if not has special",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_LEFTSHIFT: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionDown: {
								{Action: chkb.KbActionPushLayer, LayerName: "test"},
							},
							chkb.KeyActionUp: {
								{Action: chkb.KbActionPopLayer, LayerName: "test"},
							},
							chkb.KeyActionMap: {
								{Action: chkb.KbActionMap, KeyCode: chkb.KEY_LEFTSHIFT},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionTap},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{Action: chkb.KbActionPushLayer, LayerName: "test"},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{Action: chkb.KbActionPopLayer, LayerName: "test"},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Hold keys if has special",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_B: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{Action: chkb.KbActionNil},
							},
							chkb.KeyActionHold: {
								{Action: chkb.KbActionPushLayer, LayerName: "test"},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionTap},
			},
			[]chkb.MapEvent{
				{Action: chkb.KbActionPushLayer, LayerName: "test"},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Ignore special",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_B: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{Action: chkb.KbActionMap, KeyCode: chkb.KEY_A},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionTap},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionUp},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Smart Tap-Tap",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_B: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{Action: chkb.KbActionMap, KeyCode: chkb.KEY_A},
							},
							chkb.KeyActionTap: {
								{Action: chkb.KbActionMap, KeyCode: chkb.KEY_B},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionTap},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionUp},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionTap},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Smart Tap-Hold",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_B: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{Action: chkb.KbActionMap, KeyCode: chkb.KEY_A},
							},
							chkb.KeyActionHold: {
								{Action: chkb.KbActionDown, KeyCode: chkb.KEY_LEFTSHIFT},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionTap},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Smart Hold-Hold-Tap",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_A: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionHold: {
								{Action: chkb.KbActionDown, KeyCode: chkb.KEY_LEFTSHIFT},
							},
						},
						chkb.KEY_B: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionHold: {
								{Action: chkb.KbActionDown, KeyCode: chkb.KEY_LEFTCTRL},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionTap},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_LEFTCTRL, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_LEFTCTRL, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Smart Hold-Hold-Tap Block map",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_A: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{Action: chkb.KbActionNil},
							},
							chkb.KeyActionTap: {
								{Action: chkb.KbActionTap, KeyCode: chkb.KEY_A},
							},
							chkb.KeyActionHold: {
								{Action: chkb.KbActionDown, KeyCode: chkb.KEY_LEFTSHIFT},
							},
						},
						chkb.KEY_B: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{Action: chkb.KbActionNil},
							},
							chkb.KeyActionTap: {
								{Action: chkb.KbActionTap, KeyCode: chkb.KEY_B},
							},
							chkb.KeyActionHold: {
								{Action: chkb.KbActionDown, KeyCode: chkb.KEY_LEFTCTRL},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				//HOLD
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionTap},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},

				//TAP
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionTap},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_LEFTCTRL, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_LEFTCTRL, Action: chkb.KbActionUp},

				{KeyCode: chkb.KEY_LEFTCTRL, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionTap},
				{KeyCode: chkb.KEY_LEFTCTRL, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("onMiss do not trigger when key is mapped",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					OnMiss: []chkb.MapEvent{
						{Action: chkb.KbActionPopLayer, LayerName: "test"},
					},
					KeyMap: chkb.KeyMap{
						chkb.KEY_A: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{KeyCode: chkb.KEY_B},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionTap},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("onMiss do not trigger when key contains actions",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					OnMiss: []chkb.MapEvent{
						{Action: chkb.KbActionPopLayer, LayerName: "test"},
					},
					KeyMap: chkb.KeyMap{
						chkb.KEY_A: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionTap: {
								{Action: chkb.KbActionChangeLayer, LayerName: "other"},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionTap},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
				{Action: chkb.KbActionChangeLayer, LayerName: "other"},
			},
			false,
		),
		Entry("onMiss trigger when key is not mapped",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					OnMiss: []chkb.MapEvent{
						{Action: chkb.KbActionPopLayer, LayerName: "test"},
					},
					KeyMap: chkb.KeyMap{
						chkb.KEY_A: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{KeyCode: chkb.KEY_B},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionTap},
			},
			[]chkb.MapEvent{
				{Action: chkb.KbActionPopLayer, LayerName: "test"},
			},
			false,
		),
		Entry("onMiss forward keys if there is a map action",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					OnMiss: []chkb.MapEvent{
						{Action: chkb.KbActionPopLayer, LayerName: "test"},
						{Action: chkb.KbActionMap},
					},
					KeyMap: chkb.KeyMap{
						chkb.KEY_A: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionMap: {
								{KeyCode: chkb.KEY_B},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_C, Action: chkb.KeyActionTap},
			},
			[]chkb.MapEvent{
				{Action: chkb.KbActionPopLayer, LayerName: "test"},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionUp},
			},
			false,
		),
	)

})
