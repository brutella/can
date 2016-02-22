package can

// The Handler interfaces defines a method to receive a frame.
type Handler interface {
	Handle(frame Frame)
}

type HandlerFunc func(frame Frame)

type handler struct {
	fn HandlerFunc
}

func NewHandler(fn HandlerFunc) Handler {
	return &handler{fn}
}

func (h *handler) Handle(frame Frame) {
	h.fn(frame)
}
