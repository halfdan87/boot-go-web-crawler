package main

import (
	"net/url"
	"strings"
)

func normalizeURL(urlStr string) (string, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	host := url.Hostname()
	return host + strings.TrimSuffix(url.Path, "/"), nil
}
