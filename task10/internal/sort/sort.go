package sort

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/pozedorum/WB_project_2.8/task10/internal/options"
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

func (ss *SortStruct) StringsSort() {
	sort.SliceStable(ss.lines, func(i, j int) bool {
		res1 := ss.getKey(i)
		res2 := ss.getKey(j)
		if res1 == "0" {
			return true
		}
		if res2 == "0" {
			return false
		}
		if *ss.fs.RFlag {
			return res1 > res2
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
func (ss *SortStruct) getKey(strInd int) string {

	var resPart string
	parts := strings.Fields(ss.lines[strInd])

	if *ss.fs.KFlag < 1 || *ss.fs.KFlag > len(parts) { // k flag -- sort by column number N
		log.Printf("Warning: -k %d is out of range (max %d), using first column\n",
			*ss.fs.KFlag, len(parts))
		resPart = parts[0]
	} else {
		resPart = parts[*ss.fs.KFlag-1]
	}

	if *ss.fs.BFlag { // b flag -- ignore trailing blanks
		resPart = strings.TrimSpace(resPart)
	}

	if (*ss.fs.HFlag || *ss.fs.NFlag) && *ss.fs.MFlag {
		panic("flags -n/-h and -m are mutually exclusive")
	}

	switch {
	case *ss.fs.HFlag:
		return parseHumanNumber(resPart)
	case *ss.fs.NFlag:
		return parseNumber(resPart)
	case *ss.fs.MFlag:
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
	numStr := strings.TrimRight(resPart, "KMGT")
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return "0"
	}
	switch strings.ToUpper(resPart[len(resPart)-1:]) {
	case "K":
		num *= 1e3
	case "M":
		num *= 1e6
	case "G":
		num *= 1e9
	case "T":
		num *= 1e12
	}
	return fmt.Sprintf("%020.0f", num)
}

func parseMonth(resPart string) string {
	if num, ok := monthOrder[strings.ToLower(resPart)]; ok {
		return fmt.Sprintf("%02d", num) // Форматируем как 01, 02, ... 12
	} else {
		return "0"
	}
}
