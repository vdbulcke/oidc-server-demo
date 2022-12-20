package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetZapLogger returns a zap.Logger
func GetZapLogger(Debug bool) *zap.Logger {

	// Zap Logger
	var logger *zap.Logger
	var err error
	if Debug {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// override time format
		zapConfig := zap.NewProductionEncoderConfig()
		zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleEncoder := zapcore.NewJSONEncoder(zapConfig)

		// First, define our level-handling logic.
		errorPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl > zapcore.DebugLevel && lvl < zapcore.ErrorLevel
		})

		// default writer for logger
		consoleStdout := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		// set log level to writer
		core := zapcore.NewTee(
			// zapcore.NewCore(consoleEncoder, consoleDebugging, zap.DebugLevel),
			zapcore.NewCore(consoleEncoder, consoleStdout, infoPriority),
			zapcore.NewCore(consoleEncoder, consoleErrors, infoPriority),
			zapcore.NewCore(consoleEncoder, consoleErrors, errorPriority),
		)

		// add function caller and stack trace on error
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return logger

}
