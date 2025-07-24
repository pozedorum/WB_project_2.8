// main.go
package main

import (
	"log"
	"os"

	"github.com/pozedorum/WB_project_2/task13/internal/cut"
	"github.com/pozedorum/WB_project_2/task13/options"
)

func main() {
	fs, args := options.ParseOptions()

	// Если нет аргументов - читаем из stdin
	if len(args) == 0 {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			err := cut.Cut(os.Stdin, *fs, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("cut: no input provided (use file argument or pipe/redirect)")
		}
		return
	} else {
		cut.ProcessFiles(*fs, args)
	}

}
