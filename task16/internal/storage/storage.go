package storage

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"task16/internal/models"
)

type Storage struct {
	root      *models.Resource            // Корневая страница (например, index.html)
	baseDir   string                      // Базовая директория для сохранения
	mu        sync.Mutex                  // Для потокобезопасности
	resources map[string]*models.Resource // Мапа всех ресурсов [URL → *Resource]
}

func NewStorage(baseDir string, rootBaseDomain string, rootURL *url.URL) *Storage {
	if err := os.Mkdir(baseDir, 0o755); err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create directory: %v", err)
	}
	return &Storage{
		resources: make(map[string]*models.Resource),
		baseDir:   baseDir,
		root: &models.Resource{
			URL:       rootURL,
			LocalPath: rootBaseDomain,
			IsHTML:    true,
		},
	}
}

func NewResource(u *url.URL, content []byte, contentType string) *models.Resource {
	rs := &models.Resource{
		URL:         u,
		LocalPath:   MakeLocalPath(u),
		Content:     content,
		ContentType: contentType,
		IsHTML:      strings.Contains(contentType, "text/html"),
	}
	return rs
}

func MakeLocalPath(u *url.URL) string {
	path := u.Path
	if path == "" || strings.HasSuffix(path, "/") {
		path += "index.html"
	} else if filepath.Ext(path) == "" {
		// Если путь не имеет расширения и не заканчивается на /, добавляем .html
		path += ".html"
	}

	// Удаляем лишние символы из имени файла
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		base += ".html"
	} else if len(ext) > 5 { // Слишком длинное расширение
		base = base[:len(base)-len(ext)] + ".html"
	}

	// Ограничиваем длину имени файла
	if len(base) > 100 {
		base = base[:100] + filepath.Ext(base)
	}

	return filepath.Join(u.Host, filepath.Dir(path), base)
}

func (s *Storage) Save(rs *models.Resource) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := rs.URL.String()
	key2, _ := strings.CutSuffix(rs.URL.String(), "/")
	_, ok := s.resources[key]
	_, ok2 := s.resources[key2]
	if ok || ok2 {
		return nil
	}
	s.resources[key] = rs
	if err := s.saveToDisk(rs); err != nil {
		delete(s.resources, key)
		return err
	}

	return nil
}

func (s *Storage) saveToDisk(rs *models.Resource) error {
	fullPath := filepath.Join(s.baseDir, rs.LocalPath)
	if _, err := os.Stat(fullPath); err == nil {
		return nil // Файл уже существует
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, rs.Content, 0o644)
}

func (s *Storage) AddLink(parent, child *models.Resource) {
	s.mu.Lock()
	defer s.mu.Unlock()

	parent.Links = append(parent.Links, child)
}

func (s *Storage) Get(urlKey string) (*models.Resource, bool) {
	s.mu.Lock()
	res, ok := s.resources[urlKey]
	s.mu.Unlock()
	return res, ok
}

func (s *Storage) Clean() error {
	return os.RemoveAll(s.baseDir)
}
