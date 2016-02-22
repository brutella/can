package can

type filter struct {
	id      uint32
	handler Handler
}

func newFilter(id uint32, handler Handler) Handler {
	return &filter{
		id:      id,
		handler: handler,
	}
}

func (f *filter) Handle(frame Frame) {
	if frame.ID == f.id {
		f.handler.Handle(frame)
	}
}
