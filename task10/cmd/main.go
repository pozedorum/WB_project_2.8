package main

import (
	"fmt"

	"github.com/pozedorum/WB_project_2.8/task10/internal/options"
)

func main() {
	fs, args := options.ParceOptions()
	fs.PrintFlags()
	fmt.Println(args)
}
