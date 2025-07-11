package task9

import (
	"fmt"
	"testing"
)

type testStruct struct {
	input  string
	output string
	err    bool
}

func TestUnpackString(t *testing.T) {
	testArr := []testStruct{
		{"abcd", "abcd", false},
		{"a4bc2d5e", "aaaabccddddde", false},
		{"45", "", true},
		{"", "", false},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
	}

	for ind, test := range testArr {
		t.Run(fmt.Sprintf("test %d", ind), func(t *testing.T) {
			output, err := UnpackString(test.input)

			if (err != nil) != test.err {
				t.Errorf("wrong error\nexpected: %v\nactual: %v", test.err, err)
			}

			if output != test.output {
				t.Errorf("wrong output string\nexpected: %s\nactual: %s", test.output, output)
			}
		})
	}
}
