package main

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseUrl            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	fmt.Println("Hello, World!")

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", args[0])
	baseUrl, err := url.Parse(args[0])
	if err != nil {
		fmt.Printf("Error parsing url: %v", err)
		os.Exit(1)
	}

	maxThreads := 10
	if len(args) > 1 {
		maxThreads, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error parsing max threads: %v", err)
			os.Exit(1)
		}
	}

	maxPages := 500
	if len(args) > 2 {
		maxPages, err = strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("Error parsing max pages: %v", err)
			os.Exit(1)
		}
	}

	conf := config{
		pages:              make(map[string]int),
		baseUrl:            baseUrl,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxThreads),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	conf.crawlPage(baseUrl.String())

	conf.wg.Wait()

	conf.printReport()
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.wg.Add(1)

	go func() {
		defer cfg.wg.Done()
		cfg.processPage(rawCurrentURL)
	}()
}

func (cfg *config) processPage(rawCurrentURL string) {
	fmt.Printf("Crawling: %s\n", rawCurrentURL)

	baseUrl, err := url.Parse(cfg.baseUrl.String())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	currentUrl, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	baseHostName := baseUrl.Hostname()
	currentHostName := currentUrl.Hostname()

	if baseHostName != currentHostName {
		fmt.Printf("Outside: %v, %v\n", baseHostName, currentHostName)
		return
	}

	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	cfg.mu.Lock()

	if cfg.maxPages <= len(cfg.pages) {
		cfg.mu.Unlock()
		return
	}

	if _, exists := cfg.pages[normalizedURL]; exists {
		cfg.pages[normalizedURL] += 1
		cfg.mu.Unlock()
		return
	}
	cfg.pages[normalizedURL] = 1
	cfg.mu.Unlock()

	htmlContent, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	urls, err := getURLsFromHTML(htmlContent, cfg.baseUrl.String())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found urls %d\n", len(urls))

	for _, url := range urls {
		cfg.crawlPage(url)
	}
}

func (cfg *config) printReport() {
	fmt.Printf(`
=============================
  REPORT for %s
=============================
`, cfg.baseUrl)

	type PageEntry struct {
		Page  string
		Count int
	}

	entriesList := []PageEntry{}
	for p, c := range cfg.pages {
		entry := PageEntry{
			Page:  p,
			Count: c,
		}

		entriesList = append(entriesList, entry)
	}

	sort.Slice(entriesList, func(i, j int) bool {

		if entriesList[i].Count > entriesList[j].Count {
			return true
		}
		if entriesList[i].Count < entriesList[j].Count {
			return false
		}
		if entriesList[i].Page > entriesList[j].Page {
			return true
		}
		if entriesList[i].Page < entriesList[j].Page {
			return false
		}
		return false
	})

	for _, entry := range entriesList {
		fmt.Printf("Found %d internal links to %s\n", entry.Count, entry.Page)
	}
}
