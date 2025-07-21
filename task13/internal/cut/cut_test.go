package cut

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pozedorum/WB_project_2/task13/options"
)

func TestCutFunctionality(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fs       options.FlagStruct
		expected string
	}{
		{
			name:  "basic fields selection",
			input: "a:b:c:d:e",
			fs: options.FlagStruct{
				Fields:    []int{2, 4},
				Delimiter: ':',
				SFlag:     false,
			},
			expected: "b:d\n",
		},
		{
			name:  "unordered fields",
			input: "1:2:3:4:5",
			fs: options.FlagStruct{
				Fields:    []int{5, 1},
				Delimiter: ':',
				SFlag:     false,
			},
			expected: "5:1\n",
		},
		{
			name:  "empty fields",
			input: "a::c:d:",
			fs: options.FlagStruct{
				Fields:    []int{2, 4},
				Delimiter: ':',
				SFlag:     false,
			},
			expected: ":d\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			input := strings.NewReader(tt.input)
			err := Cut(input, tt.fs, &buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestFieldOrder(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fields   []int
		expected string
	}{
		{
			name:     "reverse_order",
			input:    "1:2:3:4:5",
			fields:   []int{5, 1},
			expected: "5:1\n",
		},
		{
			name:     "repeted_fields",
			input:    "a:b:c:d:e",
			fields:   []int{2, 2, 4},
			expected: "b:b:d\n",
		},
		{
			name:     "fields out of range",
			input:    "x:y:z",
			fields:   []int{1, 5, 2},
			expected: "x:y\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			fs := options.FlagStruct{
				Fields:    tt.fields,
				Delimiter: ':',
				SFlag:     false,
			}
			processLine(fs, tt.input, &buf)
			if buf.String() != tt.expected {
				t.Errorf("Ожидалось %q, получено %q", tt.expected, buf.String())
			}
		})
	}
}
