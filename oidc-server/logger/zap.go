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

		// default writer for logger
		consoleDebugging := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		// set log level to writer
		core := zapcore.NewTee(
			// zapcore.NewCore(consoleEncoder, consoleDebugging, zap.DebugLevel),
			zapcore.NewCore(consoleEncoder, consoleDebugging, zap.InfoLevel),
			zapcore.NewCore(consoleEncoder, consoleErrors, zap.WarnLevel),
			zapcore.NewCore(consoleEncoder, consoleErrors, zap.ErrorLevel),
		)

		// add function caller and stack trace on error
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return logger

}
