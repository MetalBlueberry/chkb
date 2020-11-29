package chkb

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	// . "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Mapper", func() {
	DescribeTable("mapOne", func(m *Mapper, event KeyEvent, expected []MapEvent, errExpected bool) {
		obtained, err := m.mapOne(event)
		assert.Equal(GinkgoT(), expected, obtained)
		if errExpected {
			assert.Error(GinkgoT(), err)
		} else {
			assert.NoError(GinkgoT(), err)
		}
	},
		Entry("Forward Up",
			&Mapper{
				Layers: []*Layer{{}},
			},
			KeyEvent{KeyCode: KEY_B, Action: KeyActionDown},
			[]MapEvent{{KeyCode: KEY_B, Action: KbActionDown}},
			false,
		),
		Entry("Forward Down",
			&Mapper{
				Layers: []*Layer{{}},
			},
			KeyEvent{KeyCode: KEY_B, Action: KeyActionUp},
			[]MapEvent{{KeyCode: KEY_B, Action: KbActionUp}},
			false,
		),
		Entry("Block Tap",
			&Mapper{
				Layers: []*Layer{{}},
			},
			KeyEvent{KeyCode: KEY_B, Action: KeyActionTap},
			nil,
			true,
		),
		Entry("map Key up",
			&Mapper{
				Layers: []*Layer{
					{
						KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
					},
				},
			},
			KeyEvent{KeyCode: KEY_A, Action: KeyActionUp},
			[]MapEvent{{KeyCode: KEY_B, Action: KbActionUp}},
			false,
		),
		Entry("map Key down",
			&Mapper{
				Layers: []*Layer{
					{
						KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
					},
				},
			},
			KeyEvent{KeyCode: KEY_A, Action: KeyActionDown},
			[]MapEvent{{KeyCode: KEY_B, Action: KbActionDown}},
			false,
		),
		Entry("no map Key tap",
			&Mapper{
				Layers: []*Layer{
					{
						KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
					},
				},
			},
			KeyEvent{KeyCode: KEY_A, Action: KeyActionTap},
			nil,
			true,
		),
		Entry("Fallback default on down",
			&Mapper{
				Layers: []*Layer{
					{
						OnMiss: []MapEvent{{Action: KbActionPopLayer, LayerName: "test"}},
						KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
					},
				},
			},
			KeyEvent{KeyCode: KEY_C, Action: KeyActionDown},
			[]MapEvent{{Action: KbActionPopLayer, LayerName: "test"}},
			false,
		),
		Entry("no Fallback default on tap",
			&Mapper{
				Layers: []*Layer{
					{
						OnMiss: []MapEvent{{Action: KbActionPopLayer, LayerName: "test"}},
						KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
					},
				},
			},
			KeyEvent{KeyCode: KEY_C, Action: KeyActionTap},
			nil,
			true,
		),
		Entry("Fallback default multi layer",
			&Mapper{
				Layers: []*Layer{
					{
						OnMiss: []MapEvent{{Action: KbActionTap, KeyCode: KEY_A}},
						KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
					},
					{
						OnMiss: []MapEvent{{Action: KbActionTap, KeyCode: KEY_B}},
						KeyMap: KeyMap{KeyCode(KEY_B): {KeyActionMap: {{KeyCode: KEY_C}}}},
					},
				},
			},
			KeyEvent{KeyCode: KEY_A, Action: KeyActionDown},
			[]MapEvent{{Action: KbActionTap, KeyCode: KEY_B}},
			false,
		),
		Entry("Fallback default layer transparent multi action",
			&Mapper{
				Layers: []*Layer{
					{
						KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
					},
					{
						OnMiss: []MapEvent{{Action: KbActionPopLayer, LayerName: "test"}, {Action: KbActionMap}},
						KeyMap: KeyMap{KeyCode(KEY_B): {KeyActionMap: {{KeyCode: KEY_C}}}},
					},
				},
			},
			KeyEvent{KeyCode: KEY_C, Action: KeyActionDown},
			[]MapEvent{{Action: KbActionPopLayer, LayerName: "test"}, {Action: KbActionDown, KeyCode: KEY_C}},
			false,
		),
		Entry("Fallback default multi layer transparent one layer",
			&Mapper{
				Layers: []*Layer{
					{
						KeyMap: KeyMap{KeyCode(KEY_A): {KeyActionMap: {{KeyCode: KEY_B}}}},
					},
					{
						OnMiss: []MapEvent{{Action: KbActionMap}},
						KeyMap: KeyMap{KeyCode(KEY_B): {KeyActionMap: {{KeyCode: KEY_C}}}},
					},
				},
			},
			KeyEvent{KeyCode: KEY_A, Action: KeyActionDown},
			[]MapEvent{{Action: KbActionDown, KeyCode: KEY_B}},
			false,
		),
	)

})
