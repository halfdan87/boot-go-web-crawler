package main

import (
	"fmt"
	"io"
	"net/http"
)

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Error getting %s : %s", rawURL, resp.Status)
	}
	if err != nil {
		return "", err
	}

	return string(body), nil
}
