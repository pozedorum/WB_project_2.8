package main

import (
	"fmt"

	"github.com/pozedorum/WB_project_2/task11/anagrams"
)

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(anagrams.Anagrams(input))
}
