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
			"simple anchor",
			`<a href="/about">About</a>`,
			[]string{"https://example.com/about"},
		},
		{
			"multiple elements",
			`<img src="/img.png"><script src="/js.js"></script>`,
			[]string{"https://example.com/img.png", "https://example.com/js.js"},
		},
		{
			"srcset parsing",
			`<img srcset="/img-1x.jpg 1x, /img-2x.jpg 2x">`,
			[]string{"https://example.com/img-1x.jpg", "https://example.com/img-2x.jpg"},
		},
		{
			"ignore external",
			`<a href="https://other.com"></a><img src="/local.jpg">`,
			[]string{"https://example.com/local.jpg"},
		},
		{
			"subdomains",
			`<img src="https://cdn.example.com/img.jpg">`,
			[]string{"https://cdn.example.com/img.jpg"},
		},
		{
			"invalid schemes",
			`<a href="javascript:alert()"></a>`,
			[]string{},
		},
		{
			"fragment and query",
			`<a href="/path?foo=bar#section"></a>`,
			[]string{"https://example.com/path"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links, err := p.ExtractLinks([]byte(tt.html))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(links) != len(tt.expected) {
				t.Fatalf("Expected %d links, got %d", len(tt.expected), len(links))
			}

			for i, expected := range tt.expected {
				if links[i].String() != expected {
					t.Errorf("Link %d: expected %q, got %q", i, expected, links[i].String())
				}
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com")
	p := NewParser(baseURL)

	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"/about", "https://example.com/about", false},
		{"https://example.com/path", "https://example.com/path", false},
		{"//cdn.example.com/img.jpg", "https://cdn.example.com/img.jpg", false},
		{"https://other.com/path", "", true},
		{"javascript:alert()", "", true},
		{"/path?foo=bar#section", "https://example.com/path", false},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			u, err := p.NormalizeURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizeURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if u.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, u.String())
			}
		})
	}
}

func TestParseSrcSet(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"image-1x.jpg 1x, image-2x.jpg 2x", []string{"image-1x.jpg", "image-2x.jpg"}},
		{"  image-1x.jpg  1x,  image-2x.jpg  2x  ", []string{"image-1x.jpg", "image-2x.jpg"}},
		{"image.jpg?size=small 1x", []string{"image.jpg?size=small"}},
		{"", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseSrcSet(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d urls, got %d", len(tt.expected), len(result))
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("URL %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestIsSameDomain(t *testing.T) {
	tests := []struct {
		host     string
		baseHost string
		expected bool
	}{
		{"example.com", "example.com", true},
		{"sub.example.com", "example.com", true},
		{"example.com", "other.com", false},
		{"evil.com", "example.com", false},
		{"localhost", "localhost", true},
		{"192.168.1.1", "192.168.1.1", true},
	}

	for _, tt := range tests {
		t.Run(tt.host+"_"+tt.baseHost, func(t *testing.T) {
			result := isSameDomain(tt.host, tt.baseHost)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
