package main

import (
	"log"
	"os"

	"github.com/pozedorum/WB_project_2/task13/internal/cut"
	"github.com/pozedorum/WB_project_2/task13/options"
)

func main() {
	fs, args := options.ParseOptions()

	if len(args) > 0 {
		log.Fatal("cut: too many arguments (this version supports only stdin)")
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Режим pipe/redirect
		err := cut.Cut(os.Stdin, *fs, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Интерактивный режим
		log.Fatal("cut: no input provided (use pipe or redirect)")
	}
}
