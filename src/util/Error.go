package util

import (
	"fmt"
	"github.com/labstack/gommon/log"
)

type ErrorType string

const (
	NOT_FOUND      ErrorType = "NOT_FOUND"
	ALREADY_EXISTS           = "ALREADY_EXISTS"
	UNKNOWN_ERROR            = "UNKNOWN_ERROR"
	INVALID_INPUT            = "INVALID_INPUT"
)

func MakeError(errorType ErrorType, message string) *MarketListError {
	return &MarketListError{
		ErrorType: errorType,
		Message:   message,
	}
}

func MakeErrorUnknown(error error) *MarketListError {
	log.Error(error)
	return MakeError(UNKNOWN_ERROR, error.Error())
}

type MarketListError struct {
	ErrorType ErrorType `json:"Error"`
	Message   string    `json:"Message"`
}

func (m *MarketListError) Error() string {
	return fmt.Sprintf("%s: %s", m.Error, m.Message)
}
