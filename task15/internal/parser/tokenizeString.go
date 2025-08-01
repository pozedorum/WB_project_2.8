package parser

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type tokenWithQuotes struct {
	content   string
	inSingleQ bool
	inDoubleQ bool
}

// TODO: Сделать обработку системных переменных через $
// TODO: Опционально добить обработку оставшихся токенов из structAndConst
func tokenizeString(str string) ([]string, error) {
	var tokens []tokenWithQuotes
	var currentToken strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escapeNext := false
	strRunes := []rune(str)

	for i := 0; i < len(strRunes); i++ {
		r := strRunes[i]

		switch {
		case escapeNext:
			if r == '$' {
				currentToken.WriteRune('\\')
			}
			currentToken.WriteRune(r)
			escapeNext = false

		case r == '\\':
			if inSingleQuote {
				// В одинарных кавычках сохраняем обратный слеш
				currentToken.WriteRune(r)
			} else {
				// Вне кавычек или в двойных - экранируем следующий символ
				escapeNext = true
			}

		case r == '\'' && !inDoubleQuote && !escapeNext:
			if inSingleQuote {
				tokens = append(tokens, tokenWithQuotes{
					content:   currentToken.String(),
					inSingleQ: true,
				})
				currentToken.Reset()
			}
			inSingleQuote = !inSingleQuote

		case r == '"' && !inSingleQuote && !escapeNext:
			if inDoubleQuote {
				tokens = append(tokens, tokenWithQuotes{
					content:   currentToken.String(),
					inDoubleQ: true,
				})
				currentToken.Reset()
			}
			inDoubleQuote = !inDoubleQuote

		case isDoubleOperator(strRunes, i) && !inSingleQuote && !inDoubleQuote && !escapeNext:
			if currentToken.Len() > 0 {
				tokens = append(tokens, tokenWithQuotes{
					content: currentToken.String(),
				})
				currentToken.Reset()
			}
			tokens = append(tokens, tokenWithQuotes{
				content: string(strRunes[i : i+2]),
			})
			i++ // Пропускаем следующий символ

		case isSingleOperator(r) && !inSingleQuote && !inDoubleQuote && !escapeNext:
			if currentToken.Len() > 0 {
				tokens = append(tokens, tokenWithQuotes{
					content: currentToken.String(),
				})
				currentToken.Reset()
			}
			tokens = append(tokens, tokenWithQuotes{
				content: string(r),
			})

		case unicode.IsSpace(r) && !inSingleQuote && !inDoubleQuote:
			if currentToken.Len() > 0 {
				tokens = append(tokens, tokenWithQuotes{
					content: currentToken.String(),
				})
				currentToken.Reset()
			}

		default:
			currentToken.WriteRune(r)
			escapeNext = false
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, tokenWithQuotes{
			content:   currentToken.String(),
			inSingleQ: inSingleQuote,
			inDoubleQ: inDoubleQuote,
		})
	}

	if inSingleQuote || inDoubleQuote {
		return nil, fmt.Errorf("unclosed quotes")
	}

	return expandEnvVars(tokens), nil
}

func isDoubleOperator(runes []rune, pos int) bool {
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

func isSingleOperator(r rune) bool {
	return r == '|' || r == '>' || r == '<' || r == '&'
}

func expandEnvVars(tokens []tokenWithQuotes) []string {
	var result []string

	for _, token := range tokens {
		content := token.content

		switch {
		case token.inSingleQ:
			// Оставляем как есть без изменений
			result = append(result, content)

		case token.inDoubleQ:
			// Раскрываем переменные внутри двойных кавычек
			result = append(result, os.ExpandEnv(content))
		case strings.Contains(content, "\\$"):
			content = strings.ReplaceAll(content, "\\$", "$")
			result = append(result, content)
		case strings.Contains(content, "$"):
			// Раскрываем переменные в незакавыченных токенах
			expanded := os.ExpandEnv(content)
			if expanded != content {
				result = append(result, strings.Fields(expanded)...)
			} else {
				result = append(result, content)
			}

		default:
			result = append(result, content)
		}
	}

	return result
}
