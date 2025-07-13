package main

import (
	"fmt"

	"github.com/pozedorum/WB_project_2.8/task10/internal/options"
)

func main() {
	fs, args := options.ParceOptions()
	fs.PrintFlags()

	if len(args) != 1 {
		panic("error: expected path/to/file or \"|\" to enter strings in STDIN")
	}

	fmt.Println(args)
}
