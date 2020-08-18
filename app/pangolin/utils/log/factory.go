package log

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.p1staff.com/tsp/common/log/conf"
	"gitlab.p1staff.com/tsp/common/log/logger"
	logrusWrapper "gitlab.p1staff.com/tsp/common/log/logger/logrus"
	zapWrapper "gitlab.p1staff.com/tsp/common/log/logger/zap"
)

var std logger.Logger

func Init(c *Config) error {
	logger, err := New(c)
	if err != nil {
		return err
	}

	std = logger
	return nil
}

func New(c *Config) (logger.Logger, error) {
	if c.Core == conf.ZapCore {
		re, err := zapWrapper.NewZapWrapper(c)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to create zap wrapper")
		}

		return re, nil
	}

	if c.Core == conf.LogrusCore {
		re, err := logrusWrapper.NewLogrusWrapper(c)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to create logrus wrapper")
		}

		return re, nil
	}

	return nil, errors.New("core is not supported")
}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}
func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

func Debug(args ...interface{}) {
	std.Debug(args...)
}
func Info(args ...interface{}) {
	std.Info(args...)
}
func Warning(args ...interface{}) {
	std.Warning(args...)
}
func Error(args ...interface{}) {
	std.Error(args...)
}
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

func WithField(key string, value interface{}) logger.Logger {
	return std.WithField(key, value)
}
func WithFields(fields map[string]interface{}) logger.Logger {
	return std.WithFields(fields)
}
func WithTraceInCtx(ctx context.Context) logger.Logger {
	return std.WithTraceInCtx(ctx)
}

// Flushing any buffered log entries. Applications should take care to call Flush before exiting.
func Flush() {
	std.Flush()
}

// Release resources
func Close() {
	std.Close()
}

func StandardLogger() logger.Logger {
	return std
}
