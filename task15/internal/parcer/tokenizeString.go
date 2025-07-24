package parcer

import (
	"fmt"
	"strings"
	"unicode"
)

// TODO: Сделать обработку системных переменных через $
func tokenizeString(str string) ([]string, error) {
	var tokens []string
	var token strings.Builder
	runes := []rune(str)

	var (
		inDoubleQuotes bool
		inSingleQuotes bool
		escapeNext     bool
	)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch {
		case escapeNext:
			token.WriteRune(r)
			escapeNext = false

		case r == '\\' && !inSingleQuotes:
			escapeNext = true

		case r == '\'' && !inDoubleQuotes:
			inSingleQuotes = !inSingleQuotes

		case r == '"' && !inSingleQuotes:
			inDoubleQuotes = !inDoubleQuotes

		// Проверка двойных спецсимволов
		case isDoubleControlSymbol(runes, i) && !inSingleQuotes && !inDoubleQuotes && !escapeNext:
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
			// Добавляем двойной символ и пропускаем следующий
			tokens = append(tokens, string(runes[i:i+2]))
			i++

		// Проверка одиночных спецсимволов
		case isControlSymbol(r) && !inSingleQuotes && !inDoubleQuotes && !escapeNext:
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
			tokens = append(tokens, string(r))

		case unicode.IsSpace(r) && !inSingleQuotes && !inDoubleQuotes && !escapeNext:
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}

		default:
			token.WriteRune(r)
		}
	}

	if token.Len() > 0 {
		tokens = append(tokens, token.String())
	}

	if inSingleQuotes || inDoubleQuotes {
		return nil, fmt.Errorf("unclosed quotes")
	}

	return tokens, nil
}

// Проверяем двойные спецсимволы
func isDoubleControlSymbol(runes []rune, pos int) bool {
	if pos+1 >= len(runes) {
		return false
	}
	current := runes[pos]
	next := runes[pos+1]

	return (current == '>' && next == '>') || // >>
		(current == '<' && next == '<') || // <<
		(current == '&' && next == '&') || // &&
		(current == '|' && next == '|') // ||
}

// Проверяем одиночные спецсимволы
func isControlSymbol(r rune) bool {
	return r == '|' || r == '>' || r == '<' || r == '&'
}
