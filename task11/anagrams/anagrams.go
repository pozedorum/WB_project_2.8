package anagrams

import (
	"sort"
	"strings"
)

// anagrams группирует слова по анаграммам.
// В результирующую map попадают только группы,
// содержащие более одного слова.
func Anagrams(strs []string) map[string][]string {
	groups := make(map[string][]string)

	for _, str := range strs {
		word := strings.ToLower(str)

		runes := []rune(word)
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})

		key := string(runes)
		groups[key] = append(groups[key], word)
	}

	res := make(map[string][]string)

	for _, group := range groups {
		if len(group) < 2 {
			continue
		}

		sort.Strings(group)
		res[group[0]] = group
	}

	return res
}
