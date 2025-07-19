package sortpkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pozedorum/WB_project_2/task10/pkg/options"
)

func ProcessStdio(fs options.FlagStruct) error {
	// Создаем временный файл для входных данных
	tmpInput, err := os.CreateTemp("", "sort_input_*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpInput.Name())
	defer tmpInput.Close()

	// Копируем stdin во временный файл
	if _, err = io.Copy(tmpInput, os.Stdin); err != nil {
		return err
	}
	tmpInput.Close() // Закрываем, чтобы убедиться в записи

	// Создаем временный файл для результатов
	tmpOutput, err := os.CreateTemp("", "sort_output_*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpOutput.Name())
	defer tmpOutput.Close()

	// Выполняем сортировку
	if err = ExternalSort(tmpInput.Name(), tmpOutput.Name(), fs); err != nil {
		return err
	}

	// Копируем результат в stdout
	if _, err = tmpOutput.Seek(0, 0); err != nil {
		return err
	} // Перематываем на начало

	_, err = io.Copy(os.Stdout, tmpOutput)
	return err
}

func ExternalSortToStdout(inputFile string, fs options.FlagStruct) error {
	// Создаем временный файл для результатов
	tmpOutput, err := os.CreateTemp("", "sort_output_*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpOutput.Name())
	defer tmpOutput.Close()

	// Выполняем сортировку
	if err = ExternalSort(inputFile, tmpOutput.Name(), fs); err != nil {
		return err
	}

	// Копируем результат в stdout
	if _, err = tmpOutput.Seek(0, 0); err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, tmpOutput)
	return err
}

func ProcessInteractiveInput(fs options.FlagStruct) {
	reader := bufio.NewReader(os.Stdin)
	var lines []string

	// Выводим подсказку
	os.Stderr.WriteString("Enter text to sort (Ctrl+D to finish):\n")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // EOF или другая ошибка
		}
		lines = append(lines, strings.TrimSuffix(line, "\n"))
	}

	// Сортируем и выводим
	ss := MakeSortStruct(lines, fs)
	ss.StringsSort()

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for _, line := range ss.lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			fmt.Printf("ExternalSortStruct.sortAndSaveChunk - file.WriteString: %v", err)
		}
	}
}
