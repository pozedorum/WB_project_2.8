package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"task16/internal/models"
	"task16/internal/storage"
)

type Downloader struct {
	client  *http.Client
	storage *storage.Storage
}

func NewDownloader(storage *storage.Storage) *Downloader {
	return &Downloader{
		client:  &http.Client{Timeout: time.Second * 30},
		storage: storage,
	}
}

func (d *Downloader) Download(ctx context.Context, u *url.URL) (*models.Resource, bool, error) {
	// start := time.Now()
	fmt.Printf("Downloading: %s\n", u.String()) // 🟢

	// defer func() {
	// 	fmt.Printf("Download finished: %s (took %v)\n", u.String(), time.Since(start)) // 🟢
	// }()

	if rs, ok := d.storage.Get(u.String()); ok {
		return rs, true, nil
	}

	content, contentType, err := d.GetContent(ctx, u)
	if err != nil {
		return nil, false, err
	}

	rs := storage.NewResource(u, content, contentType)

	if err := d.storage.Save(rs); err != nil {
		return nil, false, err
	}
	return rs, false, nil
}

func (d *Downloader) GetContent(ctx context.Context, u *url.URL) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("create request failed: %w", err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return content, resp.Header.Get("Content-Type"), nil
}
