package sortpkg

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/pozedorum/WB_project_2/task10/pkg/options"
)

var monthOrder map[string]int

func init() {
	monthOrder = map[string]int{
		"jan": 1, "january": 1, "feb": 2, "february": 2,
		"mar": 3, "march": 3, "apr": 4, "april": 4,
		"may": 5, "jun": 6, "june": 6, "jul": 7, "july": 7,
		"aug": 8, "august": 8, "sep": 9, "september": 9,
		"oct": 10, "october": 10, "nov": 11, "november": 11,
		"dec": 12, "december": 12,
	}
}

type SortStruct struct {
	lines []string
	fs    options.FlagStruct
}

func MakeSortStruct(lines []string, fs options.FlagStruct) *SortStruct {
	return &SortStruct{lines, fs}
}

func (ss *SortStruct) StringsSort() error {
	var sortErr error

	sort.SliceStable(ss.lines, func(i, j int) bool {
		if sortErr != nil {
			return false
		}

		res1, err := getKey(ss.fs, ss.lines[i])
		if err != nil {
			sortErr = err
			return false
		}

		res2, err := getKey(ss.fs, ss.lines[j])
		if err != nil {
			sortErr = err
			return false
		}

		if res1 == "0" {
			return true
		}
		if res2 == "0" {
			return false
		}

		if *ss.fs.RFlag {
			return res1 > res2
		}

		return res1 < res2
	})

	return sortErr
}

// getKey возвращает ключ для сортировки на основе флагов
// Формат ключа:
//   - Числа: 20-значное число с ведущими нулями
//   - Месяцы: двузначный номер месяца (01-12)
//   - Ошибка: если флаги конфликтуют
//   - "0", если число 0, или строка без возможности сортировки

func getKey(fs options.FlagStruct, str string) (string, error) {
	if (*fs.HFlag || *fs.NFlag) && *fs.MFlag {
		return "", fmt.Errorf("flags -n/-h and -M are mutually exclusive")
	}

	parts := strings.Fields(str)

	if *fs.KFlag < 1 || *fs.KFlag > len(parts) {
		return "0", nil
	}

	resPart := parts[*fs.KFlag-1]

	if *fs.BFlag {
		resPart = strings.TrimRightFunc(resPart, unicode.IsSpace)
	}

	switch {
	case *fs.HFlag:
		return parseHumanNumber(resPart), nil
	case *fs.NFlag:
		return parseNumber(resPart), nil
	case *fs.MFlag:
		return parseMonth(resPart), nil
	default:
		return resPart, nil
	}
}

func parseNumber(resPart string) string {
	if num, err := strconv.ParseFloat(resPart, 64); err != nil {
		return "0"
	} else {
		return fmt.Sprintf("%020.0f", num)
	}
}

func parseHumanNumber(resPart string) string {
	if len(resPart) == 0 {
		return "0"
	}
	numStr := resPart
	multiplier := 1.0
	lastChar := strings.ToUpper(resPart[len(resPart)-1:])
	if strings.ContainsAny(lastChar, "KMGT") {
		numStr = strings.TrimRight(resPart, "KMGT")
		switch lastChar {
		case "K":
			multiplier = 1e3
		case "M":
			multiplier = 1e6
		case "G":
			multiplier = 1e9
		case "T":
			multiplier = 1e12
		}
	}
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return "0"
	}
	return fmt.Sprintf("%020.0f", num*multiplier)
}

func parseMonth(resPart string) string {
	if num, ok := monthOrder[strings.ToLower(resPart)]; ok {
		return fmt.Sprintf("%02d", num) // Форматируем как 01, 02, ... 12
	} else {
		return "0"
	}
}

// log.Printf("warning: -c %s file is empty", filepath)
func isSorted(filepath string, fs options.FlagStruct) (bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return false, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		// Пустой файл считается отсортированным
		return true, nil
	}

	prevLine := scanner.Text()
	prevKey, err := getKey(fs, prevLine)
	if err != nil {
		return false, err
	}
	for scanner.Scan() {
		currentLine := scanner.Text()
		currentKey, err := getKey(fs, currentLine)
		if err != nil {
			return false, err
		}
		// Сравниваем ключи с учетом флага -r
		if *fs.RFlag {
			if currentKey > prevKey {
				return false, nil
			}
		} else {
			if currentKey < prevKey {
				return false, nil
			}
		}

		prevKey = currentKey
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading file: %v", err)
	}

	return true, nil
}
