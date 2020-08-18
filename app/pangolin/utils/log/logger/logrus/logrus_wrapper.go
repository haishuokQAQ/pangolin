package logrus

import (
	"context"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/syslog"
	"gopkg.in/natefinch/lumberjack.v2"

	"gitlab.p1staff.com/tsp/common/log/conf"
	"gitlab.p1staff.com/tsp/common/log/logger"
)

type Close func()

type LogrusWrapper struct {
	logrusLogger *logrus.Logger
	fields       logrus.Fields

	fnClose Close // fnClose is used to release resources when Close()
}

func NewLogrusWrapper(options *conf.Config) (*LogrusWrapper, error) {
	logrusLogger := logrus.New()

	logrusLogger.SetLevel(convertToLogrusLevel(options.Level))

	logrusLogger.SetFormatter(convertToLogrusFormatter(options.Formatter))

	writers := make([]io.Writer, 0, len(options.Outputs))
	closers := []io.Closer{}
	for _, v := range options.Outputs {
		if v.Type == conf.OutputTypeStdout {
			writers = append(writers, os.Stdout)
		}

		if v.Type == conf.OutputTypeStderr {
			writers = append(writers, os.Stderr)
		}

		if v.Type == conf.OutputTypeFile && v.File != nil {
			file, err := os.OpenFile(*v.File, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				return nil, errors.Wrapf(err, "fail to open file %v", v.File)
			}

			writers = append(writers, file)
			closers = append(closers, file)
		}

		if v.Type == conf.OutputTypeRotateFile && v.RotateFile != nil {

			v.RotateFile.SetDefaults()

			rotateFile := &lumberjack.Logger{
				Filename:   v.RotateFile.FileName,
				MaxSize:    v.RotateFile.MaxSize,
				MaxBackups: v.RotateFile.MaxBackups,
				MaxAge:     v.RotateFile.MaxAge,
				Compress:   v.RotateFile.Compress,
			}

			writers = append(writers, rotateFile)
			closers = append(closers, rotateFile)
		}

		if v.Type == conf.OutputTypeSyslog && v.Syslog != nil {

			hook, err := syslog.NewSyslogHook(v.Syslog.Protocol, v.Syslog.Address, v.Syslog.GetFacility(), "")
			if err != nil {
				return nil, errors.Wrapf(err, "fail to connect %s using %s", v.Syslog.Address, v.Syslog.Protocol)
			}
			logrusLogger.AddHook(hook)
		}
	}

	mv := io.MultiWriter(writers...)
	logrusLogger.SetOutput(mv)
	logrusLogger.SetReportCaller(true)

	wrapper := &LogrusWrapper{
		logrusLogger: logrusLogger,
	}
	wrapper.fnClose = func() {
		for _, v := range closers {
			v.Close()
		}
	}

	return wrapper, nil
}

func (wrapper *LogrusWrapper) Debugf(format string, args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Debugf(format, args...)
	} else {
		wrapper.logrusLogger.Debugf(format, args...)
	}
}

func (wrapper *LogrusWrapper) Infof(format string, args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Infof(format, args...)
	} else {
		wrapper.logrusLogger.Infof(format, args...)
	}
}

func (wrapper *LogrusWrapper) Warningf(format string, args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Warningf(format, args...)
	} else {
		wrapper.logrusLogger.Warningf(format, args...)
	}
}

func (wrapper *LogrusWrapper) Errorf(format string, args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Errorf(format, args...)
	} else {
		wrapper.logrusLogger.Errorf(format, args...)
	}
}

func (wrapper *LogrusWrapper) Fatalf(format string, args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Fatalf(format, args...)
	} else {
		wrapper.logrusLogger.Fatalf(format, args...)
	}
}

func (wrapper *LogrusWrapper) Debug(args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Debug(args...)
	} else {
		wrapper.logrusLogger.Debug(args...)
	}
}

func (wrapper *LogrusWrapper) Info(args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Info(args...)
	} else {
		wrapper.logrusLogger.Info(args...)
	}
}

func (wrapper *LogrusWrapper) Warning(args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Warning(args...)
	} else {
		wrapper.logrusLogger.Warning(args...)
	}
}

func (wrapper *LogrusWrapper) Error(args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Error(args...)
	} else {
		wrapper.logrusLogger.Error(args...)
	}
}

func (wrapper *LogrusWrapper) Fatal(args ...interface{}) {
	if len(wrapper.fields) > 0 {
		wrapper.logrusLogger.WithFields(wrapper.fields).Fatal(args...)
	} else {
		wrapper.logrusLogger.Fatal(args...)
	}
}

func (wrapper *LogrusWrapper) WithField(key string, value interface{}) logger.Logger {
	result := &LogrusWrapper{
		logrusLogger: wrapper.logrusLogger,
	}

	// 合并 wrapper.fields 和 key:value 到 data中
	// key:value 可能会覆盖 wrapper.fields 现有项
	data := make(map[string]interface{}, len(wrapper.fields)+1)
	for k, v := range wrapper.fields {
		data[k] = v
	}
	data[key] = value

	result.fields = logrus.Fields(data)
	return result
}

func (wrapper *LogrusWrapper) WithFields(fields map[string]interface{}) logger.Logger {
	result := &LogrusWrapper{
		logrusLogger: wrapper.logrusLogger,
		fnClose:      wrapper.fnClose,
	}

	// 合并 wrapper.fields 和 key:value 到 data中
	// fields 可能会覆盖 wrapper.fields 现有项
	data := make(map[string]interface{}, len(wrapper.fields)+len(fields))
	for k, v := range wrapper.fields {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}

	result.fields = logrus.Fields(data)
	return result
}

func (wrapper *LogrusWrapper) WithTraceInCtx(ctx context.Context) logger.Logger {
	result := &LogrusWrapper{
		logrusLogger: wrapper.logrusLogger,
		fnClose:      wrapper.fnClose,
	}

	if ctx != nil {
		result.fields = logrus.Fields(logger.CtxToMap(ctx))
	}
	return result
}

func (wrapper *LogrusWrapper) Flush() {
	// Refer to: https://github.com/sirupsen/logrus/issues/435
	// logrus doesn't provide any flush or sync method. If you
	// don't want to lost message, just sleep a time before exit
	wrapper.Infof("If you see this message and `Flush` is the last logger's method called in your application, it means no log lost.")
}

func (wrapper *LogrusWrapper) Close() {
	if wrapper.fnClose != nil {
		wrapper.fnClose()
	}
}

func convertToLogrusLevel(l conf.Level) logrus.Level {
	var level logrus.Level
	switch l {
	case conf.LevelDebug:
		level = logrus.DebugLevel
	case conf.LevelInfo:
		level = logrus.InfoLevel
	case conf.LevelWarning:
		level = logrus.WarnLevel
	case conf.LevelError:
		level = logrus.ErrorLevel
	case conf.LevelFatal:
		level = logrus.FatalLevel
	default:
		level = logrus.DebugLevel
	}

	return level
}

func convertToLogrusFormatter(f conf.Formatter) logrus.Formatter {
	var formatter logrus.Formatter
	switch f {
	case conf.JSONFormater:
		formatter = &logrus.JSONFormatter{}
	case conf.ConsoleFormater:
		formatter = &logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05.000"}
	default:
		formatter = &logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05.000"}
	}
	return formatter
}
