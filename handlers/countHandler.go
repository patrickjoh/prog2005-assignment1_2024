package handlers

import (
	"encoding/json"
	as1 "example/assignment1_2024"
	"log"
	"net/http"
	"net/url"
	"strings"
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
