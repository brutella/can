package can

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
)

const maskCobID = 0x7FF

func NewReadWriteCloserForInterface(i *net.Interface) (ReadWriteCloser, error) {
	s, _ := syscall.Socket(syscall.AF_CAN, syscall.SOCK_RAW, unix.CAN_RAW)
	addr := &unix.SockaddrCAN{Ifindex: i.Index}
	if err := unix.Bind(s, addr); err != nil {
		return nil, err
	}

	f := os.NewFile(uintptr(s), fmt.Sprintf("fd %d", s))

	return &readWriteCloser{rwc: f, socket: s}, nil
}

func (rwc *readWriteCloser) setFilter(allowedIds []uint32) error {
	if rwc.socket == 0 {
		return ErrorKernelFilterNotSupported
	}

	filter := make([]unix.CanFilter, len(allowedIds))

	for i, allowedId := range allowedIds {
		filter[i].Id = allowedId
		filter[i].Mask = maskCobID
	}

	return unix.SetsockoptCanRawFilter(rwc.socket, solCANRaw, canRawFilter, filter)
}

func (rwc *readWriteCloser) deleteFilter() error {
	if rwc.socket == 0 {
		return ErrorKernelFilterNotSupported
	}
	return unix.SetsockoptCanRawFilter(rwc.socket, solCANRaw, canRawFilter, nil)
}
