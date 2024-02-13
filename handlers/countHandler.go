package handlers

import (
	"encoding/json"
	as1 "example/assignment1_2024"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func CountHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handleGetCount(w, r)
	default:
		http.Error(w, "REST method "+r.Method+" is not supported. Only "+http.MethodGet+" is supported.",
			http.StatusNotImplemented)
		return
	}
}

func handleGetCount(w http.ResponseWriter, r *http.Request) {

	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error parsing URL", err)
		return
	}

	lanParam := parsedURL.Query().Get("language")
	if lanParam == "" {
		http.Error(w, "No language parameter", http.StatusBadRequest)
		log.Println("No language parameter")
		return
	}

	// Split the languages into a slice
	languages := strings.Split(lanParam, ",")

	// Initialize a map to hold all ISO codes
	isocode := make(map[string]bool)

	// Loop through each language and add iso code to the map if it is valid
	for _, lang := range languages {
		if len(lang) != 2 {
			// Write which language code is invalid
			http.Error(w, "Invalid language code", http.StatusBadRequest)
			log.Println("Invalid language code")
			return
		}
		isocode[lang] = true
	}

	// Check if any of the codes are less or more than 2 characters

	responseData, err := getBooks(languages)
	if err != nil {
		http.Error(w, "Error during request to GutendexAPI", http.StatusInternalServerError)
		log.Println("Error during request to GutendexAPI")
		return
	}

	// Fraction - Divide the amount of books per language on the total amount of books

	marshallResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error during encoding to JSON", http.StatusInternalServerError)
		log.Println("Error during encoding to JSON")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshallResponse)

}

func getBooks(languages []string) ([]as1.BookCount, error) {

	// Initialize a http client
	client := &http.Client{
		Timeout: time.Second * 30, // Set time for requests to time out
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
	sem := make(chan struct{}, 10)

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

func getTotalBooks(bookUrl string, client *http.Client) (int, error) {

	// Make a request to the GutendexAPI
	bookRes, err := client.Get(bookUrl)
	if err != nil {
		return 0, err
	}

	// Decode the responses into a slice of Book structs
	var bookData as1.Gutendex
	err = json.NewDecoder(bookRes.Body).Decode(&bookData)
	if err != nil {
		return 0, err
	}

	// Close the response body after successfully decoding the data
	defer bookRes.Body.Close()

	return bookData.Count, nil
}
