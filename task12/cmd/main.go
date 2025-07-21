package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pozedorum/WB_project_2/task12/internal/grep"
	"github.com/pozedorum/WB_project_2/task12/pkg/options"
)

func main() {
	// Парсинг флагов (остаётся таким же)
	fs, args := options.ParseOptions()

	// Определяем источник ввода
	switch len(args) {
	case 1:
		// Обработка файла (аргумент - имя файла)
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatalf("grep: %s: %v", args[0], err)
		}
		defer file.Close()

		err = grep.Grep(file, *fs, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}

	case 0:
		// Проверяем, есть ли данные в stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Обработка pipe/redirect из stdin
			err := grep.Grep(os.Stdin, *fs, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Интерактивный режим (ожидание ввода с клавиатуры)
			fmt.Println("Waiting for input (press Ctrl+D to finish):")
			err := grep.Grep(os.Stdin, *fs, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		}

	default:
		log.Fatal("grep: too many arguments")
	}
}
