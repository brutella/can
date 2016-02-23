package can

// The Handler interfaces defines a method to receive a frame.
type Handler interface {
	Handle(frame Frame)
}

// HandlerFunc defines the function type to handle a frame.
type HandlerFunc func(frame Frame)

type handler struct {
	fn HandlerFunc
}

// NewHandler returns a new handler which calls fn when a frame is received.
func NewHandler(fn HandlerFunc) Handler {
	return &handler{fn}
}

func (h *handler) Handle(frame Frame) {
	h.fn(frame)
}
