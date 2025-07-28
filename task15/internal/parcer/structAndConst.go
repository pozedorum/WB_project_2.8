package parcer

import "errors"

// Константы ошибок
var (
	ErrEmptyString          = errors.New("empty command string")
	ErrNoFileForRedirect    = errors.New("no file specified for redirect")
	ErrUnexpectedOperator   = errors.New("unexpected operator")
	ErrMissingAfterOperator = errors.New("command expected after operator")
	ErrMultipleOperators    = errors.New("multiple control operators")
	ErrEmptyCommand         = errors.New("empty command before operator")
)
