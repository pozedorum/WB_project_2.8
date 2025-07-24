package parcer

import (
	"testing"
)

func TestTokenizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		// Базовые случаи
		{
			name:     "simple command",
			input:    "ls -l",
			expected: []string{"ls", "-l"},
		},
		{
			name:     "command with quotes",
			input:    `echo "hello world"`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "single quoted command",
			input:    `echo 'hello $USER'`,
			expected: []string{"echo", "hello $USER"},
		},

		// Тесты с экранированием
		{
			name:     "escaped space",
			input:    `echo hello\ world`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "escaped quote",
			input:    `echo \"hello\"`,
			expected: []string{"echo", `"hello"`},
		},

		// Тесты с операторами
		{
			name:     "pipe operator",
			input:    "ls | grep go",
			expected: []string{"ls", "|", "grep", "go"},
		},
		{
			name:     "double pipe no spaces",
			input:    "ls||grep",
			expected: []string{"ls", "||", "grep"},
		},
		{
			name:     "redirect with append",
			input:    "echo hello >> file.txt",
			expected: []string{"echo", "hello", ">>", "file.txt"},
		},

		// Комбинированные случаи
		{
			name:     "combined operators",
			input:    "ls -l | grep test && echo found || echo not found",
			expected: []string{"ls", "-l", "|", "grep", "test", "&&", "echo", "found", "||", "echo", "not", "found"},
		},
		{
			name:     "complex quotes and operators",
			input:    `echo "hello" > file.txt && cat << EOF`,
			expected: []string{"echo", "hello", ">", "file.txt", "&&", "cat", "<<", "EOF"},
		},

		// Ошибочные случаи
		{
			name:    "unclosed double quote",
			input:   `echo "hello`,
			wantErr: true,
		},
		{
			name:    "unclosed single quote",
			input:   `echo 'hello`,
			wantErr: true,
		},
		{
			name:     "empty input",
			input:    "",
			expected: []string{},
		},
		{
			name:     "multiple spaces",
			input:    "   ls    -l   ",
			expected: []string{"ls", "-l"},
		},
		{
			name:     "special chars in quotes",
			input:    `echo "| && > <"`,
			expected: []string{"echo", "| && > <"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenizeString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenizeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareStringSlices(got, tt.expected) {
				t.Errorf("TokenizeString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Вспомогательная функция для сравнения слайсов строк
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
