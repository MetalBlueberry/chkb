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
	"MetalBlueberry/cheap-keyboard/pkg/deliverers/vkb"
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
						KeyMap: map[chkb.KeyCode]map[chkb.KeyActions][]chkb.MapEvent{
							evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPushLayer, LayerName: "swapAB"}}},
						},
					},
					"swapAB": {
						KeyMap: map[chkb.KeyCode]map[chkb.KeyActions][]chkb.MapEvent{
							evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPopLayer, LayerName: "swapAB"}}},
							evdev.KEY_A:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_B}}},
							evdev.KEY_B:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_A}}},
						},
					},
				},
				"base",
			)
			kb.AddDeliverer(&vkb.Keyboard{mockKb})
		})

		DescribeTable("Type",
			func(press []chkb.InputEvent, expect []chkb.KeyEvent) {
				for i := range press {
					fmt.Fprintf(GinkgoWriter, "Input %v %s\n", press[i].Action, evdev.KEY[int(press[i].KeyCode)])
					events, err := kb.CaptureOne(press[i])
					assert.NoError(GinkgoT(), err, "Capture should not fail")
					mevents, err := kb.Maps(events)
					assert.NoError(GinkgoT(), err, "Maps should not fail")
					err = kb.Delivers(mevents)
					assert.NoError(GinkgoT(), err, "Deliver should not fail")
				}
				assert.Equal(GinkgoT(), expect, mockKb.Events)
			},
			Entry("Empty", []chkb.InputEvent{}, []chkb.KeyEvent{}),
			Entry("Forward AB", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionUp},
			}),
			Entry("Push layer swap AB", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
			}),
			Entry("Pop layer swap AB", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
				{Time: Elapsed(6), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(7), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_A, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_A, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionUp},
			}),
		)
	})
	Describe("push down pop up", func() {
		BeforeEach(func() {
			mockKb = &TestKeyboard{[]chkb.KeyEvent{}}
			kb = chkb.NewKeyboard(
				chkb.Book{
					"base": {
						KeyMap: map[chkb.KeyCode]map[chkb.KeyActions][]chkb.MapEvent{
							evdev.KEY_CAPSLOCK: {
								chkb.KeyActionDown: {{Action: chkb.KbActionPushLayer, LayerName: "easyenter"}},
								chkb.KeyActionUp:   {{Action: chkb.KbActionPopLayer, LayerName: "easyenter"}},
							},
						},
					},
					"easyenter": {
						KeyMap: map[chkb.KeyCode]map[chkb.KeyActions][]chkb.MapEvent{
							evdev.KEY_SEMICOLON: {chkb.KeyActionMap: {{KeyCode: evdev.KEY_ENTER}}},
						},
					},
				},
				"base",
			)
			kb.AddDeliverer(&vkb.Keyboard{mockKb})
		})
		DescribeTable("keyup must always be equal to keydown code",
			func(press []chkb.InputEvent, expect []chkb.KeyEvent) {
				for i := range press {
					fmt.Fprintf(GinkgoWriter, "Input %v %s\n", press[i].Action, evdev.KEY[int(press[i].KeyCode)])
					events, err := kb.CaptureOne(press[i])
					assert.NoError(GinkgoT(), err, "Capture should not fail")
					mevents, err := kb.Maps(events)
					assert.NoError(GinkgoT(), err, "Maps should not fail")
					err = kb.Delivers(mevents)
					assert.NoError(GinkgoT(), err, "Deliver should not fail")
				}
				assert.Equal(GinkgoT(), expect, mockKb.Events)
			},
			Entry("up after layer pop", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_ENTER, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_ENTER, Action: chkb.KeyActionUp},
			}),
			Entry("up after layer push", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionUp},
			}),
			Entry("colon/enter/colon", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(4), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
				{Time: Elapsed(5), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
				{Time: Elapsed(6), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(7), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_ENTER, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_ENTER, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionUp},
			}),
			Entry("quick colon/enter/colon", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
				{Time: Elapsed(6), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(7), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_ENTER, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_ENTER, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KeyActionUp},
			}),
		)
	})
	Describe("Map and down/up", func() {
		BeforeEach(func() {
			mockKb = &TestKeyboard{[]chkb.KeyEvent{}}
			kb = chkb.NewKeyboard(
				chkb.Book{
					"base": {
						KeyMap: map[chkb.KeyCode]map[chkb.KeyActions][]chkb.MapEvent{
							evdev.KEY_CAPSLOCK: {
								chkb.KeyActionDown: {{Action: chkb.KbActionTap, KeyCode: chkb.KEY_0}},
								chkb.KeyActionUp:   {{Action: chkb.KbActionTap, KeyCode: chkb.KEY_1}},
								chkb.KeyActionMap:  {{KeyCode: chkb.KEY_LEFTMETA}},
							},
						},
					},
				},
				"base",
			)
			kb.AddDeliverer(&vkb.Keyboard{mockKb})
		})
		DescribeTable("Should do actions up/down and mapkey",
			func(press []chkb.InputEvent, expect []chkb.KeyEvent) {
				for i := range press {
					fmt.Fprintf(GinkgoWriter, "Input %v %s\n", press[i].Action, evdev.KEY[int(press[i].KeyCode)])
					events, err := kb.CaptureOne(press[i])
					assert.NoError(GinkgoT(), err, "Capture should not fail")
					mevents, err := kb.Maps(events)
					assert.NoError(GinkgoT(), err, "Maps should not fail")
					err = kb.Delivers(mevents)
					assert.NoError(GinkgoT(), err, "Deliver should not fail")
				}
				assert.Equal(GinkgoT(), expect, mockKb.Events)
			},
			Entry("map and up/down", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_LEFTMETA, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_0, Action: chkb.KeyActionTap},
				{KeyCode: evdev.KEY_LEFTMETA, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_1, Action: chkb.KeyActionTap},
			}),
		)
	})
	Describe("Multiple actions", func() {
		BeforeEach(func() {
			mockKb = &TestKeyboard{[]chkb.KeyEvent{}}
			kb = chkb.NewKeyboard(
				chkb.Book{
					"base": {
						KeyMap: map[chkb.KeyCode]map[chkb.KeyActions][]chkb.MapEvent{
							evdev.KEY_CAPSLOCK: {
								chkb.KeyActionDown: {
									{Action: chkb.KbActionDown, KeyCode: chkb.KEY_LEFTMETA},
									{Action: chkb.KbActionPushLayer, LayerName: "swapAB"},
								},
								chkb.KeyActionUp: {
									{Action: chkb.KbActionUp, KeyCode: chkb.KEY_LEFTMETA},
									{Action: chkb.KbActionPopLayer, LayerName: "swapAB"},
								},
							},
							evdev.KEY_F1: {
								chkb.KeyActionTap: {
									{Action: chkb.KbActionTap, KeyCode: chkb.KEY_H},
									{Action: chkb.KbActionTap, KeyCode: chkb.KEY_E},
									{Action: chkb.KbActionTap, KeyCode: chkb.KEY_L},
									{Action: chkb.KbActionTap, KeyCode: chkb.KEY_L},
									{Action: chkb.KbActionTap, KeyCode: chkb.KEY_O},
								},
								chkb.KeyActionMap: {
									{Action: chkb.KbActionNil},
								},
							},
						},
					},
					"swapAB": {
						KeyMap: map[chkb.KeyCode]map[chkb.KeyActions][]chkb.MapEvent{
							evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPopLayer}}},
							evdev.KEY_A:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_B}}},
							evdev.KEY_B:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_A}}},
						},
					},
				},
				"base",
			)
			kb.AddDeliverer(&vkb.Keyboard{mockKb})
		})
		DescribeTable("Should do multiple actions",
			func(press []chkb.InputEvent, expect []chkb.KeyEvent) {
				for i := range press {
					fmt.Fprintf(GinkgoWriter, "Input %v %s\n", press[i].Action, evdev.KEY[int(press[i].KeyCode)])
					events, err := kb.CaptureOne(press[i])
					assert.NoError(GinkgoT(), err, "Capture should not fail")
					mevents, err := kb.Maps(events)
					assert.NoError(GinkgoT(), err, "Maps should not fail")
					err = kb.Delivers(mevents)
					assert.NoError(GinkgoT(), err, "Deliver should not fail")
				}
				assert.Equal(GinkgoT(), expect, mockKb.Events)
			},
			Entry("push layer and up/down", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_LEFTMETA, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionDown},
				{KeyCode: evdev.KEY_B, Action: chkb.KeyActionUp},
				{KeyCode: evdev.KEY_LEFTMETA, Action: chkb.KeyActionUp},
			}),
			Entry("type on tap", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_F1, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_F1, Action: chkb.InputActionUp},
			}, []chkb.KeyEvent{
				{KeyCode: evdev.KEY_H, Action: chkb.KeyActionTap},
				{KeyCode: evdev.KEY_E, Action: chkb.KeyActionTap},
				{KeyCode: evdev.KEY_L, Action: chkb.KeyActionTap},
				{KeyCode: evdev.KEY_L, Action: chkb.KeyActionTap},
				{KeyCode: evdev.KEY_O, Action: chkb.KeyActionTap},
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
		Action:  chkb.KeyActionTap,
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
		Action:  chkb.KeyActionDown,
	})
	fmt.Fprintf(GinkgoWriter, "Output down %s\n", evdev.KEY[key])
	return nil
}

// KeyUp will send a keyrelease event to an existing keyboard device.
// The key can be any of the predefined keycodes from keycodes.go.
func (kb *TestKeyboard) KeyUp(key int) error {
	kb.Events = append(kb.Events, chkb.KeyEvent{
		KeyCode: chkb.KeyCode(key),
		Action:  chkb.KeyActionUp,
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
