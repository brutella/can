package can

import "syscall"

// NewSockaddr returns a socket address based on the protocol and interface index.
// TODO(brutella) This method has no implementation.
func NewSockaddr(proto uint16, Ifindex int) syscall.Sockaddr {
	return &syscall.SockaddrUnix{}
}
