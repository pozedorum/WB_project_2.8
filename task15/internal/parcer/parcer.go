package parcer

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

	if tokensCount == 0 {
		return nil, ErrEmptyString
	}

	var (
		prev    *Command
		cmd     *Command
		current *Command
	)
	ind := 0
	for ind < tokensCount {
		current = &Command{}

		for ind < tokensCount && !isControlOperator(tokens[ind]) {
			if isRedirectOperator(tokens[ind]) {
				if ind+1 >= tokensCount {
					return nil, ErrNoFileForRedirect
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

		if ind >= tokensCount {
			break
		}

		if current.IsEmpty() {
			return nil, ErrEmptyCommand
		}
		switch tokens[ind] {
		case Pipe:
			ind++

		case And, Or:
			operator := tokens[ind]
			ind++

			if ind >= tokensCount {
				return nil, ErrMissingAfterOperator
			}

			if current.AndNext != nil || current.OrNext != nil {
				return nil, ErrMultipleOperators
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
			return nil, ErrUnexpectedOperator
		}
	}

	return cmd, nil
}

func (c *Command) IsEmpty() bool {
	return c.Name == "" &&
		len(c.Args) == 0 &&
		len(c.Redirects) == 0 &&
		c.PipeTo == nil &&
		c.AndNext == nil &&
		c.OrNext == nil
}

func isControlOperator(token string) bool {
	_, ok := controlOperators[token]
	return ok
}

func isRedirectOperator(token string) bool {
	_, ok := redirectOperators[token]
	return ok
}
