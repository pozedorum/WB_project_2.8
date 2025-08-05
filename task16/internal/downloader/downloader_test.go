package downloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"task16/internal/storage"
)

func TestDownloader_Download(t *testing.T) {
	// Создаем тестовый HTTP сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.Header().Set("Content-Type", "text/html")
			if _, err := w.Write([]byte("<html>test</html>")); err != nil {
				t.Fatal(err)
			}
		case "/error":
			w.WriteHeader(http.StatusNotFound)
		case "/redirect":
			w.Header().Set("Location", "/success")
			w.WriteHeader(http.StatusMovedPermanently)
		case "/slow":
			time.Sleep(2 * time.Second)
			if _, err := w.Write([]byte("ok")); err != nil {
				t.Fatal(err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	// Создаем временное хранилище
	tempDir := t.TempDir()
	store := storage.NewStorage(tempDir, "test", &url.URL{Host: "test"})
	dl := NewDownloader(store)

	t.Run("successful download", func(t *testing.T) {
		u, _ := url.Parse(ts.URL + "/success")
		rs, fromCache, err := dl.Download(context.Background(), u)

		if err != nil {
			t.Fatalf("Download failed: %v", err)
		}
		if fromCache {
			t.Error("Resource should not be from cache")
		}
		if !strings.Contains(string(rs.Content), "<html>test</html>") {
			t.Errorf("Unexpected content: %s", rs.Content)
		}
		if rs.ContentType != "text/html" {
			t.Errorf("Unexpected content type: %s", rs.ContentType)
		}
	})

	t.Run("error response", func(t *testing.T) {
		u, _ := url.Parse(ts.URL + "/error")
		_, _, err := dl.Download(context.Background(), u)

		if err == nil {
			t.Error("Expected error for 404 response")
		}
	})

	t.Run("cached resource", func(t *testing.T) {
		u, _ := url.Parse(ts.URL + "/success")
		// Первый запрос - сохраняем в кеш
		_, _, err := dl.Download(context.Background(), u)
		if err != nil {
			t.Fatal(err)
		}

		// Второй запрос - должен быть из кеша
		_, fromCache, err := dl.Download(context.Background(), u)
		if err != nil {
			t.Fatal(err)
		}
		if !fromCache {
			t.Error("Resource should be from cache")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		u, _ := url.Parse(ts.URL + "/slow")
		_, _, err := dl.Download(ctx, u)

		if err == nil {
			t.Error("Expected timeout error")
		}
	})

}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid http", "http://example.com", false},
		{"valid https", "https://example.com/path", false},
		{"empty", "", true},
		{"no scheme", "example.com", true},
		{"invalid scheme", "ftp://example.com", true},
		{"invalid chars", "https://exa<mple.com", true},
		{"too long", "https://example.com/" + strings.Repeat("a", 2000), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBaseDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com", "example.com"},
		{"https://sub.example.com", "example.com"},
		{"http://localhost:8080", "localhost"},
		{"http://192.168.1.1", "192.168.1.1"},
		{"http://user:pass@example.com", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			u, _ := url.Parse(tt.input)
			result, err := GetBaseDomain(u)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
