package asynclog

import (
	"errors"
	"io"
	"net"
)

type logstashWriter struct {
	conn    net.Conn
	dialer  func() (net.Conn, error)
	err     error
	onError func([]byte, error)
}

// async tcp writer
func NewTcpWriter(addrs []string, onError func(msg []byte, err error)) io.Writer {

	return NewTcpWriterWithDialer(func() (conn net.Conn, err error) {
		if len(addrs) == 0 {
			return nil, errors.New("no address for tcp writer")
		}
		for _, addr := range addrs {
			conn, err = net.Dial("tcp", addr)
			if err == nil {
				return conn, nil
			}
		}
		return
	}, onError)
}

func NewTcpWriterWithDialer(dialer func() (net.Conn, error), onError func(msg []byte, err error)) io.Writer {

	if onError == nil {
		panic("onError must not be nil")
	}
	w := &logstashWriter{
		dialer:  dialer,
		onError: onError,
	}
	return AsyncWriter(w)
}

func (w *logstashWriter) checkConn() net.Conn {
	if w.conn != nil {
		return w.conn
	}
	//var err error
	w.conn, w.err = w.dialer()
	return w.conn
}

func (w *logstashWriter) Write(p []byte) (n int, err error) {
	conn := w.checkConn()
	if conn == nil {
		w.onError(p, w.err)
		return 0, w.err
	}
	n, err = conn.Write(p)
	if err != nil {
		conn.Close()
		w.conn = nil
		w.onError(p, err)
	}
	return n, err
}

func (w *logstashWriter) Close() error {
	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}
