package main

import (
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)

	doc, err := html.Parse(reader)
	if err != nil {
		return []string{}, err
	}
	return getURLsFromNode(doc, rawBaseURL), nil
}

func getURLsFromNode(node *html.Node, rawBaseURL string) []string {
	urls := []string{}

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				newUrl := a.Val
				if strings.HasPrefix(newUrl, "/") {
					newUrl = rawBaseURL + newUrl
				}
				urls = append(urls, newUrl)
				break
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		childUrls := getURLsFromNode(c, rawBaseURL)
		urls = append(urls, childUrls...)
	}
	return urls
}
