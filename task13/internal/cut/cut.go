package cut

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pozedorum/WB_project_2/task13/options"
)

func Cut(input io.Reader, fs options.FlagStruct, writer io.Writer) error {

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		processLine(fs, scanner.Text(), writer)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	return nil
}

func ProcessFiles(fs options.FlagStruct, args []string) {
	for _, filename := range args {
		file, err := os.Open(filename)
		if err != nil {
			log.Printf("cut: %s: %v", filename, err)
			continue
		}

		err = Cut(file, fs, os.Stdout)
		file.Close()

		if err != nil {
			log.Printf("cut: %s: %v", filename, err)
		}
	}
}

func processLine(fs options.FlagStruct, line string, writer io.Writer) {
	if !strings.ContainsRune(line, fs.Delimiter) {
		if !fs.SFlag {
			fmt.Fprintln(writer, line)
		}
		return
	}

	// Разбиваем строку на поля
	fields := strings.Split(line, string(fs.Delimiter))
	output := make([]string, 0, len(fs.Fields))

	// Обрабатываем поля в порядке, указанном пользователем
	for _, fieldNum := range fs.Fields {
		// Проверяем, что номер поля существует
		if fieldNum > 0 && fieldNum <= len(fields) {
			output = append(output, fields[fieldNum-1]) // -1 т.к. индексация с 0
		}
	}

	// Выводим результат
	//if len(output) > 0 {
	fmt.Fprintln(writer, strings.Join(output, string(fs.Delimiter)))
	//}
}
