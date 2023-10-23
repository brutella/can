package can

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
)

const CAN_RAW_JOIN_FILTER = 0x6

func NewReadWriteCloserForInterface(i *net.Interface) (ReadWriteCloser, error) {
	s, _ := syscall.Socket(syscall.AF_CAN, syscall.SOCK_RAW, unix.CAN_RAW)
	addr := &unix.SockaddrCAN{Ifindex: i.Index}
	if err := unix.Bind(s, addr); err != nil {
		return nil, err
	}

	f := os.NewFile(uintptr(s), fmt.Sprintf("fd %d", s))

	return &readWriteCloser{rwc: f, socket: s}, nil
}

func (rwc *readWriteCloser) setPassFilter(allowedIds []uint32) error {
	if rwc.socket == 0 {
		return ErrorKernelFilterNotSupported
	}

	if len(allowedIds) >= unix.CAN_RAW_FILTER_MAX {
		return ErrorKernelFilterTooMany
	}

	filter := make([]unix.CanFilter, len(allowedIds))

	for i, allowedId := range allowedIds {
		filter[i].Id = allowedId
		filter[i].Mask = unix.CAN_SFF_MASK
	}

	return unix.SetsockoptCanRawFilter(rwc.socket, unix.SOL_CAN_RAW, unix.CAN_RAW_FILTER, filter)
}

func (rwc *readWriteCloser) setBlockFilter(disallowedIds []uint32) error {
	if rwc.socket == 0 {
		return ErrorKernelFilterNotSupported
	}

	filter := make([]unix.CanFilter, len(disallowedIds))

	for i, disallowedId := range disallowedIds {
		filter[i].Id = disallowedId | unix.CAN_INV_FILTER
		filter[i].Mask = unix.CAN_SFF_MASK
	}

	return unix.SetsockoptCanRawFilter(rwc.socket, unix.SOL_CAN_RAW, unix.CAN_RAW_JOIN_FILTER, filter)
}

func (rwc *readWriteCloser) deleteFilter() error {
	if rwc.socket == 0 {
		return ErrorKernelFilterNotSupported
	}
	return unix.SetsockoptCanRawFilter(rwc.socket, unix.SOL_CAN_RAW, unix.CAN_RAW_FILTER, nil)
}
