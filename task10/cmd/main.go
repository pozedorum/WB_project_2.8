package main

import (
	"fmt"
	"strings"

	"github.com/pozedorum/WB_project_2.8/task10/internal/options"
	"github.com/pozedorum/WB_project_2.8/task10/internal/sortpkg"
)

func main() {
	fs, args := options.ParseOptions()
	fs.PrintFlags()

	if len(args) != 1 {
		panic("error: expected path/to/file or \"|\" to enter strings in STDIN")
	}

	//fmt.Println(args)

	inputFile := args[0]
	outputFile, _ := strings.CutSuffix(inputFile, ".txt")
	outputFile += "_sorted.txt"
	if err := sortpkg.ExternalSort(inputFile, outputFile, *fs); err != nil {
		fmt.Printf("error: %v", err)
	} else {
		fmt.Printf("result of sort locales by path: %s", outputFile)
	}

}
