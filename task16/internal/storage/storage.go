package storage

import (
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
	os.Mkdir(baseDir, 0755)
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

func (s *Storage) Save(rs *models.Resource) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := rs.URL.String()
	_, ok := s.resources[key]
	if ok {
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
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, rs.Content, 0644)
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

func MakeNewResource(u *url.URL, content []byte, contentType string) *models.Resource {
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
	if u.Path == "" || strings.HasSuffix(u.Path, "/") {
		path += "index.html"
	}
	return filepath.Clean(path)
}
