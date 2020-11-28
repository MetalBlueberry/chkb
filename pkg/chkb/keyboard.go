package chkb

type Keyboard struct {
	*Captor
	*Mapper
	*Handler
}

func NewKeyboard(book Book, initialLayer string) *Keyboard {
	kb := &Keyboard{
		Captor:  NewCaptor(),
		Mapper:  NewMapper(book, initialLayer),
		Handler: NewHandler(),
	}
	kb.AddDeliverer(kb.Mapper)
	return kb
}
