package handlers

import (
	"encoding/json"
	as1 "example/assignment1_2024"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ReadershipHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handleGetReadership(w, r)
	default:
		http.Error(w, "REST method "+r.Method+" is not supported. Only "+http.MethodGet+" is supported.",
			http.StatusNotImplemented)
		return
	}
}

/*
handleGetReadership responds with amount of books and author,
and the population of all countries where that language is spoken, per country that speaks the language
*/
func handleGetReadership(w http.ResponseWriter, r *http.Request) {

	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error parsing URL", err)
		return
	}

	limitStr := parsedURL.Query().Get("limit")
	limitLen, err := strconv.Atoi(limitStr)
	if limitStr != "" && err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		log.Println("Invalid limit parameter")
		return
	}

	var limitBool bool
	// Checks if the limit parameter is set or not
	if limitStr == "" {
		limitBool = false
	} else {
		limitBool = true
	}

	// Get the language from the URL
	path := parsedURL.Path

	// Split the URL into parts delimited by "/"
	parts := strings.Split(path, "/")

	var langCode []string
	if len(parts) == 5 {
		langCode = append(langCode, parts[4]) // Append the language code to the slice
	} else {
		http.Error(w, "URL does not contain all necessary parts", http.StatusBadRequest)
		log.Println("URL in request does not contain all necessary parts")
		return
	}

	// Check if the language code is valid
	if len(langCode[0]) != 2 {
		http.Error(w, "Invalid language code. Must be a 2-letter code", http.StatusBadRequest)
		log.Println("Invalid language code. Must be a 2-letter code")
		return
	}
	// Get the books and authors from "Gutendex" API
	books, err := getBooks(langCode)
	if err != nil {
		http.Error(w, "Error during request to GutendexAPI", http.StatusInternalServerError)
		log.Println("Error during request to GutendexAPI")
		return
	}

	// Get the countries that speaks the specified language, from L2C API
	l2cData, err := getLanguage(langCode[0])
	if err != nil {
		http.Error(w, "Error during request to L2CAPI", http.StatusInternalServerError)
		log.Println("Error during request to L2CAPI")
		return
	}

	// Initialize a map to hold the name of each country with the isocode as the key
	countryMap := make(map[string]string)

	// Loop through all countries and add the name to the map
	for _, country := range l2cData {
		countryMap[country.Alpha2Code] = country.Name
	}

	// Loop through all countries and add the isocode to a slice
	var isocode []string
	for _, country := range l2cData {
		isocode = append(isocode, country.Alpha2Code)
	}

	// Get the population of the countries that speaks the specified language, from RESTcountries API
	countryData, err := getCountries(isocode)

	// Initialize a response object
	response := []as1.Readership{}
	limiter := 0

	// Loop through each country in the response
	for _, country := range countryData {
		// Limits the number of countries returned
		if limitBool && limitLen == limiter {
			break
		}
		// Create the response object for this country
		responseObj := as1.Readership{
			Country:    countryMap[country.Alpha2Code],
			Isocode:    country.Alpha2Code,
			Books:      books[0].Books,
			Authors:    books[0].Authors,
			Readership: country.Population,
		}
		// Append the response object to the response slice
		response = append(response, responseObj)
		limiter++
	}

	// Marshal the response slice to JSON
	marshallResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error during encoding to JSON", http.StatusInternalServerError)
		log.Println("Error during encoding to JSON")
		return
	}
	// Ensure that the server interprets requests as JSON from Client (browser)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshallResponse)
}

func getLanguage(langCode string) ([]as1.L2C, error) {
	// Get the language from the "L2C" API
	langUrl := as1.L2C_API + langCode
	resp, err := http.Get(langUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var langData []as1.L2C
	err = json.NewDecoder(resp.Body).Decode(&langData)
	if err != nil {
		return nil, err
	}

	return langData, nil
}

func getCountries(countryCodes []string) ([]as1.Country, error) {
	// Get the countries from the "RESTcountries" API
	countryUrl := as1.COUNTRY_API + strings.Join(countryCodes, ",")
	resp, err := http.Get(countryUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var countryData []as1.Country
	err = json.NewDecoder(resp.Body).Decode(&countryData)
	if err != nil {
		return nil, err
	}

	return countryData, nil
}
