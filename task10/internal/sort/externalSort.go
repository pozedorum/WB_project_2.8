package sort

import (
	"bufio"
	"os"
	"strconv"

	"github.com/pozedorum/WB_project_2.8/task10/internal/options"
)

type ExternalSortStruct struct {
	fs         options.FlagStruct
	filesList  []string
	outputFile string
	chunkSize  int
}

func MakeExternalSortStruct(fs options.FlagStruct, outputFile string, chunkSize int) *ExternalSortStruct {

	return &ExternalSortStruct{fs, make([]string, 0), outputFile, chunkSize}
}

func ExternalSort(inputFile, outputFile string, fs options.FlagStruct) error {

	chunks, err := splitAndSort(inputFile, 100_000, fs) // Чанки по 100k строк (тут я не знаю, насколько ограничивать размер чанка, потом подумаю)
	if err != nil {
		return err
	}
	return mergeChunks(chunks, outputFile)
}

func (ess *ExternalSortStruct) splitAndSort(filepath string, chunkSize int, fs options.FlagStruct) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(os.Stdin)
	buffer := make([]string, 0, chunkSize)
	buflen := 0

	fileInd := 0
	for scanner.Scan() {
		buffer = append(buffer, scanner.Text())
		buflen++
		if buflen >= chunkSize {
			chunkFile := sortAndSaveChunk(fileInd, buffer, fs)
			fileInd++
			filesList = append(filesList, chunkFile)
			buffer = buffer[:0] // очищаем буфер
			buflen = 0
		}
	}
	if len(buffer) > 0 {
		chunkFile := sortAndSaveChunk(fileInd, buffer, fs)
		filesList = append(filesList, chunkFile)
	}

	return
}

func (ess *ExternalSortStruct) sortAndSaveChunk(ind int, lines []string, fs options.FlagStruct) string {
	ss := MakeSortStruct(lines, fs)

	ss.StringsSort()

	tmpFile := "chunk_" + strconv.Itoa(ind) + ".tmp"
	f, _ := os.Create(tmpFile)
	defer f.Close()

	for _, line := range ss.lines {
		f.WriteString(line + "\n")
	}
	return tmpFile
}

func (ess *ExternalSortStruct) mergeChunks(filesList []string, output string) error {

	return nil
}
