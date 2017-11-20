package can

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"unsafe"
)

// The Reader interface extends the `io.Reader` interface by method
// to read a frame.
type Reader interface {
	io.Reader
	ReadFrame(*Frame) error
}

// The Writer interface extends the `io.Writer` interface by method
// to write a frame.
type Writer interface {
	io.Writer
	WriteFrame(Frame) error
}

// The ReadWriteCloser interface combines the Reader and Writer and
// `io.Closer` interface.
type ReadWriteCloser interface {
	Reader
	Writer

	io.Closer
}

// Socket protocols
const (
	Raw   uint8 = 1 // CAN_RAW
	Bcm   uint8 = 2 // CAN_BCM
	TP16  uint8 = 3
	TP20  uint8 = 4
	MCNet uint8 = 5
	ISOTp uint8 = 6
)

type readWriteCloser struct {
	rwc io.ReadWriteCloser
	s   int
}

// NewReadWriteCloserForInterface returns a ReadWriteCloser for a network interface.
func NewReadWriteCloserForInterface(i *net.Interface) (ReadWriteCloser, error) {
	proto := Raw
	s, err := syscall.Socket(AF_CAN, syscall.SOCK_RAW, int(proto) /* 0? */)
	if err != nil {
		return nil, err
	}

	addr := NewSockaddr(uint16(proto) /* can.AF_CAN? */, i.Index /* 0  for all interfaces? */)

	if err := syscall.Bind(s, addr); err != nil {
		return nil, err
	}

	f := os.NewFile(uintptr(s), fmt.Sprintf("fd %d", s))

	return &readWriteCloser{f, s}, nil
}

// NewReadWriteCloser returns a ReadWriteCloser for an `io.ReadWriteCloser`.
func NewReadWriteCloser(rwc io.ReadWriteCloser) ReadWriteCloser {
	return &readWriteCloser{rwc, 0}
}

// try to get timestamp via ioctl SIOCGSTAMP request
func (rwc *readWriteCloser) setTimestamp(frame *Frame) {
	fd := uintptr(rwc.s)
	req := uintptr(syscall.SIOCGSTAMP)
	arg := uintptr(unsafe.Pointer(&frame.Time))
	syscall.Syscall(syscall.SYS_IOCTL, fd, req, arg)
}

func (rwc *readWriteCloser) ReadFrame(frame *Frame) error {
	b := make([]byte, 256) // TODO(brutella) optimize size
	n, err := rwc.Read(b)

	if err != nil {
		return err
	}

	rwc.setTimestamp(frame)

	err = Unmarshal(b[:n], frame)

	return err
}

func (rwc *readWriteCloser) WriteFrame(frame Frame) error {
	b, err := Marshal(frame)

	if err != nil {
		return err
	}

	_, err = rwc.Write(b)

	return err
}

func (rwc *readWriteCloser) Read(b []byte) (n int, err error) {
	return rwc.rwc.Read(b)
}

func (rwc *readWriteCloser) Write(b []byte) (n int, err error) {
	return rwc.rwc.Write(b)
}

func (rwc *readWriteCloser) Close() error {
	return rwc.rwc.Close()
}
