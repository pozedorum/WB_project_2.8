package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"task16/internal/models"
	"task16/internal/storage"
	"time"
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

func (d *Downloader) Download(u *url.URL) (*models.Resource, error) {
	if rs, ok := d.storage.Get(u.String()); ok {
		return rs, nil
	}

	content, contentType, err := d.GetContent(u)
	if err != nil {
		return nil, err
	}

	rs := storage.MakeNewResource(u, content, contentType)

	if err := d.storage.Save(rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (d *Downloader) GetContent(u *url.URL) ([]byte, string, error) {
	resp, err := d.client.Get(u.String())
	if err != nil {
		return nil, "", err
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
