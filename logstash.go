package asynclog

import (
	"errors"
	"io"
	"net"
)

type logstashWriter struct {
	conn   net.Conn
	dialer func() (net.Conn, error)
	err    error
}

// async tcp writer
func NewTcpWriter(addrs ...string) io.Writer {

	w := &logstashWriter{
		dialer: func() (conn net.Conn, err error) {
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
		},
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
