package can

import (
	"io"
	"net"
)

// Bus represents the CAN bus.
// Handlers can subscribe to receive frames.
// Frame are sent using the *Publish* method.
type Bus struct {
	rwc ReadWriteCloser

	handler []Handler
}

// NewBusForInterfaceWithName returns a bus from the network interface with name ifaceName.
func NewBusForInterfaceWithName(ifaceName string) (*Bus, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}

	conn, err := NewReadWriteCloserForInterface(iface)
	if err != nil {
		return nil, err
	}

	return NewBus(conn), nil
}

// NewBus returns a new CAN bus.
func NewBus(rwc ReadWriteCloser) *Bus {
	return &Bus{
		rwc:     rwc,
		handler: make([]Handler, 0),
	}
}

// ConnectAndPublish starts handling CAN frames to publish them to handlers.
func (b *Bus) ConnectAndPublish() error {
	for {
		err := b.publishNextFrame()
		if err != nil {
			return err
		}
	}

	return nil
}

// Disconnect stops handling CAN frames.
func (b *Bus) Disconnect() error {
	return b.rwc.Close()
}

// Subscribe adds a handler to the bus.
func (b *Bus) Subscribe(handler Handler) {
	b.handler = append(b.handler, handler)
}

// SubscribeFunc adds a function as handler.
func (b *Bus) SubscribeFunc(fn HandlerFunc) {
	handler := NewHandler(fn)
	b.Subscribe(handler)
}

// Unsubscribe removes a handler.
func (b *Bus) Unsubscribe(handler Handler) {
	for i, h := range b.handler {
		if h == handler {
			b.handler = append(b.handler[:i], b.handler[i+1:]...)
			return
		}
	}
}

// Publish publishes a frame on the bus.
//
// Frames publishes with the Publish methods are not received by handlers.
func (b *Bus) Publish(frame Frame) error {
	return b.rwc.WriteFrame(frame)
}

func (b *Bus) contains(handler Handler) bool {
	for _, h := range b.handler {
		if h == handler {
			return true
		}
	}

	return false
}

func (b *Bus) publishNextFrame() error {
	frame := Frame{}
	err := b.rwc.ReadFrame(&frame)
	if err != nil {
		b.rwc.Close()

		if err != io.EOF { // EOF is not an error, it happens when calling rwc.Close()
			return err
		}

		return nil
	}

	b.publish(frame)

	return nil
}

func (b *Bus) publish(frame Frame) {
	for _, h := range b.handler {
		h.Handle(frame)
	}
}
