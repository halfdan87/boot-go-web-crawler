package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	fmt.Println("Hello, World!")

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseUrl := args[0]
	fmt.Printf("starting crawl of: %s\n", baseUrl)

	pageHtml, err := getHTML(baseUrl)
	if err != nil {
		fmt.Println("Error getting html: %v", err)
		os.Exit(1)
	}

	fmt.Println(pageHtml)
}

func printReport(pages map[string]int, baseUrl string) {
	fmt.Print(`
=============================
  REPORT for https://example.com
=============================
` + baseUrl)

	type PageEntry struct {
		Page  string
		Count int
	}

	entriesList := []PageEntry{}
	for p, c := range pages {
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
		// Now we sort by text
		if entriesList[i].Page > entriesList[j].Page {
			return true
		}
		if entriesList[i].Page < entriesList[j].Page {
			return false
		}
		return false
	})

}
