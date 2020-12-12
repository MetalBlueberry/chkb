package chkb_test

import (
	"errors"
	"syscall"
	"time"

	"github.com/benbjohnson/clock"
	evdev "github.com/gvalkov/golang-evdev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	log "github.com/sirupsen/logrus"

	// . "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"

	"MetalBlueberry/cheap-keyboard/pkg/chkb"
)

var _ = Describe("Keyboard", func() {

	var (
		clockMock *clock.Mock
		mockKb    *TestKeyboard
		kb        *chkb.Keyboard
	)

	BeforeEach(func() {
		clockMock = clock.NewMock()
		mockKb = &TestKeyboard{[]chkb.MapEvent{}}
	})
	JustBeforeEach(func() {
		kb.Captor = chkb.NewCaptorWithClock(clockMock)
	})

	RunTableTest := func(events []chkb.InputEvent, expect []chkb.MapEvent) {
		i := 0
		finished := errors.New("Finished")
		err := kb.Run(func() ([]chkb.InputEvent, error) {
			if len(events) == i {
				clockMock.Add(chkb.TapDelay)
				return nil, finished
			}
			event := events[i]
			log.Println(event)
			clockMock.Set(event.Time)
			i++
			return []chkb.InputEvent{event}, nil
		})
		assert.Equal(GinkgoT(), err, finished)
		assert.Equal(GinkgoT(), expect, mockKb.Events)
	}

	Describe("Single layer swap A-B", func() {
		BeforeEach(func() {
			kb = chkb.NewKeyboard(
				chkb.Config{
					Layers: chkb.LayerBook{
						"base": {
							KeyMap: chkb.KeyMap{
								evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPushLayer, LayerName: "swapAB"}}},
							},
						},
						"swapAB": {
							KeyMap: chkb.KeyMap{
								evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPopLayer, LayerName: "swapAB"}}},
								evdev.KEY_A:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_B}}},
								evdev.KEY_B:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_A}}},
							},
						},
					},
				},
				"base",
			)
			kb.AddDeliverer(mockKb)
		})
		DescribeTable("Type", RunTableTest,
			Entry("One", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.KbActionDown},
			}),
			Entry("Forward AB", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_A, Action: chkb.KbActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_A, Action: chkb.KbActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_B, Action: chkb.KbActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_B, Action: chkb.KbActionUp},
			}),
		)
		DescribeTable("Layers", RunTableTest,
			Entry("Push layer swap AB", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_B, Action: chkb.KbActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_B, Action: chkb.KbActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_A, Action: chkb.KbActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_A, Action: chkb.KbActionUp},
			}),
			Entry("Pop layer swap AB", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
				{Time: Elapsed(AfterTap), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(AfterTap + 1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(AfterTap + 2), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(AfterTap + 3), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(AfterTap + 4), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(AfterTap + 5), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_B, Action: chkb.KbActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_B, Action: chkb.KbActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_A, Action: chkb.KbActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_A, Action: chkb.KbActionUp},
				{Time: Elapsed(AfterTap), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{Time: Elapsed(AfterTap + 1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{Time: Elapsed(AfterTap + 2), KeyCode: evdev.KEY_A, Action: chkb.KbActionDown},
				{Time: Elapsed(AfterTap + 3), KeyCode: evdev.KEY_A, Action: chkb.KbActionUp},
				{Time: Elapsed(AfterTap + 4), KeyCode: evdev.KEY_B, Action: chkb.KbActionDown},
				{Time: Elapsed(AfterTap + 5), KeyCode: evdev.KEY_B, Action: chkb.KbActionUp},
			}),
		)
	})
	Describe("push down pop up", func() {
		BeforeEach(func() {
			kb = chkb.NewKeyboard(
				chkb.Config{
					Layers: chkb.LayerBook{
						"base": {
							KeyMap: map[chkb.KeyCode]chkb.KeyMapActions{
								evdev.KEY_CAPSLOCK: {
									chkb.KeyActionDown: {{Action: chkb.KbActionPushLayer, LayerName: "easyenter"}},
									chkb.KeyActionUp:   {{Action: chkb.KbActionPopLayer, LayerName: "easyenter"}},
								},
							},
						},
						"easyenter": {
							KeyMap: map[chkb.KeyCode]chkb.KeyMapActions{
								evdev.KEY_SEMICOLON: {chkb.KeyActionMap: {{KeyCode: evdev.KEY_ENTER}}},
							},
						},
					},
				},
				"base",
			)
			kb.AddDeliverer(mockKb)
		})
		DescribeTable("keyup must always be equal to keydown code", RunTableTest,
			Entry("up after layer pop", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(1), KeyCode: evdev.KEY_ENTER, Action: chkb.KbActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_ENTER, Action: chkb.KbActionUp},
			}),
			Entry("up after layer push", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionUp},
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
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_ENTER, Action: chkb.KbActionDown},
				{Time: Elapsed(4), KeyCode: evdev.KEY_ENTER, Action: chkb.KbActionUp},
				{Time: Elapsed(6), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionDown},
				{Time: Elapsed(7), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionUp},
			}),
			Entry("quick colon/enter/colon", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(4), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
				{Time: Elapsed(5), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
				{Time: Elapsed(6), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionDown},
				{Time: Elapsed(7), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_ENTER, Action: chkb.KbActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_ENTER, Action: chkb.KbActionUp},
				{Time: Elapsed(6), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionDown},
				{Time: Elapsed(7), KeyCode: evdev.KEY_SEMICOLON, Action: chkb.KbActionUp},
			}),
		)
	})
	Describe("Map and down/up", func() {
		BeforeEach(func() {
			kb = chkb.NewKeyboard(
				chkb.Config{
					Layers: chkb.LayerBook{
						"base": {
							KeyMap: map[chkb.KeyCode]chkb.KeyMapActions{
								evdev.KEY_CAPSLOCK: {
									chkb.KeyActionDown: {{Action: chkb.KbActionTap, KeyCode: chkb.KEY_0}},
									chkb.KeyActionUp:   {{Action: chkb.KbActionTap, KeyCode: chkb.KEY_1}},
									chkb.KeyActionMap:  {{KeyCode: chkb.KEY_LEFTMETA}},
								},
							},
						},
					},
				},
				"base",
			)
			kb.AddDeliverer(mockKb)
		})
		DescribeTable("Should do actions up/down and mapkey", RunTableTest,
			Entry("map and up/down", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTMETA, Action: chkb.KbActionDown},
				{Time: Elapsed(0), KeyCode: evdev.KEY_0, Action: chkb.KbActionTap},
				{Time: Elapsed(2), KeyCode: evdev.KEY_LEFTMETA, Action: chkb.KbActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_1, Action: chkb.KbActionTap},
			}),
		)
	})
	Describe("Multiple actions", func() {
		BeforeEach(func() {
			mockKb = &TestKeyboard{[]chkb.MapEvent{}}
			kb = chkb.NewKeyboard(
				chkb.Config{
					Layers: chkb.LayerBook{
						"base": {
							KeyMap: map[chkb.KeyCode]chkb.KeyMapActions{
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
							KeyMap: map[chkb.KeyCode]chkb.KeyMapActions{
								evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPopLayer}}},
								evdev.KEY_A:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_B}}},
								evdev.KEY_B:         {chkb.KeyActionMap: {{KeyCode: evdev.KEY_A}}},
							},
						},
					},
				},
				"base",
			)
			kb.AddDeliverer(mockKb)
		})
		DescribeTable("Should do multiple actions", RunTableTest,
			Entry("push layer and up/down", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_CAPSLOCK, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTMETA, Action: chkb.KbActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_B, Action: chkb.KbActionDown},
				{Time: Elapsed(2), KeyCode: evdev.KEY_B, Action: chkb.KbActionUp},
				{Time: Elapsed(3), KeyCode: evdev.KEY_LEFTMETA, Action: chkb.KbActionUp},
			}),
			Entry("type on tap", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_F1, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_F1, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(TapDelayMs + 1), KeyCode: evdev.KEY_H, Action: chkb.KbActionTap},
				{Time: Elapsed(TapDelayMs + 1), KeyCode: evdev.KEY_E, Action: chkb.KbActionTap},
				{Time: Elapsed(TapDelayMs + 1), KeyCode: evdev.KEY_L, Action: chkb.KbActionTap},
				{Time: Elapsed(TapDelayMs + 1), KeyCode: evdev.KEY_L, Action: chkb.KbActionTap},
				{Time: Elapsed(TapDelayMs + 1), KeyCode: evdev.KEY_O, Action: chkb.KbActionTap},
			}),
		)
	})
	Describe("Layers", func() {
		BeforeEach(func() {
			kb = chkb.NewKeyboard(
				chkb.Config{
					Layers: chkb.LayerBook{
						"base": {
							KeyMap: chkb.KeyMap{
								evdev.KEY_LEFTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionChangeLayer, LayerName: "swapAB"}}},
							},
						},
						"swapAB": {
							KeyMap: chkb.KeyMap{
								evdev.KEY_LEFTSHIFT:  {chkb.KeyActionTap: {{Action: chkb.KbActionChangeLayer, LayerName: "base"}}},
								evdev.KEY_RIGHTSHIFT: {chkb.KeyActionTap: {{Action: chkb.KbActionPopLayer, LayerName: "swapAB"}}},
								evdev.KEY_A:          {chkb.KeyActionMap: {{KeyCode: evdev.KEY_B}}},
								evdev.KEY_B:          {chkb.KeyActionMap: {{KeyCode: evdev.KEY_A}}},
							},
						},
					},
				},
				"base",
			)
			kb.AddDeliverer(mockKb)
		})
		DescribeTable("Actions", RunTableTest,
			Entry("Cannot pop last layer", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_RIGHTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_RIGHTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(6), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(7), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_RIGHTSHIFT, Action: chkb.KbActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_RIGHTSHIFT, Action: chkb.KbActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_B, Action: chkb.KbActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_B, Action: chkb.KbActionUp},
				{Time: Elapsed(6), KeyCode: evdev.KEY_A, Action: chkb.KbActionDown},
				{Time: Elapsed(7), KeyCode: evdev.KEY_A, Action: chkb.KbActionUp},
			}),
			Entry("change layer swap AB", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_B, Action: chkb.KbActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_B, Action: chkb.KbActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_A, Action: chkb.KbActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_A, Action: chkb.KbActionUp},
			}),
			Entry("change and come back layer swap AB", []chkb.InputEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
				{Time: Elapsed(AfterTap), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionDown},
				{Time: Elapsed(AfterTap + 1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.InputActionUp},
				{Time: Elapsed(AfterTap + 2), KeyCode: evdev.KEY_A, Action: chkb.InputActionDown},
				{Time: Elapsed(AfterTap + 3), KeyCode: evdev.KEY_A, Action: chkb.InputActionUp},
				{Time: Elapsed(AfterTap + 4), KeyCode: evdev.KEY_B, Action: chkb.InputActionDown},
				{Time: Elapsed(AfterTap + 5), KeyCode: evdev.KEY_B, Action: chkb.InputActionUp},
			}, []chkb.MapEvent{
				{Time: Elapsed(0), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{Time: Elapsed(1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{Time: Elapsed(2), KeyCode: evdev.KEY_B, Action: chkb.KbActionDown},
				{Time: Elapsed(3), KeyCode: evdev.KEY_B, Action: chkb.KbActionUp},
				{Time: Elapsed(4), KeyCode: evdev.KEY_A, Action: chkb.KbActionDown},
				{Time: Elapsed(5), KeyCode: evdev.KEY_A, Action: chkb.KbActionUp},
				{Time: Elapsed(AfterTap), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionDown},
				{Time: Elapsed(AfterTap + 1), KeyCode: evdev.KEY_LEFTSHIFT, Action: chkb.KbActionUp},
				{Time: Elapsed(AfterTap + 2), KeyCode: evdev.KEY_A, Action: chkb.KbActionDown},
				{Time: Elapsed(AfterTap + 3), KeyCode: evdev.KEY_A, Action: chkb.KbActionUp},
				{Time: Elapsed(AfterTap + 4), KeyCode: evdev.KEY_B, Action: chkb.KbActionDown},
				{Time: Elapsed(AfterTap + 5), KeyCode: evdev.KEY_B, Action: chkb.KbActionUp},
			}),
		)
	})
})

type TestKeyboard struct {
	Events []chkb.MapEvent
}

func (kb *TestKeyboard) Deliver(event chkb.MapEvent) (bool, error) {
	switch event.Action {
	case chkb.KbActionUp, chkb.KbActionDown, chkb.KbActionTap:
		kb.Events = append(kb.Events, event)
		log.
			WithField("Key", event.KeyCode).
			WithField("Action", event.Action).
			Print("Deliver")
		return true, nil
	}
	return false, nil
}

func Clock() func(increment time.Duration) syscall.Timeval {
	init := time.Date(2020, 11, 20, 12, 0, 0, 0, time.UTC)
	return func(increment time.Duration) syscall.Timeval {
		init.Add(increment)
		return Time(init)
	}
}
