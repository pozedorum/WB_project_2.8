// Package options содержит стуктуру FlagStruct, используется для парсинга флагов командной строки
package options

import (
	"fmt"
	"os"
	"sort"
	"strings"

	flag "github.com/spf13/pflag"
)

type FlagStruct struct {
	Fields    []int // Те же номера полей но в виде массива чисел
	Delimiter rune  // Разделитель полей (аналог -d в cut)
	SFlag     bool  // Только строки с разделителем (аналог -s в cut)
}

// ParseOptions парсит флаги командной строки
func ParseOptions() (*FlagStruct, []string) {
	var fs FlagStruct

	FFlag := flag.StringP("fields", "f", "", "Select only these fields (columns)\n"+
		"Specify as comma-separated list or ranges (e.g. 1,3-5)")

	DFlag := flag.StringP("delimiter", "d", "\t", "Use specified delimiter instead of TAB")

	SFlag := flag.BoolP("separated", "s", false,
		"Only output lines containing delimiter")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -f 1,3-5 -d ',' file.csv\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  cat file.txt | %s -f 2\n", os.Args[0])
	}

	flag.Parse()
	fs.SFlag = *SFlag
	// Валидация обязательных флагов
	if *FFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: field list is required (-f)\n")
		flag.Usage()
		os.Exit(1)
	}
	if err := fs.ParseFields(*FFlag); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}

	runeDelimiter := []rune(*DFlag)
	if len(runeDelimiter) != 1 {
		fmt.Fprintf(os.Stderr, "Error: the delimiter must be a single character\n")
		flag.Usage()
		os.Exit(1)
	}
	fs.Delimiter = runeDelimiter[0]

	return &fs, flag.Args()
}

// ParseFields парсит строку с номерами полей в массив чисел
func (fs *FlagStruct) ParseFields(FFlag string) error {
	seen := make(map[int]bool)
	parts := strings.Split(FFlag, ",")

	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return fmt.Errorf("invalid range: %s", part)
			}

			start, err := parseInt(rangeParts[0])
			if err != nil {
				return err
			}

			end, err := parseInt(rangeParts[1])
			if err != nil {
				return err
			}

			if start > end {
				return fmt.Errorf("invalid range: start > end in %s", part)
			}

			for i := start; i <= end; i++ {
				if !seen[i] {
					fs.Fields = append(fs.Fields, i)
					seen[i] = true
				}
			}
		} else {
			num, err := parseInt(part)
			if err != nil {
				return err
			}
			if !seen[num] {
				fs.Fields = append(fs.Fields, num)
				seen[num] = true
			}
		}
	}
	sort.Ints(fs.Fields)
	return nil
}

// parseInt преобразует строку в число с проверкой ошибок
func parseInt(s string) (int, error) {
	num := 0
	_, err := fmt.Sscanf(s, "%d", &num)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}
	if num < 1 {
		return 0, fmt.Errorf("field numbers must be positive")
	}
	return num, nil
}

func (fs *FlagStruct) PrintFlags() {
	fmt.Println("flag F (fields) -", fs.Fields)
	fmt.Println("flag D (delimiter) -", fs.Delimiter)
	fmt.Println("flag S (separated) -", fs.SFlag)
}
