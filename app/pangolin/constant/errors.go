package constant

import (
	"errors"
)

const (
	ErrorCodeDuplicatedPort   = -1
	ErrorCodeDuplicatedTunnel = -2
)

var (
	errorMap = map[int]error{
		ErrorCodeDuplicatedPort:   errors.New("Duplicate port!"),
		ErrorCodeDuplicatedTunnel: errors.New("Duplicate tunnel!"),
	}
	unknownError = errors.New("Unknown error!")
)

func GetErrorByErrorCode(code int) error {
	if err, ok := errorMap[code]; ok {
		return err
	}
	return unknownError
}
