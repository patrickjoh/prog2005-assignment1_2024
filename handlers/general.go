package handlers

import (
	"encoding/json"
	as1 "example/assignment1_2024"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

func getBooks(languages []string) ([]as1.BookCount, error) {

	// Initialize a http client
	client := &http.Client{
		Timeout: time.Second * 10, // Set time for requests to time out
	}

	// Construct the URL for the GutendexAPI
	bookUrl := as1.GUTENDEX_API + "?languages=" + strings.Join(languages, ",")

	// log.Println("Book URL: ", bookUrl)
	// Make a request to the GutendexAPI
	books, err := fetchBooks(bookUrl, client)
	if err != nil {
		return nil, err
	}

	// Maps to hold track of distinct authors and book counts per language
	bookCount := make(map[string]int)
	uniqueAuthors := make(map[string]map[string]bool)
	authorCount := make(map[string]int)

	// Loop through each book in each language, count number of books and find all unique authors
	for _, book := range books {
		for _, lang := range book.Languages {
			if uniqueAuthors[lang] == nil {
				uniqueAuthors[lang] = make(map[string]bool)
			}
			for _, author := range book.Authors {
				uniqueAuthors[lang][author.Name] = true
			}
			bookCount[lang]++
		}
	}
	// Count the number of unique authors per language
	for lang, authors := range uniqueAuthors {
		authorCount[lang] = len(authors)
	}

	// Make request for all books
	bookUrl = as1.GUTENDEX_API

	// log.Println("Book URL: ", bookUrl)
	// Make a request to the GutendexAPI
	totalBooks, err := getTotalBooks(bookUrl, client)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold all response objects
	var response []as1.BookCount

	// Loop through each language in the response
	for _, lang := range languages {
		fraction := float32(bookCount[lang]) / float32(totalBooks)
		// Create the response object for this language
		responseObj := as1.BookCount{
			Language: lang,
			Books:    bookCount[lang],
			Authors:  authorCount[lang],
			Fraction: fraction,
		}
		// Append the response object to the response slice
		response = append(response, responseObj)

	}

	return response, nil

}

func fetchBooks(bookUrl string, client *http.Client) ([]as1.Book, error) {
	resultsPerPage := 32

	// First, fetch the first page to get the total count
	resp, err := client.Get(bookUrl)
	if err != nil {
		return nil, err
	}

	var initialData as1.Gutendex
	if err := json.NewDecoder(resp.Body).Decode(&initialData); err != nil {
		return nil, err
	}
	resp.Body.Close() // Close the response body after successfully decoding the data

	// Calculate the total number of pages based on the count and results per page
	totalPages := (initialData.Count + resultsPerPage) / resultsPerPage // Ceiling division

	// Create a semaphore to limit the number of concurrent requests
	sem := make(chan struct{}, runtime.NumCPU())

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var allBooks []as1.Book

	fetchPage := func(page int) {
		defer wg.Done()
		sem <- struct{}{}        // Acquire a semaphore
		defer func() { <-sem }() // Release the semaphore

		pageUrl := fmt.Sprintf("%s&page=%d", bookUrl, page)
		resp, err := client.Get(pageUrl)
		if err != nil {
			log.Println("Error fetching page", page, ":", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var pageData as1.Gutendex
			if err := json.NewDecoder(resp.Body).Decode(&pageData); err != nil {
				fmt.Printf("Error decoding page %d: %v\n", page, err)
				return
			}

			mutex.Lock()
			allBooks = append(allBooks, pageData.Books...)
			mutex.Unlock()
		} else {
			fmt.Printf("Unexpected status code for page %d: %d\n", page, resp.StatusCode)
		}
	}

	for page := 1; page <= totalPages; page++ {
		wg.Add(1)
		go fetchPage(page)
	}

	wg.Wait() // Wait for all pages to be fetched

	return allBooks, nil
}
