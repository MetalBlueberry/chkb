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
		func(events []evdev.InputEvent, expected []chkb.KeyEvent) {
			captured := make([]chkb.KeyEvent, 0)
			for _, event := range events {
				c, err := captor.CaptureOne(chkb.NewKeyInputEvent(event))
				assert.NoError(GinkgoT(), err)
				captured = append(captured, c...)
			}
			assert.Equal(GinkgoT(), expected, captured)
		},
		Entry("KeyDown", []evdev.InputEvent{
			{Time: Elapsed(0), Code: evdev.KEY_A, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
		}),
		Entry("KeyUp", []evdev.InputEvent{
			{Time: Elapsed(0), Code: evdev.KEY_A, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
		}),
		Entry("KeyHold", []evdev.InputEvent{
			{Time: Elapsed(0), Code: evdev.KEY_A, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
			{Time: Elapsed(100), Code: evdev.KEY_A, Value: int32(evdev.KeyHold), Type: evdev.EV_KEY},
			{Time: Elapsed(200), Code: evdev.KEY_A, Value: int32(evdev.KeyHold), Type: evdev.EV_KEY},
			{Time: Elapsed(300), Code: evdev.KEY_A, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionHold},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
		}),
		Entry("Tap", []evdev.InputEvent{
			{Time: Elapsed(0), Code: evdev.KEY_A, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
			{Time: Elapsed(50), Code: evdev.KEY_A, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionTap},
		}),
		Entry("NotTap - slow", []evdev.InputEvent{
			{Time: Elapsed(0), Code: evdev.KEY_A, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
			{Time: Elapsed(250), Code: evdev.KEY_A, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
		}, []chkb.KeyEvent{
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
		}),
	)

})
