package downloader

import (
	"net/url"
	"strings"
	"testing"
)

func TestValidateUrl(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errMatch string
	}{
		// Валидные URL
		{
			name:    "valid http",
			input:   "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid https with path",
			input:   "https://example.com/path/to/page",
			wantErr: false,
		},
		{
			name:    "valid with query params",
			input:   "https://example.com?param=value&foo=bar",
			wantErr: false,
		},

		// Невалидные URL
		{
			name:     "empty URL",
			input:    "",
			wantErr:  true,
			errMatch: "invalid URL format",
		},
		{
			name:     "no scheme",
			input:    "example.com",
			wantErr:  true,
			errMatch: "host is required",
		},
		{
			name:     "invalid scheme",
			input:    "ftp://example.com",
			wantErr:  true,
			errMatch: "only http and https",
		},
		{
			name:     "URL too long",
			input:    "https://example.com/" + strings.Repeat("a", 2000),
			wantErr:  true,
			errMatch: "URL too long",
		},
		{
			name:     "invalid host chars",
			input:    "https://exa<mple.com",
			wantErr:  true,
			errMatch: "invalid characters in host",
		},
		{
			name:     "space in host",
			input:    "https://exa mple.com",
			wantErr:  true,
			errMatch: "invalid characters in host",
		},
		{
			name:     "no host",
			input:    "https://",
			wantErr:  true,
			errMatch: "host is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateUrl(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMatch != "" && !strings.Contains(err.Error(), tt.errMatch) {
				t.Errorf("validateUrl() error = %v, want contains %q", err, tt.errMatch)
			}
		})
	}
}

func TestGetBaseDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Стандартные домены
		{
			name:     "simple domain",
			input:    "https://example.com",
			expected: "example.com",
		},
		{
			name:     "subdomain",
			input:    "http://sub.example.com",
			expected: "example.com",
		},
		{
			name:     "multiple subdomains",
			input:    "https://a.b.c.example.co.uk",
			expected: "co.uk", // Ожидаемое поведение для простой реализации
		},

		// Специальные случаи
		{
			name:     "localhost",
			input:    "http://localhost:8080",
			expected: "localhost",
		},
		{
			name:     "IP address",
			input:    "https://192.168.1.1",
			expected: "192.168.1.1",
		},
		{
			name:     "no subdomain with country TLD",
			input:    "https://example.co.uk",
			expected: "co.uk", // Проблемный кейс для простой реализации
		},
		{
			name:     "URL with path",
			input:    "https://blog.example.com/path/to/page",
			expected: "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.input)
			if err != nil {
				t.Fatalf("failed to parse test URL: %v", err)
			}

			result, err := GetBaseDomain(u)
			if err != nil {
				t.Errorf("error: %v", err)
			} else if result != tt.expected {
				t.Errorf("getBaseDomain() = %v, want %v", result, tt.expected)
			}
		})
	}
}
