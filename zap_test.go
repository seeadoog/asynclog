package asynclog

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_newLogger(t *testing.T) {
	logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		StacktraceKey: "sk",
	}), zapcore.AddSync(os.Stdout), zapcore.DebugLevel), zap.AddStacktrace(zapcore.ErrorLevel))

	logger.Fatal("tar")

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
		ExtraWriters: []io.Writer{os.Stdout},

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
	time.Sleep(1 * time.Second)
}

func BenchmarkLogger(b *testing.B) {
	// TODO: Initialize
	lg, err := NewLogger(&LogConf{
		Level:    "error",
		Filename: "../test.log",
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
