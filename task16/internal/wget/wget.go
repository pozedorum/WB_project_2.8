package wget

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"task16/config"
	"task16/internal/downloader"
	"task16/internal/parser"
	"task16/internal/storage"
)

type linkWithDepth struct {
	u     *url.URL
	depth int
}

func NewLinkWithDepth(u *url.URL, depth int) *linkWithDepth {
	return &linkWithDepth{u: u, depth: depth}
}

func NewListOfLinks(uList []*url.URL, depth int) []*linkWithDepth {
	if uList == nil {
		return nil
	}
	res := make([]*linkWithDepth, 0, len(uList))
	for _, u := range uList {
		res = append(res, NewLinkWithDepth(u, depth))
	}
	return res
}

func Wget(urlFrom string, depth int) error {
	cnf := config.DefaultConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cnf.Timeout)
	defer cancel()
	setupSignalHandling(cancel)

	u, err := downloader.ValidateURL(urlFrom)
	if err != nil {
		return err
	}
	baseURL, err := downloader.GetBaseDomain(u)
	if err != nil {
		return err
	}

	// Инициализация компонентов
	store := storage.NewStorage(baseURL, baseURL, u)
	dl := downloader.NewDownloader(store)
	p := parser.NewParser(u)

	taskQueue := make(chan *linkWithDepth, 2000)
	resQueue := make(chan error, 2000)
	allGood := make(chan struct{})

	defer closeChannels(taskQueue, resQueue)

	var taskWg sync.WaitGroup
	taskQueue <- NewLinkWithDepth(u, 0)
	taskWg.Add(1)

	for i := 0; i < cnf.WorkersCount; i++ {
		go worker(ctx, dl, p, depth, taskQueue, resQueue, &taskWg)
	}
	go warden(ctx, &taskWg, allGood)
	// Начальная задача

	var firstErr error
	errorCount := 0
	for {
		select {
		case err := <-resQueue:
			if err != nil {
				errorCount++
				if firstErr == nil {
					firstErr = err
					cancel()
				}
				fmt.Printf("Error processing URL: %v\n", err)
			}
		case <-allGood:
			if errorCount > 0 {
				fmt.Printf("Completed with %d errors\n", errorCount)
			}
			return firstErr
		case <-ctx.Done():
			fmt.Println("returning by context")
			return firstErr

		}
	}
}

func setupSignalHandling(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
		<-sig // Второй сигнал - немедленный выход
		os.Exit(1)
	}()
}

// Тут я воспользовался идеей с or, но наоборот, жду завершения всех воркеров
func warden(ctx context.Context, taskWg *sync.WaitGroup, allGood chan<- struct{}) {
	done := make(chan struct{})
	go func() {
		taskWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		close(allGood)
	case <-ctx.Done():
		return
	}
}

func worker(ctx context.Context, dl *downloader.Downloader, p *parser.Parser,
	maxDepth int, queue chan *linkWithDepth, results chan<- error, taskWg *sync.WaitGroup,
) {
	for {
		select {
		case <-ctx.Done():
			// fmt.Println("Worker canceled via ctx")
			return
		case task, ok := <-queue:
			// fmt.Printf("Processing URL (depth %d/%d): %s\n",
			// 	newURL.depth, maxDepth, newURL.u.String())
			if !ok {
				return
			}
			processTask(ctx, dl, p, task, maxDepth, queue, results, taskWg)

		}
	}
	// воркер может несколько раз подряд обрабатывать ссылки одного слоя рекурсии, нужна отдельная функция следящая за уровнем рекурсии
}

func processTask(ctx context.Context, dl *downloader.Downloader, p *parser.Parser,
	task *linkWithDepth, maxDepth int, queue chan<- *linkWithDepth, results chan<- error, taskWg *sync.WaitGroup,
) {
	defer taskWg.Done() // Уменьшаем счётчик при завершении обработки

	// Скачивание и обработка
	content, alreadyDownloaded, err := dl.Download(ctx, task.u)
	if err != nil {
		results <- err
		return
	}
	if alreadyDownloaded || task.depth >= maxDepth {
		return
	}

	// Парсинг и добавление новых задач
	links, err := p.ExtractLinks(content.Content)
	if err != nil {
		results <- err
		return
	}

	for _, link := range NewListOfLinks(links, task.depth+1) {
		taskWg.Add(1) // Увеличиваем счётчик перед добавлением

		select {
		case queue <- link:
		default:
			taskWg.Done() // Если очередь переполнена
			results <- fmt.Errorf("queue full, dropped: %s", link.u.String())
		}
	}
}

func closeChannels(taskQueue chan *linkWithDepth, resQueue chan error) {
	close(taskQueue)
	close(resQueue)
}
