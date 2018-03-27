package can

import (
	"bytes"
	"io"
	"time"
	"syscall"
)

type echoReadWriteCloser struct {
	buf         bytes.Buffer
	closed      bool
	writeSocket int
}

// NewEchoReadWriteCloser returns a ReadWriteCloser which echoes received bytes
// via an Unix domain socket pair.
func NewEchoReadWriteCloser() ReadWriteCloser {
	pair, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_DGRAM, 0)
	if err != nil {
		panic(err)
	}
	return NewReadWriteCloser(&echoReadWriteCloser{writeSocket: pair[0]}, pair[1])
}

func (rw *echoReadWriteCloser) Read(b []byte) (n int, err error) {
	for {
		if rw.buf.Len() > 0 {
			return rw.buf.Read(b)
		}

		if rw.closed == true {
			break
		}

		<-time.After(time.Millisecond * 1)
	}

	return 0, io.EOF
}

func (rw *echoReadWriteCloser) Write(b []byte) (n int, err error) {
	err = syscall.Sendmsg(rw.writeSocket, b, nil, nil, 0)
	n = len(b)

	return
}

func (rw *echoReadWriteCloser) Close() error {
	rw.closed = true
	return nil
}
