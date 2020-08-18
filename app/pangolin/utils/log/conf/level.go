package conf

import (
	"fmt"
	"strings"
)

type Level int8

const (
	LevelFatal   Level = iota // 0, Fatal level. Logs and then calls `logger.Exit(1)`.
	LevelError                // 1, Error level.
	LevelWarning              // 2, Warning level.
	LevelInfo                 // 3, Info level.
	LevelDebug                // 4, Debug level.

	_minLevel = LevelDebug
	_maxLevel = LevelFatal
)

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "fatal":
		return LevelFatal, nil
	case "error":
		return LevelError, nil
	case "warning":
		return LevelWarning, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid Level: %v", lvl)
}

// implement TextMarshaler in encoding package
func (level *Level) UnmarshalText(text []byte) error {
	l, err := ParseLevel(string(text))
	if err != nil {
		return err
	}

	*level = Level(l)

	return nil
}

// implement TextUnmarshaler in encoding package
func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case LevelDebug:
		return []byte("debug"), nil
	case LevelInfo:
		return []byte("info"), nil
	case LevelWarning:
		return []byte("warning"), nil
	case LevelError:
		return []byte("error"), nil
	case LevelFatal:
		return []byte("fatal"), nil
	}

	return nil, fmt.Errorf("not a valid level %v", level)
}

func isValidLevel(level Level) bool {
	return level >= _minLevel && level <= _maxLevel
}
