package parser

import (
	"net/url"
	"testing"
)

func TestParser_ExtractLinks(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com")
	p := NewParser(baseURL)

	tests := []struct {
		name     string
		html     string
		expected []string
	}{
		{
			name: "simple anchor",
			html: `<a href="/about">About</a>`,
			expected: []string{
				"https://example.com/about",
			},
		},
		{
			name: "multiple elements",
			html: `
				<img src="/img/logo.png">
				<script src="/js/app.js"></script>
				<link href="/css/style.css" rel="stylesheet">
			`,
			expected: []string{
				"https://example.com/img/logo.png",
				"https://example.com/js/app.js",
				"https://example.com/css/style.css",
			},
		},
		{
			name: "srcset parsing",
			html: `<img srcset="/img-1x.jpg 1x, /img-2x.jpg 2x">`,
			expected: []string{
				"https://example.com/img-1x.jpg",
				"https://example.com/img-2x.jpg",
			},
		},
		{
			name: "ignore external links",
			html: `
				<a href="https://other.com/external"></a>
				<img src="https://example.com/local.jpg">
			`,
			expected: []string{
				"https://example.com/local.jpg",
			},
		},
		{
			name: "subdomains",
			html: `<img src="https://cdn.example.com/image.jpg">`,
			expected: []string{
				"https://cdn.example.com/image.jpg",
			},
		},
		{
			name: "ignore invalid schemes",
			html: `
				<a href="mailto:test@example.com"></a>
				<script src="javascript:alert()"></script>
			`,
			expected: []string{},
		},
		{
			name: "fragment and query params",
			html: `<a href="/path?foo=bar#section"></a>`,
			expected: []string{
				"https://example.com/path",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links, err := p.ExtractLinks([]byte(tt.html))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(links) != len(tt.expected) {
				t.Fatalf("expected %d links, got %d", len(tt.expected), len(links))
			}

			for i, expected := range tt.expected {
				if actual := links[i].String(); actual != expected {
					t.Errorf("link %d: expected %q, got %q", i, expected, actual)
				}
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com")
	p := NewParser(baseURL)

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "relative path",
			input:    "/about",
			expected: "https://example.com/about",
		},
		{
			name:     "absolute url",
			input:    "https://example.com/path",
			expected: "https://example.com/path",
		},
		{
			name:     "subdomain",
			input:    "https://cdn.example.com/img.jpg",
			expected: "https://cdn.example.com/img.jpg",
		},
		{
			name:    "external domain",
			input:   "https://other.com/path",
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			input:   "javascript:alert()",
			wantErr: true,
		},
		{
			name:     "clean fragments and queries",
			input:    "/path?foo=bar#section",
			expected: "https://example.com/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := p.NormalizeURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got := u.String(); got != tt.expected {
				t.Errorf("NormalizeURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseSrcSet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple",
			input:    "image-1x.jpg 1x, image-2x.jpg 2x",
			expected: []string{"image-1x.jpg", "image-2x.jpg"},
		},
		{
			name:     "with spaces",
			input:    "  image-1x.jpg  1x,  image-2x.jpg  2x  ",
			expected: []string{"image-1x.jpg", "image-2x.jpg"},
		},
		{
			name:     "with parameters",
			input:    "image-1x.jpg?size=small 1x, image-2x.jpg?size=large 2x",
			expected: []string{"image-1x.jpg?size=small", "image-2x.jpg?size=large"},
		},
		{
			name:     "empty",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSrcSet(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d urls, got %d", len(tt.expected), len(got))
			}
			for i, expected := range tt.expected {
				if got[i] != expected {
					t.Errorf("url %d: expected %q, got %q", i, expected, got[i])
				}
			}
		})
	}
}
