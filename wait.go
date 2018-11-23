package can

import (
	"fmt"
	"time"
)

// A WaitResponse encapsulates the response of waiting for a frame.
type WaitResponse struct {
	Frame Frame
	Err   error
}

// The Waiter interfaces defines a method to wait for a frame.
type WaiterInterface interface {
	Wait(timeout time.Duration) <-chan WaitResponse
}

type Waiter struct {
	fn WaitFunc
}

type WaitFunc func(bus *Bus, id uint32, timeout time.Duration) <-chan WaitResponse

type waiter struct {
	id     uint32
	wait   chan WaitResponse
	bus    *Bus
	filter Handler
	fn     WaitFunc
}

// NewWaiter returns a new waiter which calls fn when going to wait for a frame.
func NewWaiter(fn WaitFunc) *Waiter {
	return &Waiter{}
}

// Wait returns a channel, which receives a frame or an error, if the
// frame with the expected id didn't arrive on time.
func (w *Waiter) Wait(bus *Bus, id uint32, timeout time.Duration) <-chan WaitResponse {
	waiter := waiter{
		id:   id,
		wait: make(chan WaitResponse),
		bus:  bus,
	}

	ch := make(chan WaitResponse)

	go func() {
		select {
		case resp := <-waiter.wait:
			ch <- resp
		case <-time.After(timeout):
			err := fmt.Errorf("Timeout error waiting for %X", id)
			ch <- WaitResponse{Frame{}, err}
		}
	}()

	waiter.filter = newFilter(id, &waiter)
	bus.Subscribe(waiter.filter)

	return ch
}

// ensure backward compatibility
// Wait returns a channel, which receives a frame or an error, if the
// frame with the expected id didn't arrive on time.
func Wait(bus *Bus, id uint32, timeout time.Duration) <-chan WaitResponse {
	waiter := waiter{
		id:   id,
		wait: make(chan WaitResponse),
		bus:  bus,
	}

	ch := make(chan WaitResponse)

	go func() {
		select {
		case resp := <-waiter.wait:
			ch <- resp
		case <-time.After(timeout):
			err := fmt.Errorf("Timeout error waiting for %X", id)
			ch <- WaitResponse{Frame{}, err}
		}
	}()

	waiter.filter = newFilter(id, &waiter)
	bus.Subscribe(waiter.filter)

	return ch
}

func (w *waiter) Handle(frame Frame) {
	w.bus.Unsubscribe(w.filter)
	w.wait <- WaitResponse{frame, nil}
}
