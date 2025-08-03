package downloader

import (
	"net/url"
	"strings"
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// Валидные URL
		{
			name:    "valid http",
			input:   "http://example.com",
			wantErr: nil,
		},
		{
			name:    "valid https with path",
			input:   "https://example.com/path/to/page",
			wantErr: nil,
		},

		// Невалидные URL
		{
			name:    "empty URL",
			input:   "",
			wantErr: ErrNilURL,
		},
		{
			name:    "no scheme",
			input:   "example.com",
			wantErr: ErrWrongProtocol,
		},
		{
			name:    "invalid scheme",
			input:   "ftp://example.com",
			wantErr: ErrWrongProtocol,
		},
		{
			name:    "URL too long",
			input:   "https://example.com/" + strings.Repeat("a", 2000),
			wantErr: ErrTooLongURL,
		},
		{
			name:    "invalid host chars",
			input:   "https://exa<mple.com",
			wantErr: ErrWrongHost,
		},
		{
			name:    "no host",
			input:   "https://",
			wantErr: ErrNulHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateURL(tt.input)
			if err != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, want %v", err, tt.wantErr)
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
			name:     "localhost",
			input:    "http://localhost:8080",
			expected: "localhost",
		},
		{
			name:     "IP address",
			input:    "https://192.168.1.1",
			expected: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := url.Parse(tt.input)
			result, _ := GetBaseDomain(u)
			if result != tt.expected {
				t.Errorf("getBaseDomain() = %v, want %v", result, tt.expected)
			}
		})
	}
}
