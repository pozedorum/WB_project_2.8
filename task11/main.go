package main

import (
	"fmt"
	"sort"
)

func anagrams(strs []string) map[string][]string {
	res := make(map[string][]string)
	preRes := make(map[string][]string)

	for _, str := range strs {
		runeStr := []rune(str)
		sort.Slice(runeStr, func(i, j int) bool { return runeStr[i] < runeStr[j] })
		newStr := string(runeStr)
		if _, ok := preRes[newStr]; !ok {
			preRes[newStr] = make([]string, 0, 1)
			preRes[newStr] = append(preRes[newStr], str)
		} else {
			preRes[newStr] = append(preRes[newStr], str)
		}
	}

	for _, value := range preRes {
		if len(value) > 1 {
			res[value[0]] = value
		}

	}

	return res
}

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(anagrams(input))
}
