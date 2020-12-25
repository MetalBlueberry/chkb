/*
Copyright © 2020 Víctor Pérez @MetalBlueberry

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package chkb

import (
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
)

// Captor processes stream of keyup/keydown events to generate advanced actions
// such as Tap/DoubleTap/Hold
type Captor struct {
	Clock  clock.Clock
	Config *Config

	idleTimers map[KeyCode]*idleTimer
	lastKey    *idleTimer
	m          sync.Mutex
}

// NewCaptor creates a new captor
func NewCaptor(config *Config) *Captor {
	return NewCaptorWithClock(clock.New(), config)
}

// NewCaptorWithClock creates a new captor with an specific clock, mainly for testing purposes
func NewCaptorWithClock(clock clock.Clock, config *Config) *Captor {
	return &Captor{
		Clock:      clock,
		Config:     config,
		idleTimers: make(map[KeyCode]*idleTimer),
		m:          sync.Mutex{},
	}
}

type idleTimer struct {
	Timeout *clock.Timer
	Time    time.Time
	Count   int
	KeyCode KeyCode
}

func (c *Captor) newIdleTimer(keyCode KeyCode, send MapFunc) *idleTimer {
	timer := &idleTimer{
		Timeout: c.Clock.AfterFunc(c.Config.GetTapDelay(), c.idle(keyCode, send)),
		KeyCode: keyCode,
	}
	timer.Timeout.Stop()
	c.idleTimers[keyCode] = timer
	return timer
}

// InputEvent represents an input keyboard event.
type InputEvent struct {
	Time    time.Time
	KeyCode KeyCode
	Action  InputActions
}

//go:generate stringer -type=InputActions -trimprefix InputAction
// InputActions represents an InputAction iota
type InputActions int

const (
	// InputActionNil identifies an invalid or undefined action
	InputActionNil InputActions = iota
	// InputActionDown identifies a keypress
	InputActionDown
	// InputActionUp identifies a keyrelease
	InputActionUp
)

// MapFunc is used to send KeyEvents to the Mapper module
type MapFunc func([]KeyEvent) error

// CaptureFunc is used to capture input events
type CaptureFunc func() ([]InputEvent, error)

func (c *Captor) loop(capture CaptureFunc, send MapFunc) error {
	events, err := capture()
	if err != nil {
		log.WithError(err).Error("cannot capture more events")
		return err
	}
	c.m.Lock()
	defer c.m.Unlock()

	captured, err := c.capture(events)
	if err != nil {
		log.WithError(err).Error("Cannot capture events")
		return err
	}

	for _, event := range captured {
		timer, ok := c.idleTimers[event.KeyCode]
		switch event.Action {
		case KeyActionDown:
			if !ok {
				timer = c.newIdleTimer(event.KeyCode, send)
			}
			if timer.Timeout.Stop() {
				timer.Count++
			} else {
				timer.Count = 0
			}
			timer.Timeout.Reset(c.Config.GetTapDelay())
			timer.Time = event.Time

			if c.lastKey != nil && c.lastKey != timer {
				// Resolve event due to another keypress
				if c.lastKey.Timeout.Stop() {
					c.deliverEvent(c.lastKey, send)
				}
			}
		case KeyActionUp:
			if ok {
				if timer.Timeout.Stop() {
					timer.Count++
					timer.Timeout.Reset(c.Config.GetTapDelay())
					timer.Time = event.Time
				}
			}
		}

		if timer != nil {
			c.lastKey = timer
		}
	}

	return send(captured)
}

// Run creates a capture infinite loop, The loop will call capture function to get the input events
// and call send function when it captures something.
// capture should block is new data is not available.
// send can be called from different gorouties but will not perform calls in parallel.
func (c *Captor) Run(capture CaptureFunc, send MapFunc) error {
	for {
		err := c.loop(capture, send)
		if err != nil {
			return err
		}
	}
}

func (c *Captor) deliverEvent(event *idleTimer, send func([]KeyEvent) error) {
	if event.Count%2 == 0 {
		// Odd count means that the key is down
		err := send([]KeyEvent{
			NewKeyEv(c.Clock.Now(), event.KeyCode, KeyActionHold),
		})
		if err != nil {
			log.WithError(err).Error("Cannot send key event")
		}
	} else {
		// Even count means that the key is up
		switch event.Count / 2 {
		case 0:
			err := send([]KeyEvent{
				NewKeyEv(c.Clock.Now(), event.KeyCode, KeyActionTap),
			})
			if err != nil {
				log.WithError(err).Error("Cannot send key event")
			}
		case 1:
			err := send([]KeyEvent{
				NewKeyEv(c.Clock.Now(), event.KeyCode, KeyActionDoubleTap),
			})
			if err != nil {
				log.WithError(err).Error("Cannot send key event")
			}
		default:
			err := send([]KeyEvent{
				NewKeyEv(c.Clock.Now(), event.KeyCode, KeyActionNil),
			})
			if err != nil {
				log.WithError(err).Error("Cannot send key event")
			}
		}
	}
}

func (c *Captor) idle(keyCode KeyCode, send func([]KeyEvent) error) func() {
	return func() {
		c.m.Lock()
		defer c.m.Unlock()

		event := c.idleTimers[keyCode]
		c.deliverEvent(event, send)
	}
}

func (c *Captor) capture(events []InputEvent) ([]KeyEvent, error) {
	captured := make([]KeyEvent, 0)
	for i := range events {
		switch events[i].Action {
		case InputActionDown:
			captured = append(captured, NewKeyEv(c.Clock.Now(), events[i].KeyCode, KeyActionDown))
		case InputActionUp:
			captured = append(captured, NewKeyEv(c.Clock.Now(), events[i].KeyCode, KeyActionUp))
		default:
			log.
				WithField("event", events[i]).
				Debug("Skip event")
			continue
		}
	}
	return captured, nil
}

func NewKeyEv(time time.Time, keyCode KeyCode, action KeyActions) KeyEvent {
	return KeyEvent{
		Time:    time,
		KeyCode: keyCode,
		Action:  action,
	}
}
