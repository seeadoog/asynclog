### async log base on [zap](https://github.com/uber-go/zap) and [lumberjack](https://github.com/natefinch/lumberjack)



### example 

````
    lg, err := NewLogger(&LogConf{
		Level:    "error",
		Filename: "./test.log",
		// Caller:   true,
        //copy log to stdout
		ExtraWriters: []io.Writer{os.Stdout},
	})
	if err != nil {
		panic(err)
	}
	slog := lg.Sugar()

	slog.Errorw("error", "dd", "x", "mm", map[string]any{"xx": "xx"})

````