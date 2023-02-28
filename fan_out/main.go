package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

func scrapeWebsite(url string, wg *sync.WaitGroup, results chan<- string) {
	// Decrement the WaitGroup counter when the function completes
	defer wg.Done()

	// Send a GET request to the website
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error scraping %s: %s\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Read the response body into a string
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error scraping %s: %s\n", url, err)
		return
	}
	bodyString := string(bodyBytes)

	// Find and extract the desired data from the HTML
	// (Assuming the data is contained in a <div> tag with class="data")
	dataStart := strings.Index(bodyString, "<div class=\"data\">") + len("<div class=\"data\">")
	dataEnd := strings.Index(bodyString, "</div>")
	data := bodyString[dataStart:dataEnd]

	// Send the data to the results channel
	results <- data
}

func main() {
	urls := []string{"http://example.com", "http://example.org", "http://example.net"}

	// Create a WaitGroup to keep track of the workers
	var wg sync.WaitGroup
	wg.Add(len(urls))

	// Create a results channel to collect the scraped data
	results := make(chan string, len(urls))

	// Start a separate goroutine for each website URL
	for _, url := range urls {
		go scrapeWebsite(url, &wg, results)
	}

	// Wait for all workers to complete
	wg.Wait()

	// Close the results channel to signal that all data has been sent
	close(results)

	// Read the scraped data from the results channel
	for data := range results {
		fmt.Printf("Scraped data: %s\n", data)
	}
}
