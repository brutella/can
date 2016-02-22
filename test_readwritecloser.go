package can

import (
	"bytes"
)

type echoReadWriteCloser struct {
	buf bytes.Buffer
}

func (rw *echoReadWriteCloser) Read(b []byte) (n int, err error)  { return rw.buf.Read(b) }
func (rw *echoReadWriteCloser) Write(b []byte) (n int, err error) { return rw.buf.Write(b) }
func (rw *echoReadWriteCloser) Close() error                      { return nil }

func NewEchoReadWriteCloser() ReadWriteCloser {
	return NewReadWriteCloser(&echoReadWriteCloser{})
}
