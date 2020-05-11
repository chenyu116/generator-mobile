package logger

import (
	"os"

	"github.com/utahta/go-cronowriter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger ZapLogger
var ZapLogger *zap.Logger

// InitLogger InitLogger
func InitLogger(debug bool, logLevel string) {
	caller := ""
	if debug {
		caller = "caller"
		logLevel = "debug"
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      caller,
		MessageKey:     "msg",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	lv := zapcore.InfoLevel
	err := lv.Set(logLevel)
	if err != nil {
		lv = zapcore.InfoLevel
	}
	levelPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= lv
	})
	//lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//  return lvl < zapcore.ErrorLevel
	//})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	//consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		//zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(cronowriter.MustNew("./logs/%Y%m%d.log")), levelPriority),
	)

	ZapLogger = zap.New(core, zap.AddCaller())
}
