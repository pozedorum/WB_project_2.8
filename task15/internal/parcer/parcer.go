package parcer

import (
	"fmt"
	"strings"
	"unicode"
)

func ParceLine(str string) error {
	words, err := TokenizeString(str)
	if err != nil {
		return err
	}

}

func TokenizeString(str string) ([]string, error) {
	var token strings.Builder
	res := make([]string, 0, strings.Count(str, " "))

	var (
		isIntoDoubleQuotes = false
		isIntoSingleQuotes = false
		isShielded         = false
	)
	for _, rn := range []rune(str) {
		switch {
		case isShielded:
			token.WriteRune(rn)
			isShielded = false
		case rn == '\\' && !isIntoSingleQuotes:
			isShielded = true
		case rn == '\'' && !isIntoDoubleQuotes:
			isIntoSingleQuotes = !isIntoSingleQuotes
		case rn == '"' && !isIntoSingleQuotes:
			isIntoDoubleQuotes = !isIntoDoubleQuotes
		case isControlSymbol(rn) && !isIntoSingleQuotes && !isIntoDoubleQuotes && !isShielded:
			if token.Len() > 0 {
				res = append(res, token.String())
				token.Reset()
			}
			res = append(res, string(rn)) // Добавляем сам спецсимвол как отдельный токен
		case unicode.IsSpace(rn) && !isIntoSingleQuotes && !isIntoDoubleQuotes && !isShielded:
			if token.Len() > 0 {
				res = append(res, token.String())
				token.Reset()
			}
		default:
			token.WriteRune(rn)
		}
	}
	if isIntoSingleQuotes || isIntoDoubleQuotes {
		return nil, fmt.Errorf("unclosed quotes in input")
	}
	if token.Len() > 0 {
		res = append(res, token.String())
	}
	return res, nil
}

func isControlSymbol(r rune) bool {

}
