package sortpkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

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

// res1 := getKey(ss.fs, ss.lines[i])
// 		res2 := getKey(ss.fs, ss.lines[j])

func (ss *SortStruct) StringsSort() {
	sort.SliceStable(ss.lines, func(i, j int) bool {
		res1 := getKey(ss.fs, ss.lines[i])
		res2 := getKey(ss.fs, ss.lines[j])
		if res1 == "0" {
			return true
		}
		if res2 == "0" {
			return false
		}
		if *ss.fs.RFlag {
			return res1 >= res2
		}
		return res1 < res2 // Устойчивая сортировка
	})
}

// getKey возвращает ключ для сортировки на основе флагов
// Формат ключа:
//   - Числа: 20-значное число с ведущими нулями
//   - Месяцы: двузначный номер месяца (01-12)
//   - Ошибка: если флаги конфликтуют
//   - "0", если число 0, или строка без возможности сортировки

func getKey(fs options.FlagStruct, str string) string {
	var resPart string
	parts := strings.Fields(str)

	if *fs.KFlag < 1 || *fs.KFlag > len(parts) { // k flag -- sort by column number N
		return "0"
	} else {
		resPart = parts[*fs.KFlag-1]
	}

	if *fs.BFlag { // b flag -- ignore trailing blanks
		resPart = strings.TrimSpace(resPart)
	}

	if (*fs.HFlag || *fs.NFlag) && *fs.MFlag {
		panic("flags -n/-h and -m are mutually exclusive")
	}

	switch {
	case *fs.HFlag:
		return parseHumanNumber(resPart)
	case *fs.NFlag:
		return parseNumber(resPart)
	case *fs.MFlag:
		return parseMonth(resPart)
	default:
		return resPart
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
func isSorted(filepath string, fs options.FlagStruct) bool {
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		// Пустой файл считается отсортированным
		return true
	}

	prevLine := scanner.Text()
	prevKey := getKey(fs, prevLine)

	for scanner.Scan() {
		currentLine := scanner.Text()
		currentKey := getKey(fs, currentLine)

		// Сравниваем ключи с учетом флага -r
		if *fs.RFlag {
			if currentKey > prevKey {
				return false
			}
		} else {
			if currentKey < prevKey {
				return false
			}
		}

		prevKey = currentKey
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file: %v", err)
		return false
	}

	return true
}
