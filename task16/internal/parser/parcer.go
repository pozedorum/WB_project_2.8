package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
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

	var (
		links   []*url.URL
		visited = make(map[string]struct{})
	)

	for node := range doc.Descendants() {
		if node.Type != html.ElementNode {
			continue
		}
		// Тут как я понимаю может быть несколько ссылок в одной ноде и надо как-то это обработать
		var attr string
		switch node.Data {
		case "a", "link", "area":
			attr = "href"
		case "img", "script", "iframe", "embed", "source":
			attr = "src"
		case "form":
			attr = "action"
		default:
			continue
		}

		for _, nAttr := range node.Attr {
			if nAttr.Key != attr {
				continue
			}
			u, err := p.NormalizeURL(nAttr.Val)
			if err != nil {
				continue
			}
			if _, ok := visited[u.String()]; !ok {
				visited[u.String()] = struct{}{}
				links = append(links, u)
			}
		}
	}

	return links, nil
}

func (p *Parser) NormalizeURL(rawURL string) (*url.URL, error) {
	if rawURL == "" || strings.HasPrefix(rawURL, "javascript:") ||
		strings.HasPrefix(rawURL, "mailto:") || strings.HasPrefix(rawURL, "tel:") {
		return nil, fmt.Errorf("unsupported URL scheme")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(u.Hostname(), p.BaseURL.Hostname()) {
		return nil, fmt.Errorf("external link") // Пропускаем внешние ссылки
	}
	if !u.IsAbs() {
		u = p.BaseURL.ResolveReference(u)
	}

	// Очищаем URL от лишней информации
	u.Fragment = "" // (#section)
	u.RawQuery = "" // query-параметры (?foo=bar)

	return u, nil
}
