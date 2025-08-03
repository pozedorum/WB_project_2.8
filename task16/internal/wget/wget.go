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
	store := storage.NewStorage(cnf.ResultDir, baseURL, u)
	dl := downloader.NewDownloader(store)
	p := parser.NewParser(u)

	taskQueue := make(chan *url.URL, 100)
	resQueue := make(chan error, 100)
	doneList := make([]chan struct{}, cnf.WorkersCount)
	allGood := make(chan struct{})
	defer closeChannels(taskQueue, resQueue, doneList, allGood)

	for i := 0; i < cnf.WorkersCount; i++ {
		doneList[i] = make(chan struct{})
		go worker(ctx, dl, p, depth, taskQueue, resQueue, doneList[i])
	}
	go warden(ctx, doneList, cnf.WorkersCount, allGood)
	// Начальная задача
	taskQueue <- u

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
			return firstErr

		}
	}
	return nil
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
func warden(ctx context.Context, channels []chan struct{}, workersCount int, allGood chan<- struct{}) {
	counter := 0
	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan struct{}) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case <-ch:
				counter++
				if counter == workersCount {
					close(allGood)
				}
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(allGood)
	}()
}

func worker(ctx context.Context, dl *downloader.Downloader, p *parser.Parser,
	maxDepth int, queue chan *url.URL, results chan<- error, done chan<- struct{}) {
	defer func() {
		close(done)
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case newURL, ok := <-queue:
			if !ok {
				return
			}
			if err := processURL(ctx, dl, p, newURL, maxDepth, queue); err != nil {
				select {
				case results <- err:
				case <-ctx.Done():
					return
				}
			}
			maxDepth--
		}
	}
	// воркер может несколько раз подряд обрабатывать ссылки одного слоя рекурсии, нужна отдельная функция следящая за уровнем рекурсии
}

func processURL(ctx context.Context, dl *downloader.Downloader, p *parser.Parser,
	u *url.URL, maxDepth int, queue chan<- *url.URL) error {

	// 1. Скачиваем страницу
	content, err := dl.Download(ctx, u)
	if err != nil {
		return err
	}

	// 2. Парсим ссылки (если не превышена глубина)
	if maxDepth > 0 {
		links, err := p.ExtractLinks(content.Content)
		if err != nil {
			return err
		}

		// 3. Добавляем новые ссылки в очередь
		for _, link := range links {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case queue <- link:
			default:
				// Если очередь переполнена, пропускаем ссылку
				fmt.Printf("Queue full, skipping URL: %s\n", link)
			}
		}
	}

	return nil
}

func closeChannels(taskQueue chan *url.URL, resQueue chan error, doneList []chan struct{}, allGood chan struct{}) {
	close(taskQueue)
	close(resQueue)
	for _, ch := range doneList {
		close(ch)
	}
	close(allGood)
}
