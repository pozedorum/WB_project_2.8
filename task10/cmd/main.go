package main

import (
	"log"
	"os"

	"github.com/pozedorum/WB_project_2.8/task10/internal/options"
	"github.com/pozedorum/WB_project_2.8/task10/internal/sortpkg"
)

func main() {
	fs, args := options.ParseOptions()

	// Определяем источник ввода
	if len(args) == 1 {
		//fmt.Println("here")
		// Сортировка файла с выводом в stdout
		err := sortpkg.ExternalSortToStdout(args[0], *fs)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(args) == 0 {
		// Проверяем, есть ли данные в stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Сортировка stdin с выводом в stdout
			err := sortpkg.ProcessStdio(*fs)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			sortpkg.ProcessInteractiveInput(*fs)
			//log.Fatal("No input file or stdin data provided")
		}
	} else {
		log.Fatal("too much files ")
	}
}
