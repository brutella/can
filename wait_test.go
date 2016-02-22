package can

import (
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	rwc := NewEchoReadWriteCloser()
	bus := NewBus(rwc)

	handler := newTestHandler()
	bus.Subscribe(handler)

	go func() {
		<-time.After(time.Millisecond * 50)
		bus.publish(testFrame)
	}()

	resp := <-Wait(bus, 0x5FAF, time.Millisecond*500)

	if x := resp.Err; x != nil {
		t.Fatal(x)
	}
}

func TestTimeoutErr(t *testing.T) {
	rwc := NewEchoReadWriteCloser()
	bus := NewBus(rwc)

	handler := newTestHandler()
	bus.Subscribe(handler)

	resp := <-Wait(bus, 0x5FAF, time.Millisecond*50)

	if x := resp.Err; x == nil {
		t.Fatal(x)
	}
}
