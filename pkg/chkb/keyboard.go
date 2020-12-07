package chkb

import "time"

const (
	TapDelay = 200 * time.Millisecond
)

type Keyboard struct {
	*Captor
	*Mapper
	*Handler
}

func NewKeyboard(book Book, initialLayer string) *Keyboard {
	kb := &Keyboard{
		Captor: NewCaptor(),
		// Mapper:  NewMapper(book, initialLayer),
		Mapper:  NewMapper(),
		Handler: NewHandler(),
	}
	kb.AddDeliverer(kb.Mapper)
	return kb
}

func (kb *Keyboard) Run(event func() ([]InputEvent, error)) error {
	return kb.Captor.Run(event, func(captured []KeyEvent) error {
		mapped, err := kb.Maps(captured)
		if err != nil {
			return err
		}

		err = kb.Delivers(mapped)
		if err != nil {
			return err
		}
		return nil
	})
}
