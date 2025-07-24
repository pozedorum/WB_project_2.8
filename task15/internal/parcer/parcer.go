package parcer

import (
	"fmt"
)

// TODO: Добавить обработку FD редиректов (2>&1)

func ParceLine(str string) (*Command, error) {
	tokens, err := tokenizeString(str)
	if err != nil {
		return nil, err
	}
	cmd, err := parceTokens(tokens)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func parceTokens(tokens []string) (*Command, error) {
	tokensCount := len(tokens)
	var prev *Command
	var cmd *Command

	ind := 0
	for ind < tokensCount {
		current := &Command{}

		for ind < tokensCount && !isControlOperator(tokens[ind]) {
			if isRedirectOperator(tokens[ind]) {
				if ind+1 >= tokensCount {
					return nil, fmt.Errorf("error: pipeline has no file")
				}

				comRedirect := Redirect{
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

		for ind < tokensCount {
			switch tokens[ind] {
			case Pipe:
				ind++
			case And, Or:
				operator := tokens[ind]
				ind++

				if ind >= tokensCount {
					return nil, fmt.Errorf("error: after %s a command is expected", operator)
				}

				if current.AndNext != nil || current.OrNext != nil {
					return nil, fmt.Errorf("error: multiple control operators")
				}
				nextCommand, err := parceTokens(tokens[ind:])
				if err != nil {
					return nil, err
				}

				if operator == And {
					current.AndNext = nextCommand
				} else {
					current.OrNext = nextCommand
				}

				return cmd, nil
			default:
				return nil, fmt.Errorf("error: unexpected operator: %s", tokens[ind])
			}
		}
	}

	return cmd, nil
}

func isControlOperator(token string) bool {
	_, ok := controlOperators[token]
	return ok
}

func isRedirectOperator(token string) bool {
	_, ok := redirectOperators[token]
	return ok
}
