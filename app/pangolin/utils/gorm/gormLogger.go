package gorm

import (
	"pangolin/app/pangolin/utils/log"
	"pangolin/app/pangolin/utils/log/logger"
)

type GormLogger struct {
	logger logger.Logger
	level  log.Level
}

// 目前只支持 info 和 debug 级别
func (gl *GormLogger) Print(args ...interface{}) {
	if gl.level == log.LevelInfo {
		gl.logger.Info(args...)
	} else {
		gl.logger.Debug(args...)
	}
}

func NewGormLogger(logger logger.Logger) *GormLogger {
	return &GormLogger{
		logger: logger,
		level:  log.LevelDebug,
	}
}

func NewGormLoggerWithLevel(logger logger.Logger, level log.Level) *GormLogger {
	return &GormLogger{
		logger: logger,
		level:  level,
	}
}
