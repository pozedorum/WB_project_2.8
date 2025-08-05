package storage

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStorage(t *testing.T) {
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

		if err := store.Save(rs); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		stored, exists := store.Get(u.String())
		if !exists {
			t.Fatal("Resource not found")
		}

		if string(stored.Content) != string(content) {
			t.Errorf("Content mismatch: got %q, want %q", stored.Content, content)
		}

		expectedPath := filepath.Join("example.com", "page1.html")
		if stored.LocalPath != expectedPath {
			t.Errorf("LocalPath mismatch: got %q, want %q", stored.LocalPath, expectedPath)
		}
	})

	t.Run("MakeLocalPath for different URLs", func(t *testing.T) {
		testCases := []struct {
			url      string
			expected string
		}{
			{"https://example.com", filepath.Join("example.com", "index.html")},
			{"https://example.com/", filepath.Join("example.com", "index.html")},
			{"https://example.com/page", filepath.Join("example.com", "page.html")},
			{"https://example.com/path/to/page", filepath.Join("example.com", "path", "to", "page.html")},
			{"https://example.com/dir/", filepath.Join("example.com", "dir", "index.html")},
			{"https://example.com/script.js", filepath.Join("example.com", "script.js")},
			{"https://example.com/long" + strings.Repeat("a", 200), filepath.Join("example.com", "long"+strings.Repeat("a", 100-4)+".html")},
		}

		for _, tc := range testCases {
			u, _ := url.Parse(tc.url)
			path := MakeLocalPath(u)
			if path != tc.expected {
				t.Errorf("For URL %q expected path %q, got %q", tc.url, tc.expected, path)
			}
		}
	})

	t.Run("Save resources with same path but different hosts", func(t *testing.T) {
		u1, _ := url.Parse("https://sub1.example.com/page")
		u2, _ := url.Parse("https://sub2.example.com/page")
		content := []byte("content")

		rs1 := NewResource(u1, content, "text/html")
		rs2 := NewResource(u2, content, "text/html")

		if err := store.Save(rs1); err != nil {
			t.Fatal(err)
		}
		if err := store.Save(rs2); err != nil {
			t.Fatal(err)
		}

		if _, exists := store.Get(u1.String()); !exists {
			t.Error("First resource not found")
		}
		if _, exists := store.Get(u2.String()); !exists {
			t.Error("Second resource not found")
		}

		// Проверяем что пути разные
		if rs1.LocalPath == rs2.LocalPath {
			t.Error("Resources with different hosts should have different paths")
		}
	})
}
