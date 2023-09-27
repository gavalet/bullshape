package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *MyLogger

type MyLogger struct {
	Slogger *zap.SugaredLogger
}

func GetLogger() *MyLogger {
	if logger == nil {
		atom := zap.NewAtomicLevel()
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "datetime"
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z"))
		}

		sLogger := zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.Lock(os.Stdout),
			atom,
		), zap.AddCaller()).Sugar()
		atom.SetLevel(zap.InfoLevel)
		logger = &MyLogger{Slogger: sLogger}

	}
	// if extra != nil {
	// 	logger.Slogger = logger.Slogger.With(extra...)
	// }
	return logger
	// if extra != nil {
	// 	logger = logger.With(extra...)
	// }

	//return &MyLogger{logger: }
}
func (ml *MyLogger) SetExtra(extra ...interface{}) {
	logger.Slogger = logger.Slogger.With(extra...)
}

func NewZapSugaredLogger(w zapcore.WriteSyncer, extra ...interface{}) *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "datetime"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z"))
	}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(w),
		atom,
	), zap.AddCaller()).Sugar()
	atom.SetLevel(zap.InfoLevel)

	if extra != nil {
		logger = logger.With(extra...)
	}

	return logger
}
