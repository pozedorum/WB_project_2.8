package main

import (
	"log"
	"os"

	"github.com/pozedorum/WB_project_2/task10/internal/sortpkg"
	"github.com/pozedorum/WB_project_2/task10/pkg/options"
)

func main() {
	fs, args := options.ParseOptions()
	if len(args) == 1 {
		err := sortpkg.ExternalSortToStdout(args[0], *fs)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(args) == 0 {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			err := sortpkg.ProcessStdio(*fs)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := sortpkg.ProcessInteractiveInput(*fs)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal("too much files ")
	}
}
