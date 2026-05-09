// Package sortpkg contains realisation og external sort for file
package sortpkg

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"

	"github.com/pozedorum/WB_project_2/task10/pkg/options"
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
		is, err := isSorted(ess.inputFile, fs)
		if err != nil {
			return err
		} else if is {
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
			if err := ess.sortAndSaveChunk(fileInd, buffer); err != nil {
				return err
			}
			fileInd++
			buffer = buffer[:0] // очищаем буфер
			buflen = 0
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(buffer) > 0 {
		if err := ess.sortAndSaveChunk(fileInd, buffer); err != nil {
			return err
		}
	}

	return nil
}

func (ess *ExternalSortStruct) sortAndSaveChunk(ind int, lines []string) error {
	ss := MakeSortStruct(lines, ess.fs)

	if err := ss.StringsSort(); err != nil {
		return err
	}

	tmpFileName := "chunk_" + strconv.Itoa(ind) + ".tmp"
	tmpFile, err := os.CreateTemp("", tmpFileName)
	if err != nil {
		return fmt.Errorf("error create temp chunk file: %w", err)
	}
	defer tmpFile.Close()

	for _, line := range ss.lines {
		if _, err := tmpFile.WriteString(line + "\n"); err != nil {
			fmt.Printf("ExternalSortStruct.sortAndSaveChunk - file.WriteString: %v", err)
			return err
		}
	}

	ess.tempFilesList = append(ess.tempFilesList, tmpFile.Name())
	return nil
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

	h := &MinHeap{
		fs: ess.fs,
	}
	heap.Init(h)

	for i, sc := range scanners {
		if sc.Scan() {
			line := sc.Text()

			key, err := getKey(ess.fs, line)
			if err != nil {
				return err
			}

			heap.Push(h, HeapItem{
				line:  line,
				key:   key,
				index: i,
			})
		}
	}

	var lastLine string
	firstLine := true

	for h.Len() > 0 {
		item := heap.Pop(h).(HeapItem)

		if *ess.fs.UFlag {
			if firstLine {
				lastLine = item.line

				if _, err := writer.WriteString(item.line + "\n"); err != nil {
					return err
				}

				firstLine = false
			} else {
				itemKey, err := getKey(ess.fs, item.line)
				if err != nil {
					return err
				}

				lastKey, err := getKey(ess.fs, lastLine)
				if err != nil {
					return err
				}

				if itemKey != lastKey {
					lastLine = item.line

					if _, err := writer.WriteString(item.line + "\n"); err != nil {
						return err
					}
				}
			}
		} else {
			if _, err := writer.WriteString(item.line + "\n"); err != nil {
				return err
			}
		}

		// Продвигаем сканер и добавляем следующую строку в кучу
		if scanners[item.index].Scan() {
			nextLine := scanners[item.index].Text()
			if !(*ess.fs.UFlag && nextLine == lastLine) {
				key, err := getKey(ess.fs, nextLine)
				if err != nil {
					return err
				}

				heap.Push(h, HeapItem{
					line:  nextLine,
					key:   key,
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
