package storage

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestStorage(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "storage_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	baseURL, _ := url.Parse("https://example.com")
	store := NewStorage(tempDir, "example.com", baseURL)

	t.Run("Save and Get resource", func(t *testing.T) {
		u, _ := url.Parse("https://example.com/page1")
		content := []byte("test content")
		rs := NewResource(u, content, "text/html")

		// Сохраняем ресурс
		if err := store.Save(rs); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		// Проверяем получение
		stored, exists := store.Get(u.String())
		if !exists {
			t.Fatal("Resource not found")
		}

		if string(stored.Content) != string(content) {
			t.Errorf("Content mismatch: got %q, want %q", stored.Content, content)
		}

		// Проверяем что файл создан
		filePath := filepath.Join(tempDir, "page1")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("File not created: %s", filePath)
		}
	})

	t.Run("Save duplicate resource", func(t *testing.T) {
		u, _ := url.Parse("https://example.com/page2")
		rs := NewResource(u, []byte("content"), "text/html")

		// Первое сохранение
		if err := store.Save(rs); err != nil {
			t.Fatalf("First save failed: %v", err)
		}

		// Второе сохранение (дубликат)
		if err := store.Save(rs); err != nil {
			t.Fatalf("Duplicate save should not fail, got: %v", err)
		}
	})

	t.Run("AddLink between resources", func(t *testing.T) {
		parentURL, _ := url.Parse("https://example.com/parent")
		childURL, _ := url.Parse("https://example.com/child")

		parent := NewResource(parentURL, []byte("parent"), "text/html")
		child := NewResource(childURL, []byte("child"), "text/html")

		// Сохраняем оба ресурса
		if err := store.Save(parent); err != nil {
			t.Fatal(err)
		}
		if err := store.Save(child); err != nil {
			t.Fatal(err)
		}

		// Добавляем связь
		store.AddLink(parent, child)

		// Проверяем что связь добавлена
		storedParent, exists := store.Get(parentURL.String())
		if !exists {
			t.Fatal("Parent not found")
		}

		if len(storedParent.Links) != 1 || storedParent.Links[0].URL.String() != childURL.String() {
			t.Errorf("Link not added correctly")
		}
	})

	t.Run("MakeLocalPath for different URLs", func(t *testing.T) {
		testCases := []struct {
			url      string
			expected string
		}{
			{"https://example.com", "index.html"},
			{"https://example.com/", "index.html"},
			{"https://example.com/page", "page"},
			{"https://example.com/path/to/page", filepath.Join("path", "to", "page")},
			{"https://example.com/dir/", filepath.Join("dir", "index.html")},
		}

		for _, tc := range testCases {
			u, _ := url.Parse(tc.url)
			path := MakeLocalPath(u)
			if path != tc.expected {
				t.Errorf("For URL %q expected path %q, got %q", tc.url, tc.expected, path)
			}
		}
	})

	t.Run("Clean storage", func(t *testing.T) {
		if err := store.Clean(); err != nil {
			t.Fatalf("Clean failed: %v", err)
		}

		// Проверяем что директория удалена
		if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
			t.Errorf("Directory was not removed")
		}
	})
}

func TestNewResource(t *testing.T) {
	u, _ := url.Parse("https://example.com/test")
	content := []byte("content")
	contentType := "text/html"

	rs := NewResource(u, content, contentType)

	if rs.URL.String() != u.String() {
		t.Errorf("URL mismatch: got %q, want %q", rs.URL, u)
	}

	if string(rs.Content) != string(content) {
		t.Errorf("Content mismatch: got %q, want %q", rs.Content, content)
	}

	if rs.ContentType != contentType {
		t.Errorf("ContentType mismatch: got %q, want %q", rs.ContentType, contentType)
	}

	if !rs.IsHTML {
		t.Error("IsHTML should be true for text/html")
	}

	// Проверка для не-HTML ресурса
	rs = NewResource(u, content, "text/css")
	if rs.IsHTML {
		t.Error("IsHTML should be false for text/css")
	}
}
