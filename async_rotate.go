package asynclog

import (
	"bufio"
	"expvar"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type asyncRotate struct {
	w         WriteBuffer
	buf       chan []byte
	bp        sync.Pool
	onLogLost func([]byte)
}

type logOptions struct {
	onLogLost          func([]byte)
	writerBufferSize   int
	maxPendingMessages int
}

type logOption func(*logOptions)

// 同步io 转异步 io
// 达到同步io的即时写入到文件的效果，同时不阻塞写入，但是buffer 满了时会丢失日志
func newAsyncRotate(w io.Writer, opts ...logOption) io.Writer {
	opt := new(logOptions)
	for _, o := range opts {
		o(opt)
	}
	if opt.maxPendingMessages <= 0 {
		opt.maxPendingMessages = 10000
	}
	if opt.writerBufferSize <= 0 {
		opt.writerBufferSize = 1024 * 320
	}

	bw, ok := w.(WriteBuffer)
	if !ok {
		bw = bufio.NewWriterSize(w, opt.writerBufferSize)
	}
	ar := &asyncRotate{w: bw, buf: make(chan []byte, opt.maxPendingMessages), onLogLost: opt.onLogLost}
	go ar.run()
	return ar
}

func SetWriterBufferSize(writerBufferSize int) logOption {
	return func(o *logOptions) {
		o.writerBufferSize = writerBufferSize
	}
}

func SetWriterMaxPendingMessages(maxPendingMessages int) logOption {
	return func(o *logOptions) {
		o.maxPendingMessages = maxPendingMessages
	}
}

func SetWriterOnLogLost(onLogLost func([]byte)) logOption {
	return func(o *logOptions) {
		o.onLogLost = onLogLost
	}
}

func AsyncWriter(w io.Writer, opts ...logOption) io.Writer {
	return newAsyncRotate(w, opts...)
}

func (a *asyncRotate) getBf() []byte {
	b, ok := a.bp.Get().([]byte)
	if ok {
		return b
	}
	return make([]byte, 256)
}

func (a *asyncRotate) run() {

	tick := time.NewTicker(1 * time.Second)
	for {
		select {
		case buf := <-a.buf:
			a.w.Write(buf)
			a.bp.Put(buf)

		case <-tick.C:
			a.w.Flush()

		}
	}
}

var (
	flushCount = atomic.Int64{}
)

var (
	logBufferNIl = expvar.NewInt("log_write_buffer_full")
)

func (a *asyncRotate) Write(p []byte) (n int, err error) {
	b := a.getBf()
	b = append(b[:0], p...)
	select {
	case a.buf <- b:
	default:
		// 没有写进去，释放buf
		logBufferNIl.Add(1)
		if a.onLogLost != nil {
			a.onLogLost(b)
		}
		a.bp.Put(b)
	}
	return len(p), nil
}
