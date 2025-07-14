package sortpkg

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"

	"github.com/pozedorum/WB_project_2.8/task10/internal/options"
)

const ConstChunkSize = 10_000 // Чанки по 10k строк (тут я не знаю, насколько ограничивать размер чанка, потом подумаю)

type ExternalSortStruct struct {
	fs            options.FlagStruct
	tempFilesList []string
	inputFile     string
	outputFile    string
	chunkSize     int
}

func MakeExternalSortStruct(fs options.FlagStruct, inputFile, outputFile string, chunkSize int) *ExternalSortStruct {

	return &ExternalSortStruct{fs, make([]string, 0), inputFile, outputFile, chunkSize}
}

func ExternalSort(inputFile, outputFile string, fs options.FlagStruct) error {

	ess := MakeExternalSortStruct(fs, inputFile, outputFile, ConstChunkSize)
	if *ess.fs.CFlag {
		if isSorted(ess.inputFile, fs) {
			fmt.Println("File is sorted")
			return nil
		} else {
			fmt.Println("File is not sorted")
			return nil
		}
	}

	err := ess.splitAndSort()
	if err != nil {
		return err
	}
	return ess.mergeChunks()
}

func (ess *ExternalSortStruct) splitAndSort() error {
	file, err := os.Open(ess.inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buffer := make([]string, 0, ess.chunkSize)
	buflen := 0

	fileInd := 0
	for scanner.Scan() {
		buffer = append(buffer, scanner.Text())
		buflen++
		if buflen >= ess.chunkSize {
			ess.sortAndSaveChunk(fileInd, buffer)
			fileInd++
			buffer = buffer[:0] // очищаем буфер
			buflen = 0
		}
	}
	if len(buffer) > 0 {
		ess.sortAndSaveChunk(fileInd, buffer)

	}

	return nil
}

func (ess *ExternalSortStruct) sortAndSaveChunk(ind int, lines []string) {
	ss := MakeSortStruct(lines, ess.fs)

	ss.StringsSort()

	tmpFile := "chunk_" + strconv.Itoa(ind) + ".tmp"
	f, _ := os.Create(tmpFile)
	defer f.Close()

	for _, line := range ss.lines {
		f.WriteString(line + "\n")
	}

	ess.tempFilesList = append(ess.tempFilesList, tmpFile)
}

func (ess *ExternalSortStruct) mergeChunks() error {
	out, err := os.Create(ess.outputFile)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(out)
	defer func() {
		writer.Flush()
		out.Close()
	}()

	// Открываем все чанки
	files := make([]*os.File, len(ess.tempFilesList))
	scanners := make([]*bufio.Scanner, len(ess.tempFilesList))
	for i, chunk := range ess.tempFilesList {
		f, err := os.Open(chunk)
		if err != nil {
			return err
		}
		files[i] = f
		scanners[i] = bufio.NewScanner(f)
	}

	// Инициализация кучи
	h := &MinHeap{}
	heap.Init(h)

	// Загружаем первые строки каждого файла
	for i, sc := range scanners {
		if sc.Scan() {
			heap.Push(h, HeapItem{line: sc.Text(), index: i})
		}
	}

	var lastLine string
	firstLine := true

	for h.Len() > 0 {
		item := heap.Pop(h).(HeapItem)

		if *ess.fs.UFlag {
			if firstLine {
				lastLine = item.line
				writer.WriteString(item.line + "\n")
				firstLine = false
			} else if item.line != lastLine {
				lastLine = item.line
				writer.WriteString(item.line + "\n")
			}
		} else {
			writer.WriteString(item.line + "\n")
		}

		// Продвигаем сканер и добавляем следующую строку в кучу
		if scanners[item.index].Scan() {
			nextLine := scanners[item.index].Text()
			// Если включён -u, пропускаем строки, которые уже равны lastLine
			if !(*ess.fs.UFlag && nextLine == lastLine) {
				heap.Push(h, HeapItem{
					line:  nextLine,
					index: item.index,
				})
			}
		}
	}

	// Закрытие файлов
	for _, f := range files {
		f.Close()
		os.Remove(f.Name())
	}

	return nil
}
