package chkb_test

import (
	"fmt"
	"syscall"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	// . "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"

	"MetalBlueberry/cheap-keyboard/pkg/chkb"
)

var _ = Describe("Keyboard", func() {

	var (
		mockKb *TestKeyboard
		kb     *chkb.Keyboard
	)

	Describe("Single layer swap A-B", func() {
		BeforeEach(func() {
			mockKb = &TestKeyboard{[]chkb.KeyEvent{}}
			kb = chkb.NewKeyboard(
				chkb.Book{
					"base": {
						KeyMap: map[chkb.KeyCode]map[chkb.Actions]chkb.MapEvent{
							evdev.KEY_LEFTSHIFT: {chkb.ActionTap: {Action: chkb.ActionPush, LayerName: "swapAB"}},
						},
					},
					"swapAB": {
						KeyMap: map[chkb.KeyCode]map[chkb.Actions]chkb.MapEvent{
							evdev.KEY_LEFTSHIFT: {chkb.ActionTap: {Action: chkb.ActionPop}},
							evdev.KEY_A:         {chkb.ActionMap: {KeyCode: evdev.KEY_B}},
							evdev.KEY_B:         {chkb.ActionMap: {KeyCode: evdev.KEY_A}},
						},
					},
				},
				"base",
				mockKb,
			)
		})

		DescribeTable("Type",
			func(press []evdev.InputEvent, expect []chkb.KeyEvent) {
				for i := range press {
					fmt.Fprintf(GinkgoWriter, "Input %v %s\n", evdev.KeyEventState(press[i].Value), evdev.KEY[int(press[i].Code)])
					events, err := kb.CaptureOne(press[i])
					assert.NoError(GinkgoT(), err, "Capture should not fail")
					mevents, err := kb.Maps(events)
					assert.NoError(GinkgoT(), err, "Maps should not fail")
					err = kb.Deliver(mevents)
					assert.NoError(GinkgoT(), err, "Deliver should not fail")
				}
				assert.Equal(GinkgoT(), expect, mockKb.Events)
			},
			Entry("Empty", []evdev.InputEvent{}, []chkb.KeyEvent{}),
			Entry("Forward AB", []evdev.InputEvent{
				{Time: Elapsed(0), Code: evdev.KEY_A, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(1), Code: evdev.KEY_A, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
				{Time: Elapsed(2), Code: evdev.KEY_B, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(3), Code: evdev.KEY_B, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_A, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_A, Action: chkb.ActionUp},
				{KeyCode: evdev.KEY_B, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.ActionUp},
			}),
			Entry("Push layer swap AB", []evdev.InputEvent{
				{Time: Elapsed(0), Code: evdev.KEY_LEFTSHIFT, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(1), Code: evdev.KEY_LEFTSHIFT, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
				{Time: Elapsed(2), Code: evdev.KEY_A, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(3), Code: evdev.KEY_A, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
				{Time: Elapsed(4), Code: evdev.KEY_B, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(5), Code: evdev.KEY_B, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.ActionUp},
				{KeyCode: evdev.KEY_B, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.ActionUp},
				{KeyCode: evdev.KEY_A, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_A, Action: chkb.ActionUp},
			}),
			Entry("Pop layer swap AB", []evdev.InputEvent{
				{Time: Elapsed(0), Code: evdev.KEY_LEFTSHIFT, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(1), Code: evdev.KEY_LEFTSHIFT, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
				{Time: Elapsed(2), Code: evdev.KEY_A, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(3), Code: evdev.KEY_A, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
				{Time: Elapsed(4), Code: evdev.KEY_B, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(5), Code: evdev.KEY_B, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
				{Time: Elapsed(6), Code: evdev.KEY_LEFTSHIFT, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(7), Code: evdev.KEY_LEFTSHIFT, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
				{Time: Elapsed(2), Code: evdev.KEY_A, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(3), Code: evdev.KEY_A, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
				{Time: Elapsed(4), Code: evdev.KEY_B, Value: int32(evdev.KeyDown), Type: evdev.EV_KEY},
				{Time: Elapsed(5), Code: evdev.KEY_B, Value: int32(evdev.KeyUp), Type: evdev.EV_KEY},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.ActionUp},
				{KeyCode: evdev.KEY_B, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.ActionUp},
				{KeyCode: evdev.KEY_A, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_A, Action: chkb.ActionUp},
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.ActionUp},
				{KeyCode: evdev.KEY_A, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_A, Action: chkb.ActionUp},
				{KeyCode: evdev.KEY_B, Action: chkb.ActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.ActionUp},
			}),
		)
	})
})

type TestKeyboard struct {
	Events []chkb.KeyEvent
}

// KeyPress will cause the key to be pressed and immediately released.
func (kb *TestKeyboard) KeyPress(key int) error {
	kb.Events = append(kb.Events, chkb.KeyEvent{
		KeyCode: chkb.KeyCode(key),
		Action:  chkb.ActionTap,
	})

	fmt.Fprintf(GinkgoWriter, "Output tap %s\n", evdev.KEY[key])
	return nil
}

// KeyDown will send a keypress event to an existing keyboard device.
// The key can be any of the predefined keycodes from keycodes.go.
// Note that the key will be "held down" until "KeyUp" is called.
func (kb *TestKeyboard) KeyDown(key int) error {
	kb.Events = append(kb.Events, chkb.KeyEvent{
		KeyCode: chkb.KeyCode(key),
		Action:  chkb.ActionDown,
	})
	fmt.Fprintf(GinkgoWriter, "Output down %s\n", evdev.KEY[key])
	return nil
}

// KeyUp will send a keyrelease event to an existing keyboard device.
// The key can be any of the predefined keycodes from keycodes.go.
func (kb *TestKeyboard) KeyUp(key int) error {
	kb.Events = append(kb.Events, chkb.KeyEvent{
		KeyCode: chkb.KeyCode(key),
		Action:  chkb.ActionUp,
	})
	fmt.Fprintf(GinkgoWriter, "Output up %s\n", evdev.KEY[key])
	return nil
}

func (kb *TestKeyboard) Close() error {
	panic("not implemented") // TODO: Implement
}

func Clock() func(increment time.Duration) syscall.Timeval {
	init := time.Date(2020, 11, 20, 12, 0, 0, 0, time.UTC)
	return func(increment time.Duration) syscall.Timeval {
		init.Add(increment)
		return Time(init)
	}
}

func Time(t time.Time) syscall.Timeval {
	return syscall.Timeval{
		Sec:  t.Unix(),
		Usec: int64(t.UnixNano() / 1000 % 1000000),
	}
}

func Elapsed(ms int64) syscall.Timeval {
	return Time(
		time.
			Date(2020, 11, 20, 12, 0, 0, 0, time.UTC).
			Add(time.Duration(ms) * time.Millisecond),
	)
}
