package can

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func NewReadWriteCloserForInterface(i *net.Interface) (ReadWriteCloser, error) {
	s, _ := syscall.Socket(syscall.AF_CAN, syscall.SOCK_RAW, unix.CAN_RAW)
	addr := &unix.SockaddrCAN{Ifindex: i.Index}
	if err := unix.Bind(s, addr); err != nil {
		return nil, err
	}

	f := os.NewFile(uintptr(s), fmt.Sprintf("fd %d", s))

	return &readWriteCloser{rwc: f, socket: s}, nil
}
