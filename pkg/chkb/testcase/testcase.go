package testcase

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MetalBlueberry/chkb/pkg/chkb"
	"github.com/benbjohnson/clock"
)

var (
	InitialTime             = time.Date(2020, 11, 20, 12, 0, 0, 0, time.UTC)
	DefaultDelayBetweenKeys = 10 * time.Millisecond
)

func MustNew(seq string) []chkb.InputEvent {
	events, err := Input(seq)
	if err != nil {
		panic(err)
	}
	return events
}

func Output(seq string) ([]chkb.InputEvent, error) {

	events := make([]chkb.InputEvent, 0)
	stepper := NewStepper(InitialTime, DefaultDelayBetweenKeys)

	fragments := strings.Split(seq, " ")
	for _, fragment := range fragments {
		if len(fragment) == 0 {
			continue
		}
		switch fragment[0] {
		case 'T':
			keyCode, err := parseKey(fragment[1:])
			if err != nil {
				return nil, err
			}
			events = append(events, chkb.InputEvent{
				Time:    stepper.Step(),
				KeyCode: keyCode,
				Action:  chkb.InputActionDown,
			})
			events = append(events, chkb.InputEvent{
				Time:    stepper.Step(),
				KeyCode: keyCode,
				Action:  chkb.InputActionUp,
			})
		case 'P':
			keyCode, err := parseKey(fragment[1:])
			if err != nil {
				return nil, err
			}
			events = append(events, chkb.InputEvent{
				Time:    stepper.Step(),
				KeyCode: keyCode,
				Action:  chkb.InputActionDown,
			})
		case 'R':
			keyCode, err := parseKey(fragment[1:])
			if err != nil {
				return nil, err
			}
			events = append(events, chkb.InputEvent{
				Time:    stepper.Step(),
				KeyCode: keyCode,
				Action:  chkb.InputActionUp,
			})
		default:
			delay, err := strconv.Atoi(fragment)
			if err != nil {
				return nil, fmt.Errorf("Invalid time fragment, %w", err)
			}
			stepper.Add(time.Duration(delay) * time.Millisecond)
		}
	}

	return events, nil
}

func Input(seq string) ([]chkb.InputEvent, error) {

	events := make([]chkb.InputEvent, 0)
	stepper := NewStepper(InitialTime, DefaultDelayBetweenKeys)

	fragments := strings.Split(seq, " ")
	for _, fragment := range fragments {
		if len(fragment) == 0 {
			continue
		}
		switch fragment[0] {
		case 'T':
			keyCode, err := parseKey(fragment[1:])
			if err != nil {
				return nil, err
			}
			events = append(events, chkb.InputEvent{
				Time:    stepper.Step(),
				KeyCode: keyCode,
				Action:  chkb.InputActionDown,
			})
			events = append(events, chkb.InputEvent{
				Time:    stepper.Step(),
				KeyCode: keyCode,
				Action:  chkb.InputActionUp,
			})
		case 'P':
			keyCode, err := parseKey(fragment[1:])
			if err != nil {
				return nil, err
			}
			events = append(events, chkb.InputEvent{
				Time:    stepper.Step(),
				KeyCode: keyCode,
				Action:  chkb.InputActionDown,
			})
		case 'R':
			keyCode, err := parseKey(fragment[1:])
			if err != nil {
				return nil, err
			}
			events = append(events, chkb.InputEvent{
				Time:    stepper.Step(),
				KeyCode: keyCode,
				Action:  chkb.InputActionUp,
			})
		default:
			delay, err := strconv.Atoi(fragment)
			if err != nil {
				return nil, fmt.Errorf("Invalid time fragment, %w", err)
			}
			stepper.Add(time.Duration(delay) * time.Millisecond)
		}
	}

	return events, nil
}

func parseKey(key string) (chkb.KeyCode, error) {
	switch key {
	case "a":
		return chkb.KEY_A, nil
	case "b":
		return chkb.KEY_B, nil
	case "c":
		return chkb.KEY_C, nil
	case "d":
		return chkb.KEY_D, nil
	}
	return chkb.KEY_RESERVED, fmt.Errorf("Invalid Key %s", key)
}

type Stepper struct {
	*clock.Mock
	Initial  time.Time
	Offset   time.Duration
	StepSize time.Duration
}

func NewStepper(initialTime time.Time, stepSize time.Duration) *Stepper {
	s := &Stepper{
		Mock:     &clock.Mock{},
		Initial:  initialTime,
		StepSize: stepSize,
	}
	s.Set(initialTime)
	return s
}

func (s *Stepper) update() {
	s.Set(s.Initial.Add(s.Offset))
}

func (s *Stepper) Step() time.Time {
	defer s.update()
	s.Offset += s.StepSize
	return s.Now()
}

func (s *Stepper) Add(d time.Duration) time.Time {
	defer s.update()
	s.Offset += d
	return s.Now()
}

func Elapsed(steps int, ms int64) time.Time {
	return Step(steps).
		Add(time.Duration(ms) * time.Millisecond)
}

func Step(n int) time.Time {
	return InitialTime.
		Add(time.Duration(n) * DefaultDelayBetweenKeys)
}
