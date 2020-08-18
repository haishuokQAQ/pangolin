package log

import "pangolin/app/pangolin/utils/log/conf"

type Config = conf.Config
type Core = conf.Core
type Formatter = conf.Formatter
type Level = conf.Level
type Output = conf.Output

const (
	ZapCore    conf.Core = conf.ZapCore
	LogrusCore           = conf.LogrusCore
)

const (
	JSONFormater    conf.Formatter = conf.JSONFormater
	ConsoleFormater                = conf.ConsoleFormater
)

const (
	LevelFatal   conf.Level = conf.LevelFatal
	LevelError              = conf.LevelError
	LevelWarning            = conf.LevelWarning
	LevelInfo               = conf.LevelInfo
	LevelDebug              = conf.LevelDebug
)

const (
	OutputTypeStdout     string = conf.OutputTypeStdout
	OutputTypeStderr            = conf.OutputTypeStderr
	OutputTypeFile              = conf.OutputTypeFile
	OutputTypeRotateFile        = conf.OutputTypeRotateFile
	OutputTypeSyslog            = conf.OutputTypeSyslog
)
