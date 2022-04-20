package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var BaseLogger *zap.Logger
var Logger *zap.SugaredLogger

func LogInit(level string, format string) error {
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err // and leave the default level
	}
	cfg := zap.Config{
		Level:             l,
		Development:       false,
		DisableCaller:     true,
		Encoding:          format,
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		DisableStacktrace: true,
	}
	if format == "console" {
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	BaseLogger, err = cfg.Build()
	if err != nil {
		return err
	}

	Logger = BaseLogger.Sugar()
	return nil
}

func DriverLog(fields ...interface{}) error {
	LogText := "Driver Log"
	if len(fields)%2 != 0 {
		panic("driverLog expects even number of arguments")
	}
	Logger.Debugw(LogText, fields...)
	return nil
}

func DriverTrace(fields ...interface{}) error {
	LogText := "Driver Trace"
	if len(fields)%2 != 0 {
		panic("driverTrace expects even number of arguments")
	}
	Logger.Infow(LogText, fields...)
	return nil
}

func IsLevelEnabled(level string) bool {
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return false
	}

	return BaseLogger.Core().Enabled(l.Level())
}
