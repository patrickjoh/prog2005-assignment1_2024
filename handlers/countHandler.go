package handlers

import (
	"encoding/json"
	as1 "example/assignment1_2024"
	"log"
	"net/http"
	"net/url"
	"strings"
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

	// Initialize a slice to hold all ISO codes
	var isocode []string

	// Loop through each university in the response
	for _, lang := range languages {
		// Check if the iso code is already in the slice
		found := false
		// Loop through each iso code in the slice
		for _, code := range isocode {
			// If the iso code is found, break out of the loop
			if lang == code {
				found = true
				break
			}
		}
		// If the iso code is not found, append it to the slice
		if !found {
			isocode = append(isocode, lang)
		}
	}

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

	// Construct the URL for the GutendexAPI
	bookUrl := as1.GUTENDEX_API + "?languages=" + strings.Join(languages, ",")

	log.Println("Book URL: ", bookUrl)
	// Make a request to the GutendexAPI
	books, err := fetchBooks(bookUrl)
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

	log.Println("Book URL: ", bookUrl)
	// Make a request to the GutendexAPI
	totalBooks, err := getTotalBooks(bookUrl)
	if err != nil {
		return nil, err
	}

	log.Println("Total books: ", totalBooks)
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

func fetchBooks(bookUrl string) ([]as1.Book, error) {
	var allBooks []as1.Book

	// Initialize a http client
	client := &http.Client{
		Timeout: time.Second * 10, // Set time for requests to time out
	}

	for bookUrl != "" {
		// Make a request to the GutendexAPI
		bookRes, err := client.Get(bookUrl)
		if err != nil {
			return nil, err
		}

		// Decode the responses into a slice of Book structs
		var bookData as1.Gutendex
		err = json.NewDecoder(bookRes.Body).Decode(&bookData)
		if err != nil {
			return nil, err
		}

		// Close the response body after successfully decoding the data
		bookRes.Body.Close()

		// Append the book data to the allBooks slice
		allBooks = append(allBooks, bookData.Books...)
		// Set the next URL to the next page of data
		bookUrl = bookData.Next
	}

	return allBooks, nil
}

func getTotalBooks(bookUrl string) (int, error) {
	// Initialize a http client
	client := &http.Client{
		Timeout: time.Second * 10, // Set time for requests to time out
	}

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
