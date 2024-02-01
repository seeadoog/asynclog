package asynclog

import (
	"fmt"
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type WriteBuffer interface {
	io.Writer
	Flush() error
}

const (
	FileNameDiscard = "null"
	FileNameStdio   = "stdout"
)

type logger struct {
	log zap.Logger
}

func newLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)

	}
	logger.Info("hello", zapcore.Field{Key: "name", String: "hello"})
}

type Logger = zap.Logger

type SugaredLogger = zap.SugaredLogger

type LogConf struct {
	Filename string `json:"filename"`
	// log level: error warn info debug panic
	Level string `json:"level"`
	// if write log sync
	Sync       bool `json:"sync"`
	MaxSize    int  `json:"max_size" default:"500"`
	MaxAge     int  `json:"max_age" default:"30"`
	MaxBackups int  `json:"max_backups" default:"30"`
	LocalTime  bool `json:"local_time"`
	Compress   bool `json:"compress"`
	Caller     bool `json:"caller"`
	CallSkip   int  `json:"call_skip"`
	//write log to this writer
	Writer io.Writer `json:"-"`
	//copy write log to other loggers
	ExtraWriters []io.Writer `json:"-"`

	ZapEncConf func(c *zapcore.EncoderConfig) error
	ZapOptions []zap.Option
}

func (lc *LogConf) init() {
	if lc.Level == "" {
		lc.Level = "info"
	}
	if lc.MaxSize <= 0 {
		lc.MaxSize = 500
	}

	if lc.MaxAge == 0 {
		lc.MaxAge = 30
	}

	if lc.MaxBackups == 0 {
		lc.MaxBackups = 30
	}

}

var (
	levelNone = zapcore.Level(-5)
)

func getLevel(s string) (zapcore.Level, error) {
	switch s {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	case "none":
		return levelNone, nil
	default:
		return 0, fmt.Errorf("invalid log level:%s", s)
	}
}

func NewLogger(lc *LogConf) (*Logger, error) {

	lc.init()

	var lw io.Writer
	if lc.Writer != nil {
		lw = lc.Writer
	} else if lc.Filename == FileNameDiscard {
		lw = io.Discard
	} else if lc.Filename == FileNameStdio {
		lw = os.Stdout
	} else {
		w := &lumberjack.Logger{
			Filename:   lc.Filename,
			MaxSize:    lc.MaxSize,
			MaxAge:     lc.MaxAge,
			MaxBackups: lc.MaxBackups,
			LocalTime:  lc.LocalTime,
			Compress:   lc.Compress,
		}
		lw = w
	}

	for _, w := range lc.ExtraWriters {
		lw = io.MultiWriter(lw, w)
	}

	if !lc.Sync {
		lw = AsyncWriter(lw)
	}

	level, err := getLevel(lc.Level)
	if err != nil {
		return nil, err
	}

	if level == levelNone {
		lw = io.Discard
	}

	opts := []zap.Option{}
	if lc.Caller {
		opts = append(opts, zap.AddCaller())
	}
	if lc.CallSkip > 0 {
		opts = append(opts, zap.AddCallerSkip(lc.CallSkip))
	}
	opts = append(opts, lc.ZapOptions...)

	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	if lc.ZapEncConf != nil {
		err := lc.ZapEncConf(&zapConfig)
		if err != nil {
			return nil, err
		}
	}

	logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zapConfig), zapcore.AddSync(lw), level), opts...)

	return logger, nil

}

func NewSugarLogger(lc *LogConf) (*zap.SugaredLogger, error) {
	lg, err := NewLogger(lc)
	if err != nil {
		return nil, err
	}
	return lg.Sugar(), nil
}
