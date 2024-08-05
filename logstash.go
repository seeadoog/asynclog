package asynclog

import (
	"io"
	"net"
)

type logstashWriter struct {
	conn   net.Conn
	dialer func() (net.Conn, error)
	err    error
}

// async tcp writer
func NewTcpWriter(addr string) io.Writer {
	w := &logstashWriter{
		dialer: func() (net.Conn, error) {
			return net.Dial("tcp", addr)
		},
	}
	return AsyncWriter(w)
}

func (w *logstashWriter) checkConn() net.Conn {
	if w.conn != nil {
		return w.conn
	}
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
		return 0, w.err
	}
	n, err = conn.Write(p)
	if err != nil {
		conn.Close()
		w.conn = nil
	}
	return n, err
}

func (w *logstashWriter) Close() error {
	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}
