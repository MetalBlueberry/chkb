package chkb_test

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	log "github.com/sirupsen/logrus"

	// . "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Mapper", func() {
	DescribeTable("mapOne", func(m *chkb.Mapper, event chkb.KeyEvent, expected []chkb.MapEvent, errExpected bool) {
		obtained, err := m.MapOne(event)
		assert.Equal(GinkgoT(), expected, obtained)
		if errExpected {
			assert.Error(GinkgoT(), err)
		} else {
			assert.NoError(GinkgoT(), err)
		}
	},
		Entry("Forward Up",
			&chkb.Mapper{
				Layers: []*chkb.Layer{{}},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
			[]chkb.MapEvent{{KeyCode: chkb.KEY_B, Action: chkb.KbActionDown}},
			false,
		),
		Entry("Forward Down",
			&chkb.Mapper{
				Layers: []*chkb.Layer{{}},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
			[]chkb.MapEvent{{KeyCode: chkb.KEY_B, Action: chkb.KbActionUp}},
			false,
		),
		Entry("Block Tap",
			&chkb.Mapper{
				Layers: []*chkb.Layer{{}},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_B, Action: chkb.KeyActionTap},
			nil,
			true,
		),
		Entry("map Key up",
			&chkb.Mapper{
				Layers: []*chkb.Layer{
					{
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_A): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_B}}}},
					},
				},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
			[]chkb.MapEvent{{KeyCode: chkb.KEY_B, Action: chkb.KbActionUp}},
			false,
		),
		Entry("map Key down",
			&chkb.Mapper{
				Layers: []*chkb.Layer{
					{
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_A): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_B}}}},
					},
				},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
			[]chkb.MapEvent{{KeyCode: chkb.KEY_B, Action: chkb.KbActionDown}},
			false,
		),
		Entry("no map Key tap",
			&chkb.Mapper{
				Layers: []*chkb.Layer{
					{
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_A): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_B}}}},
					},
				},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_A, Action: chkb.KeyActionTap},
			nil,
			true,
		),
		Entry("Fallback default on down",
			&chkb.Mapper{
				Layers: []*chkb.Layer{
					{
						OnMiss: []chkb.MapEvent{{Action: chkb.KbActionPopLayer, LayerName: "test"}},
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_A): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_B}}}},
					},
				},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_C, Action: chkb.KeyActionDown},
			[]chkb.MapEvent{{Action: chkb.KbActionPopLayer, LayerName: "test"}},
			false,
		),
		Entry("no Fallback default on tap",
			&chkb.Mapper{
				Layers: []*chkb.Layer{
					{
						OnMiss: []chkb.MapEvent{{Action: chkb.KbActionPopLayer, LayerName: "test"}},
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_A): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_B}}}},
					},
				},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_C, Action: chkb.KeyActionTap},
			nil,
			true,
		),
		Entry("Fallback default multi layer",
			&chkb.Mapper{
				Layers: []*chkb.Layer{
					{
						OnMiss: []chkb.MapEvent{{Action: chkb.KbActionTap, KeyCode: chkb.KEY_A}},
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_A): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_B}}}},
					},
					{
						OnMiss: []chkb.MapEvent{{Action: chkb.KbActionTap, KeyCode: chkb.KEY_B}},
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_B): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_C}}}},
					},
				},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
			[]chkb.MapEvent{{Action: chkb.KbActionTap, KeyCode: chkb.KEY_B}},
			false,
		),
		Entry("Fallback default layer transparent multi action",
			&chkb.Mapper{
				Layers: []*chkb.Layer{
					{
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_A): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_B}}}},
					},
					{
						OnMiss: []chkb.MapEvent{{Action: chkb.KbActionPopLayer, LayerName: "test"}, {Action: chkb.KbActionMap}},
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_B): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_C}}}},
					},
				},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_C, Action: chkb.KeyActionDown},
			[]chkb.MapEvent{{Action: chkb.KbActionPopLayer, LayerName: "test"}, {Action: chkb.KbActionDown, KeyCode: chkb.KEY_C}},
			false,
		),
		Entry("Fallback default multi layer transparent one layer",
			&chkb.Mapper{
				Layers: []*chkb.Layer{
					{
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_A): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_B}}}},
					},
					{
						OnMiss: []chkb.MapEvent{{Action: chkb.KbActionMap}},
						KeyMap: chkb.KeyMap{chkb.KeyCode(chkb.KEY_B): {chkb.KeyActionMap: {{KeyCode: chkb.KEY_C}}}},
					},
				},
			},
			chkb.KeyEvent{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
			[]chkb.MapEvent{{Action: chkb.KbActionDown, KeyCode: chkb.KEY_B}},
			false,
		),
	)

	DescribeTable("Map", func(m *chkb.Mapper, events []chkb.KeyEvent, expected []chkb.MapEvent, errExpected bool) {
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
		Entry("Forward Real key if not mapped",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionNil},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionNil},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("forward if mapped",
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
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionNil},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionNil},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_C, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Not hold if is forwarding",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{
					KeyMap: chkb.KeyMap{
						chkb.KEY_LEFTSHIFT: map[chkb.KeyActions][]chkb.MapEvent{
							chkb.KeyActionTap: {
								{Action: chkb.KbActionPushLayer, LayerName: "test"},
							},
						},
					},
				},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KeyActionNil},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KeyActionNil},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("hold keys if maps to nil",
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
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionHold},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionTap},
				{KeyCode: chkb.KEY_A, Action: chkb.KeyActionUp},
			},
			[]chkb.MapEvent{
				{Action: chkb.KbActionPushLayer, LayerName: "test"},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_A, Action: chkb.KbActionUp},
			},
			false,
		),
		Entry("Ignore Special",
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
		Entry("SmartTap/Tap",
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
		Entry("SmartTap/Hold",
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
		Entry("Smart/Hold-Hold-Tap",
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
		Entry("Smart/Hold-Hold-Tap Block map",
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
	)

})
