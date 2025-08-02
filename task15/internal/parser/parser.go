package parser

import (
	"errors"
	"task15/internal/core"
)

// Константы ошибок
var (
	ErrEmptyString          = errors.New("empty command string")
	ErrNoFileForRedirect    = errors.New("no file specified for redirect")
	ErrUnexpectedOperator   = errors.New("unexpected operator")
	ErrMissingAfterOperator = errors.New("command expected after operator")
	ErrMultipleOperators    = errors.New("multiple control operators")
	ErrEmptyCommand         = errors.New("empty command before operator")
)

func ParseLine(str string) (*core.Command, error) {
	tokens, err := tokenizeString(str)
	if err != nil {
		return nil, err
	}
	cmd, err := parseTokens(tokens)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func parseTokens(tokens []string) (*core.Command, error) {
	tokensCount := len(tokens)

	if tokensCount == 0 {
		return nil, ErrEmptyString
	}

	var (
		prev    *core.Command
		cmd     *core.Command
		current *core.Command
	)
	ind := 0
	for ind < tokensCount {
		current = &core.Command{}

		for ind < tokensCount && !isControlOperator(tokens[ind]) {
			if isRedirectOperator(tokens[ind]) {
				if ind+1 >= tokensCount {
					return nil, ErrNoFileForRedirect
				}

				comRedirect := core.Redirect{
					Type: tokens[ind],
					File: tokens[ind+1],
				}
				current.Redirects = append(current.Redirects, comRedirect)
				ind += 2
			} else {
				if current.Name == "" {
					current.Name = tokens[ind]
				} else {
					current.Args = append(current.Args, tokens[ind])
				}
				ind++
			}
		}

		if cmd == nil {
			cmd = current
		} else {
			prev.PipeTo = current
		}
		prev = current

		if ind >= tokensCount {
			break
		}

		if current.IsEmpty() {
			return nil, ErrEmptyCommand
		}
		switch tokens[ind] {
		case core.Pipe:
			ind++

		case core.And, core.Or:
			operator := tokens[ind]
			ind++

			if ind >= tokensCount {
				return nil, ErrMissingAfterOperator
			}

			if current.AndNext != nil || current.OrNext != nil {
				return nil, ErrMultipleOperators
			}
			nextCommand, err := parseTokens(tokens[ind:])
			if err != nil {
				return nil, err
			}
			if operator == core.And {
				current.AndNext = nextCommand
			} else {
				current.OrNext = nextCommand
			}

			return cmd, nil
		default:
			return nil, ErrUnexpectedOperator
		}
	}

	return cmd, nil
}

func isControlOperator(token string) bool {
	_, ok := core.ControlOperators[token]
	return ok
}

func isRedirectOperator(token string) bool {
	_, ok := core.RedirectOperators[token]
	return ok
}
