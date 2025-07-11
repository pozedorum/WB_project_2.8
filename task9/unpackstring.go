package task9

import (
	"fmt"
	"strings"
	"unicode"
)

// UnpackString take string and for each pair "letter, digit" make sequence of digit letters
func UnpackString(input string) (string, error) {
	inputRune := []rune(input)
	inputLen := len(inputRune)
	bldr := strings.Builder{}

	var letter rune = 0
	for ind := 0; ind < inputLen; ind++ {
		if inputRune[ind] == '\\' {
			if ind == inputLen-1 {
				return "", fmt.Errorf("wrong string format")
			}
			if letter != 0 {
				bldr.WriteRune(letter)
			}
			ind++
			letter = inputRune[ind]
		} else if !unicode.IsDigit(inputRune[ind]) {
			if letter != 0 {
				_, err := bldr.WriteRune(letter)
				if err != nil {
					return "", err
				}
			}
			letter = inputRune[ind]
		} else if unicode.IsDigit(inputRune[ind]) {
			if letter == 0 {
				return "", fmt.Errorf("wrong string format")
			}

			for range int(inputRune[ind] - '0') {
				bldr.WriteRune(letter)
			}
			letter = 0
		} else {
			bldr.WriteRune(inputRune[ind])
			letter = 0
		}
	}
	if letter != 0 {
		bldr.WriteRune(letter)
	}
	return bldr.String(), nil
}
