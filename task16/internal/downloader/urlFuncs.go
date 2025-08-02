package downloader

import (
	"errors"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func ValidateUrl(rawUrl string) (*url.URL, error) {
	if len(rawUrl) > 2000 {
		return nil, errors.New("URL too long")
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, errors.New("invalid URL format")
	}
	if strings.ContainsAny(u.Host, " <>\"'{}|\\^`") {
		return nil, errors.New("invalid characters in host")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("only http and https protocols are allowed")
	}

	if u.Host == "" {
		return nil, errors.New("host is required")
	}
	return u, nil
}

func GetBaseDomain(u *url.URL) (string, error) {
	domain, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname())
	if err != nil {
		return u.Hostname(), nil // fallback
	}
	return domain, nil
}
