package chkb_test

import (
	evdev "github.com/gvalkov/golang-evdev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/stretchr/testify/assert"

	// . "github.com/onsi/gomega"

	"MetalBlueberry/cheap-keyboard/pkg/chkb"
)

var _ = Describe("Captor", func() {

	var (
		captor *chkb.Captor
	)
	BeforeEach(func() {
		captor = chkb.NewCaptor()
	})

	DescribeTable("Capture",
		func(events []chkb.InputEvent, expected []chkb.KeyEvent) {
			captured := make([]chkb.KeyEvent, 0)
			for _, event := range events {
				c, err := captor.CaptureOne(event)
				assert.NoError(GinkgoT(), err)
				captured = append(captured, c...)
			}
			assert.Equal(GinkgoT(), expected, captured)
		},
		Entry("KeyDown", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
		}),
		Entry("KeyUp", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
		}),
		Entry("KeyHold", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(100), KeyCode: evdev.KEY_A, Action: chkb.InputActionHold},
			{Time: Elapsed(200), KeyCode: evdev.KEY_A, Action: chkb.InputActionHold},
			{Time: Elapsed(300), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionHold},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
		}),
		Entry("Tap", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(50), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionTap},
		}),
		Entry("NotTap - slow", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(250), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
		}),
	)

})
