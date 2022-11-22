package can

import (
	"fmt"
	"net"
)

func NewReadWriteCloserForInterface(i *net.Interface) (ReadWriteCloser, error) {
	return nil, fmt.Errorf("Binding to can interface no supported on Windows")
}

func (rwc *readWriteCloser) setPassFilter(allowedIds []uint32) error {
	return ErrorKernelFilterNotSupported
}

func (rwc *readWriteCloser) setBlockFilter(disallowedIds []uint32) error {
	return ErrorKernelFilterNotSupported
}

func (rwc *readWriteCloser) deleteFilter() error {
	return ErrorKernelFilterNotSupported
}
