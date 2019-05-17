package can

import (
	"fmt"
	"net"
)

func NewReadWriteCloserForInterface(i *net.Interface) (ReadWriteCloser, error) {
	return nil, fmt.Errorf("Binding to can interface no supported on Darwin")
}
