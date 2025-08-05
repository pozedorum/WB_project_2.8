package parser

import (
	"bytes"
	"errors"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Мапа тегов и их атрибутов, содержащих ссылки
var linkAttributes = map[string][]string{
	"a":      {"href"},
	"link":   {"href"},
	"area":   {"href"},
	"img":    {"src", "data-src", "srcset"},
	"script": {"src"},
	"iframe": {"src"},
	"embed":  {"src"},
	"source": {"src", "srcset"},
	"form":   {"action"},
	"object": {"data"},
	"video":  {"src", "poster"},
	"audio":  {"src"},
	"track":  {"src"},
}

var (
	ErrUnsupportedScheme = errors.New("unsupported URL scheme")
	ErrExternalLink      = errors.New("external link")
)

type Parser struct {
	BaseURL *url.URL
}

func NewParser(baseURL *url.URL) *Parser {
	return &Parser{BaseURL: baseURL}
}

func (p *Parser) ExtractLinks(htmlContent []byte) ([]*url.URL, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var links []*url.URL
	visited := make(map[string]struct{})

	for node := range doc.Descendants() {
		if node.Type != html.ElementNode {
			continue
		}
		// Получаем атрибуты для тега
		attrs, ok := linkAttributes[node.Data]
		if !ok {
			continue
		}

		for _, attr := range attrs {
			attrValue := getAttr(node, attr)
			if attrValue == "" {
				continue
			}

			// Обработка srcset (специальный формат)
			if attr == "srcset" {
				urls := parseSrcSet(attrValue)

				for _, rawU := range urls {
					u, err := p.NormalizeURL(rawU)
					if err != nil {
						continue
					}
					if _, exists := visited[u.String()]; !exists {
						visited[u.String()] = struct{}{}
						links = append(links, u)
					}
				}

				continue
			}

			if u, err := p.NormalizeURL(attrValue); err == nil {
				if _, exists := visited[u.String()]; !exists {
					visited[u.String()] = struct{}{}
					links = append(links, u)
				}
			}
		}
	}

	return links, nil
}

func (p *Parser) NormalizeURL(rawURL string) (*url.URL, error) {
	if rawURL == "" || strings.HasPrefix(rawURL, "javascript:") ||
		strings.HasPrefix(rawURL, "mailto:") || strings.HasPrefix(rawURL, "tel:") {
		return nil, ErrUnsupportedScheme
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if u.Hostname() == "" {
		u = p.BaseURL.ResolveReference(u)
	} else if !isSameDomain(u.Hostname(), p.BaseURL.Hostname()) {
		return nil, ErrExternalLink
	}
	if !u.IsAbs() {
		u = p.BaseURL.ResolveReference(u)
	}

	// Очищаем URL от лишней информации
	u.Fragment = "" // (#section)
	u.RawQuery = "" // query-параметры (?foo=bar)

	return u, nil
}

func parseSrcSet(srcset string) []string {
	var urls []string
	for _, part := range strings.Split(srcset, ",") {
		if url := strings.Split(strings.TrimSpace(part), " ")[0]; url != "" {
			urls = append(urls, url)
		}
	}
	return urls
}

func getAttr(n *html.Node, attrName string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func isSameDomain(host, baseHost string) bool {
	if host == baseHost {
		return true
	}
	// Для поддоменов: a.b.example.com и example.com
	return strings.HasSuffix(host, "."+baseHost)
}
