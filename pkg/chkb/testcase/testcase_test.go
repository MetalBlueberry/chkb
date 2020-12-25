package testcase_test

import (
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
	"github.com/MetalBlueberry/chkb/pkg/chkb"
	"github.com/MetalBlueberry/chkb/pkg/chkb/testcase"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Testcase", func() {
	DescribeTable("Map",
		func(seq string, expected []chkb.InputEvent, errExpected bool) {
			events, err := testcase.Input(seq)
			if errExpected {
				assert.Error(GinkgoT(), err)
			} else {
				assert.NoError(GinkgoT(), err)
				assert.Equal(GinkgoT(), expected, events)
			}
		},
		Entry("Tap A",
			"Ta",
			[]chkb.InputEvent{
				{Time: testcase.Step(0), KeyCode: chkb.KEY_A, Action: chkb.InputActionDown},
				{Time: testcase.Step(1), KeyCode: chkb.KEY_A, Action: chkb.InputActionUp},
			},
			nil,
		),
		Entry("Tap B",
			"Tb",
			[]chkb.InputEvent{
				{Time: testcase.Step(0), KeyCode: chkb.KEY_B, Action: chkb.InputActionDown},
				{Time: testcase.Step(1), KeyCode: chkb.KEY_B, Action: chkb.InputActionUp},
			},
			nil,
		),
		Entry("Tap AB",
			"Ta Tb",
			[]chkb.InputEvent{
				{Time: testcase.Step(0), KeyCode: chkb.KEY_A, Action: chkb.InputActionDown},
				{Time: testcase.Step(1), KeyCode: chkb.KEY_A, Action: chkb.InputActionUp},
				{Time: testcase.Step(2), KeyCode: chkb.KEY_B, Action: chkb.InputActionDown},
				{Time: testcase.Step(3), KeyCode: chkb.KEY_B, Action: chkb.InputActionUp},
			},
			nil,
		),
		Entry("Press A",
			"Pa",
			[]chkb.InputEvent{
				{Time: testcase.Step(0), KeyCode: chkb.KEY_A, Action: chkb.InputActionDown},
			},
			nil,
		),
		Entry("Release A",
			"Ra",
			[]chkb.InputEvent{
				{Time: testcase.Step(0), KeyCode: chkb.KEY_A, Action: chkb.InputActionUp},
			},
			nil,
		),
		Entry("Press B",
			"Pb",
			[]chkb.InputEvent{
				{Time: testcase.Step(0), KeyCode: chkb.KEY_B, Action: chkb.InputActionDown},
			},
			nil,
		),
		Entry("Release B",
			"Rb",
			[]chkb.InputEvent{
				{Time: testcase.Step(0), KeyCode: chkb.KEY_B, Action: chkb.InputActionUp},
			},
			nil,
		),
		Entry("Type ABCD",
			"Ta Tb Tc Td",
			[]chkb.InputEvent{
				{Time: testcase.Step(0), KeyCode: chkb.KEY_A, Action: chkb.InputActionDown},
				{Time: testcase.Step(1), KeyCode: chkb.KEY_A, Action: chkb.InputActionUp},
				{Time: testcase.Step(2), KeyCode: chkb.KEY_B, Action: chkb.InputActionDown},
				{Time: testcase.Step(3), KeyCode: chkb.KEY_B, Action: chkb.InputActionUp},
				{Time: testcase.Step(4), KeyCode: chkb.KEY_C, Action: chkb.InputActionDown},
				{Time: testcase.Step(5), KeyCode: chkb.KEY_C, Action: chkb.InputActionUp},
				{Time: testcase.Step(6), KeyCode: chkb.KEY_D, Action: chkb.InputActionDown},
				{Time: testcase.Step(7), KeyCode: chkb.KEY_D, Action: chkb.InputActionUp},
			},
			nil,
		),
		Entry("Hold delay",
			"Pa 200 Ra",
			[]chkb.InputEvent{
				{Time: testcase.Step(0), KeyCode: chkb.KEY_A, Action: chkb.InputActionDown},
				{Time: testcase.Elapsed(1, 200), KeyCode: chkb.KEY_A, Action: chkb.InputActionUp},
			},
			nil,
		),

		Entry("Invalid, No space",
			"Pa200Ra",
			[]chkb.InputEvent{},
			true,
		),
		Entry("Invalid, lowercase",
			"pa",
			[]chkb.InputEvent{},
			true,
		),
	)
})
