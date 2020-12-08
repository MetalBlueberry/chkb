package chkb_test

import (
	"errors"

	"github.com/benbjohnson/clock"
	evdev "github.com/gvalkov/golang-evdev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/stretchr/testify/assert"

	// . "github.com/onsi/gomega"

	"MetalBlueberry/cheap-keyboard/pkg/chkb"
)

var _ = Describe("Captor", func() {

	var (
		clockMock *clock.Mock
		captor    *chkb.Captor
	)
	BeforeEach(func() {
		clockMock = clock.NewMock()
		captor = chkb.NewCaptorWithClock(clockMock)
	})

	DescribeTable("Capture",
		func(events []chkb.InputEvent, expected []chkb.KeyEvent) {
			clockMock.Set(InitialTime)

			allcaptured := make([]chkb.KeyEvent, 0)
			i := 0
			captor.Run(func() ([]chkb.InputEvent, error) {
				if len(events) == i {
					clockMock.Add(chkb.TapDelay)
					return nil, errors.New("Finished")
				}
				event := events[i]
				if clockMock.Now().After(event.Time) {
					Fail("events do not follow the timeline")
				}
				clockMock.Set(event.Time)
				i++
				return []chkb.InputEvent{event}, nil
			}, func(ke []chkb.KeyEvent) error {
				allcaptured = append(allcaptured, ke...)
				return nil
			})
			assert.Equal(GinkgoT(), expected, allcaptured)
		},
		Entry("KeyHold", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(AfterTap), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{Time: Elapsed(TapDelayMs), KeyCode: evdev.KEY_A, Action: chkb.KeyActionHold},
			{Time: Elapsed(AfterTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
		}),
		Entry("Tap", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{Time: Elapsed(BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			{Time: Elapsed(BeforeTap + TapDelayMs), KeyCode: evdev.KEY_A, Action: chkb.KeyActionTap},
		}),
		Entry("Tap quick. Tap event must be fired as soon as other key is pressed", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
			{Time: Elapsed(10 + BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
			{Time: Elapsed(10 + 2*BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{Time: Elapsed(BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			{Time: Elapsed(10 + BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionTap},
			{Time: Elapsed(10 + BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.KeyActionDown},
			{Time: Elapsed(10 + 2*BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.KeyActionUp},
			{Time: Elapsed(10 + 2*BeforeTap + TapDelayMs), KeyCode: evdev.KEY_B, Action: chkb.KeyActionTap},
		}),
		Entry("Hold quick. Hold event must be fired after TapDelay, even if other key is pressed", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(10), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
			{Time: Elapsed(10 + BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			{Time: Elapsed(AfterTap), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{Time: Elapsed(10), KeyCode: evdev.KEY_B, Action: chkb.KeyActionDown},
			{Time: Elapsed(10 + BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.KeyActionUp},
			{Time: Elapsed(TapDelayMs), KeyCode: evdev.KEY_A, Action: chkb.KeyActionHold},
			{Time: Elapsed(AfterTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			{Time: Elapsed(10 + BeforeTap + TapDelayMs), KeyCode: evdev.KEY_B, Action: chkb.KeyActionTap},
		}),
		Entry("Not Tap/Not Hold, tap must be press/release of the same key without others in between. Hold requires a minimum time", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(50), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
			{Time: Elapsed(BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
			{Time: Elapsed(50 + BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{Time: Elapsed(50), KeyCode: evdev.KEY_B, Action: chkb.KeyActionDown},
			{Time: Elapsed(BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionNil},
			{Time: Elapsed(BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			{Time: Elapsed(50 + BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.KeyActionNil},
			{Time: Elapsed(50 + BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.KeyActionUp},
		}),
		Entry("Hold - Tap", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(25), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
			{Time: Elapsed(25 + BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			{Time: Elapsed(50 + AfterTap), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{Time: Elapsed(25), KeyCode: evdev.KEY_B, Action: chkb.KeyActionDown},
			{Time: Elapsed(25 + BeforeTap), KeyCode: evdev.KEY_B, Action: chkb.KeyActionUp},
			{Time: Elapsed(TapDelayMs), KeyCode: evdev.KEY_A, Action: chkb.KeyActionHold},
			{Time: Elapsed(50 + AfterTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			{Time: Elapsed(25 + BeforeTap + TapDelayMs), KeyCode: evdev.KEY_B, Action: chkb.KeyActionTap},
		}),
		Entry("DoubleTap", []chkb.InputEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
			{Time: Elapsed(2 * BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			{Time: Elapsed(3 * BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
		}, []chkb.KeyEvent{
			{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{Time: Elapsed(BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			{Time: Elapsed(2 * BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
			{Time: Elapsed(3 * BeforeTap), KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			{Time: Elapsed(3*BeforeTap + TapDelayMs), KeyCode: evdev.KEY_A, Action: chkb.KeyActionDoubleTap},
		}),
	)
})
