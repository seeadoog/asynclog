package asynclog

import (
	"bufio"
	"expvar"
	"io"
	"sync"
)

type asyncRotate struct {
	w   WriteBuffer
	buf chan []byte
	bp  sync.Pool
}

// 同步io 转异步 io
// 达到同步io的即时写入到文件的效果，同时不阻塞写入，但是buffer 满了时会丢失日志
func newAsyncRotate(w io.Writer) io.Writer {
	bw, ok := w.(WriteBuffer)
	if !ok {
		bw = bufio.NewWriterSize(w, 1024*320)
	}
	ar := &asyncRotate{w: bw, buf: make(chan []byte, 100000)}
	go ar.run()
	return ar
}

func AsyncWriter(w io.Writer) io.Writer {
	return newAsyncRotate(w)
}

func (a *asyncRotate) getBf() []byte {
	b, ok := a.bp.Get().([]byte)
	if ok {
		return b
	}
	return make([]byte, 256)
}

func (a *asyncRotate) run() {

	for {
		select {
		case buf := <-a.buf:
			a.w.Write(buf)
			a.bp.Put(buf)

		default:
			a.w.Flush()

			buf := <-a.buf
			a.w.Write(buf)
			a.bp.Put(buf)
		}
	}
}

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
		a.bp.Put(b)
		logBufferNIl.Add(1)
	}
	return len(p), nil
}
