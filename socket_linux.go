package can

import (
	"syscall"
)

func NewSockaddr(proto uint16, Ifindex int) syscall.Sockaddr {
	return &syscall.SockaddrLinklayer{
		Protocol: proto,
		Ifindex:  Ifindex,
	}
}
