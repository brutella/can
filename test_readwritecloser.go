package can

import (
	"bytes"
)

type echoReadWriteCloser struct {
	buf bytes.Buffer
}

// NewEchoReadWriteCloser returns a ReadWriteCloser which echoes received bytes.
func NewEchoReadWriteCloser() ReadWriteCloser {
	return NewReadWriteCloser(&echoReadWriteCloser{})
}

func (rw *echoReadWriteCloser) Read(b []byte) (n int, err error)  { return rw.buf.Read(b) }
func (rw *echoReadWriteCloser) Write(b []byte) (n int, err error) { return rw.buf.Write(b) }
func (rw *echoReadWriteCloser) Close() error                      { return nil }
