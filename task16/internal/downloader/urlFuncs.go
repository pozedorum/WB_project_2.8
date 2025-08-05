package downloader

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var (
	ErrNilURL           = errors.New("URL is nil")
	ErrTooLongURL       = errors.New("URL too long")
	ErrInvalidFormatURL = errors.New("invalid URL format")
	ErrWrongHost        = errors.New("invalid characters in host")
	ErrWrongProtocol    = errors.New("only http and https protocols are allowed")
	ErrNulHost          = errors.New("host is required")
	ErrNotFound         = errors.New("resource is not found")
)

func ValidateURL(rawURL string) (*url.URL, error) {
	if rawURL == "" {
		return nil, ErrNilURL
	}
	if len(rawURL) > 2000 {
		return nil, ErrTooLongURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, ErrInvalidFormatURL
	}
	if strings.ContainsAny(u.Host, " <>\"'{}|\\^`") {
		return nil, ErrWrongHost
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, ErrWrongProtocol
	}

	if u.Host == "" {
		return nil, ErrNulHost
	}
	return u, nil
}

func GetBaseDomain(u *url.URL) (string, error) {
	hostname := u.Hostname()

	// Специальная обработка для IP-адресов
	if ip := net.ParseIP(hostname); ip != nil {
		return hostname, nil
	}

	// Специальная обработка для localhost
	if hostname == "localhost" {
		return hostname, nil
	}

	// Обычная обработка доменов
	domain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return "", fmt.Errorf("failed to get base domain: %w", err)
	}
	return domain, nil
}
