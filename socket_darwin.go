package can

import "syscall"

func NewSockaddr(proto uint16, Ifindex int) syscall.Sockaddr {
	// TODO(brutella) needs implementation
	return &syscall.SockaddrUnix{}
}
