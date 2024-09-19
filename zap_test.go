package asynclog

import (
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"log/slog"

	"go.uber.org/zap/zapcore"
)

// elasticsearch
func Test_newLogger(t *testing.T) {

}

func BenchmarkLog(b *testing.B) {
	// TODO: Initialize
	log := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{}))
	log.Debug("name", "name", "1")

	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		log.Error("error", "name", "chenjian", "age", 5, "som", "hello world")

	}

}

func Test_Slog(t *testing.T) {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	log.Error("name", "name", "1")

}

func Test_Logger(t *testing.T) {
	lg, err := NewLogger(&LogConf{
		Level:    "error",
		Filename: "./test.log",
		// Caller:   true,
		ExtraWriters: []io.Writer{NewTcpWriter("127.0.0.1:9000")},

		ZapEncConf: func(c *zapcore.EncoderConfig) error {
			c.LevelKey = "lv"
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	slog := lg.Sugar()

	slog.Errorw("error", "dd", "x", "mm", map[string]any{"xx": "xx"})
	slog.Errorw("error", "dd", "x", "mm", map[string]any{"xx": "xx"})
	slog.Errorw("error", "dd", "x", "mm", map[string]any{"xx": "xx"})
	slog.Errorw("error", "dd", "x", "mm", map[string]any{"xx": "xx"})
	time.Sleep(1 * time.Second)
}

func BenchmarkLogger(b *testing.B) {
	// TODO: Initialize
	lg, err := NewLogger(&LogConf{
		Level:    "none",
		Filename: "../test.log",
		// Filename: FileNameDiscard,
		// Caller:   true,
		Sync: false,
	})

	if err != nil {
		panic(err)
	}

	slog := lg.Sugar()

	slog.Errorw("error", "dd", "x")
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		slog.Errorw("error", "name", "chenjian", "age", 5, "som", "hello worldxlpolgmyjgroojsjdofdsjfds99dsfsdnfndsjfjdsojfosdjfojdsjf")
	}
	fmt.Println(logBufferNIl.String())
}

func Test_LogWatch(t *testing.T) {
	ew := []io.Writer{os.Stdout}
	lw, err := NewLogWatch[SugaredLogger](&LogConf{
		Level: "error",
		// Filename: "../test.log",
		Filename: FileNameDiscard,
		// Caller:   true,
		Sync:         true,
		ExtraWriters: ew,
	}, NewSugarLogger)
	if err != nil {
		panic(err)
	}

	lw.Get().Infow("hello world 1 ")

	err, ok := lw.Update(&LogConf{
		Level:    "info",
		Filename: "./test.log",
		// Filename: FileNameDiscard,
		// Caller:   true,
		Sync:         true,
		ExtraWriters: ew,
	})

	fmt.Println(err, ok)
	fmt.Println()
	lw.Get().Infow("hello world 2")

	err, ok = lw.Update(&LogConf{
		Level:    "info",
		Filename: "./test2.log",
		// Filename: FileNameDiscard,
		// Caller:   true,
		Sync:         true,
		ExtraWriters: ew,
	})

	fmt.Println(err, ok)

	lw.Get().Infow("helo3 ")

}

func TestTcpServer(t *testing.T) {
	ls, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer ls.Close()
	for {
		conn, err := ls.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			defer conn.Close()
			buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					return
				}
				fmt.Println(string(buf[:n]))
			}
		}()
	}
}

func BenchmarkINT(b *testing.B) {
	// TODO: Initialize

	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
	}
}
