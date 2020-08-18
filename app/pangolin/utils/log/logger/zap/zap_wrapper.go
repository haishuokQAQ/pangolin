package zap

import (
	"context"
	"fmt"
	"log/syslog"
	"os"

	"github.com/pkg/errors"
	"gitlab.p1staff.com/tsp/common/log/conf"
	"gitlab.p1staff.com/tsp/common/log/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapWrapper struct {
	zapLogger *zap.Logger
}

const (
	defaultLevel   = zapcore.InfoLevel
	jsonEncoder    = "json"
	consoleEncoder = "console"
)

func NewZapWrapper(options *conf.Config) (*ZapWrapper, error) {
	level := convertToZapLevel(options.Level)
	encoder := convertToZapEncoder(options.Formatter)

	syncers := make([]zapcore.WriteSyncer, 0)

	// outputs
	for _, v := range options.Outputs {
		if v.Type == conf.OutputTypeStdout {
			syncers = append(syncers, zapcore.AddSync(os.Stdout))
		}

		if v.Type == conf.OutputTypeStderr {
			syncers = append(syncers, zapcore.AddSync(os.Stderr))
		}

		if v.Type == conf.OutputTypeFile && v.File != nil {
			file, err := os.OpenFile(*v.File, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				return nil, errors.Wrapf(err, "fail to open file %v", v.File)
			}
			syncers = append(syncers, zapcore.AddSync(file))
		}

		if v.Type == conf.OutputTypeRotateFile && v.RotateFile != nil {

			v.RotateFile.SetDefaults()

			hook := lumberjack.Logger{
				Filename:   v.RotateFile.FileName,   // file to store log file
				MaxSize:    v.RotateFile.MaxSize,    // max size of each log file, unit: MB
				MaxBackups: v.RotateFile.MaxBackups, // backup count
				MaxAge:     v.RotateFile.MaxAge,     // how long to store log file, unit: days
				Compress:   v.RotateFile.Compress,
				LocalTime:  v.RotateFile.LocalTime,
			}
			syncers = append(syncers, zapcore.AddSync(&hook))
		}

		if v.Type == conf.OutputTypeSyslog && v.Syslog != nil {
			w, err := syslog.Dial(v.Syslog.Protocol, v.Syslog.Address, v.Syslog.GetFacility(), "")
			if err != nil {
				return nil, errors.Wrapf(err, "fail to connect %s using %s", v.Syslog.Address, v.Syslog.Protocol)
			}
			syncers = append(syncers, zapcore.AddSync(w))
		}
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(syncers...),
		level,
	)
	return &ZapWrapper{
		zapLogger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.FatalLevel)),
	}, nil
}

func (wrapper *ZapWrapper) Debugf(format string, args ...interface{}) {
	wrapper.zapLogger.Debug(fmt.Sprintf(format, args...))
}

func (wrapper *ZapWrapper) Infof(format string, args ...interface{}) {
	wrapper.zapLogger.Info(fmt.Sprintf(format, args...))
}

func (wrapper *ZapWrapper) Warningf(format string, args ...interface{}) {
	wrapper.zapLogger.Warn(fmt.Sprintf(format, args...))
}

func (wrapper *ZapWrapper) Errorf(format string, args ...interface{}) {
	wrapper.zapLogger.Error(fmt.Sprintf(format, args...))
}

func (wrapper *ZapWrapper) Fatalf(format string, args ...interface{}) {
	wrapper.zapLogger.Fatal(fmt.Sprintf(format, args...))
}

func (wrapper *ZapWrapper) Debug(args ...interface{}) {
	wrapper.zapLogger.Debug(fmt.Sprint(args...))
}

func (wrapper *ZapWrapper) Info(args ...interface{}) {
	wrapper.zapLogger.Info(fmt.Sprint(args...))
}

func (wrapper *ZapWrapper) Warning(args ...interface{}) {
	wrapper.zapLogger.Warn(fmt.Sprint(args...))
}

func (wrapper *ZapWrapper) Error(args ...interface{}) {
	wrapper.zapLogger.Error(fmt.Sprint(args...))
}

func (wrapper *ZapWrapper) Fatal(args ...interface{}) {
	wrapper.zapLogger.Fatal(fmt.Sprint(args...))
}

func (wrapper *ZapWrapper) WithField(key string, value interface{}) logger.Logger {
	return &ZapWrapper{
		zapLogger: wrapper.zapLogger.With(zap.Any(key, value)),
	}
}

func (wrapper *ZapWrapper) WithFields(fields map[string]interface{}) logger.Logger {
	zapFields := make([]zapcore.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return &ZapWrapper{
		zapLogger: wrapper.zapLogger.With(zapFields...),
	}
}

func (wrapper *ZapWrapper) WithTraceInCtx(ctx context.Context) logger.Logger {
	if ctx != nil {
		return wrapper.WithFields(logger.CtxToMap(ctx))
	}
	return wrapper
}

func (wrapper *ZapWrapper) Flush() {
	wrapper.zapLogger.Sync()
}

func (wrapper *ZapWrapper) Close() {
	wrapper.zapLogger.Sync()
	// TODO: release resources
}

func convertToZapLevel(l conf.Level) zapcore.Level {
	var result zapcore.Level
	switch l {
	case conf.LevelDebug:
		result = zapcore.DebugLevel
	case conf.LevelInfo:
		result = zapcore.InfoLevel
	case conf.LevelWarning:
		result = zapcore.WarnLevel
	case conf.LevelError:
		result = zapcore.ErrorLevel
	case conf.LevelFatal:
		result = zapcore.FatalLevel
	default:
		result = zapcore.DebugLevel
	}
	return result
}

func convertToZapEncoder(f conf.Formatter) zapcore.Encoder {
	var encoder zapcore.Encoder

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "zaplogger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	switch f {
	case conf.ConsoleFormater:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case conf.JSONFormater:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return encoder
}
