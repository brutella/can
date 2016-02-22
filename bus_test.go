package can

import (
	"reflect"
	"testing"
	"time"
)

var testFrame = Frame{
	ID:     0x5FAF,
	Length: 0x8,
	Flags:  0x0,
	Res0:   0x0,
	Res1:   0x0,
	Data:   [MaxFrameDataLength]uint8{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8},
}

type testHandler struct {
	frame Frame
}

func newTestHandler() *testHandler {
	return &testHandler{}
}

func (h *testHandler) Handle(frame Frame) {
	h.frame = frame
}

func TestPublish(t *testing.T) {
	rwc := NewEchoReadWriteCloser()
	bus := NewBus(rwc)

	handler := newTestHandler()
	bus.Subscribe(handler)

	go bus.publishNextFrame()

	rwc.WriteFrame(testFrame)

	<-time.After(time.Millisecond * 50)

	if is, want := handler.frame, testFrame; reflect.DeepEqual(is, want) == false {
		t.Fatalf("is=% X want=% X", is, want)
	}
}

func TestSubscribe(t *testing.T) {
	rwc := NewEchoReadWriteCloser()
	bus := NewBus(rwc)

	handler := newTestHandler()
	bus.Subscribe(handler)
	bus.publish(testFrame)

	if is, want := handler.frame, testFrame; reflect.DeepEqual(is, want) == false {
		t.Fatalf("is=% X want=% X", is, want)
	}
}

func TestUnsubscribe(t *testing.T) {
	rwc := NewEchoReadWriteCloser()
	bus := NewBus(rwc)

	handler := newTestHandler()
	bus.Subscribe(handler)

	if x := bus.contains(handler); x != true {
		t.Fatal(x)
	}

	bus.Unsubscribe(handler)

	if x := bus.contains(handler); x != false {
		t.Fatal(x)
	}
}
