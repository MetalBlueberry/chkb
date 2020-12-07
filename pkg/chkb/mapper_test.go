package chkb_test

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

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

	FDescribeTable("Map", func(m *chkb.Mapper, events []chkb.KeyEvent, expected []chkb.MapEvent, errExpected bool) {
		obtained, err := m.Maps(events)
		assert.Equal(GinkgoT(), expected, obtained)
		if errExpected {
			assert.Error(GinkgoT(), err)
		} else {
			assert.NoError(GinkgoT(), err)
		}
	},
		Entry("Forward key",
			chkb.NewMapper().WithLayers(chkb.Layers{
				&chkb.Layer{},
			}),
			[]chkb.KeyEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KeyActionUp},
			},
			[]chkb.MapEvent{
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionDown},
				{KeyCode: chkb.KEY_B, Action: chkb.KbActionUp},
			},
			false,
		),
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
			},
			[]chkb.MapEvent{
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
		Entry("SmartTap Tap",
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
	)

})
