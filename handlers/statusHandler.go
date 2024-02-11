package handlers

import (
	"encoding/json"
	as1 "example/assignment1_2024"
	"fmt"
	"net/http"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

// Utility function to perform HTTP GET requests and handle common errors.
func fetchAPIStatus(url string) string {

	client := &http.Client{
		Timeout: time.Second * 10, // Set time for requests to time out
	}

	resp, err := client.Get(url)
	if err != nil {
		return http.StatusText(http.StatusServiceUnavailable)
	}
	defer resp.Body.Close()
	return resp.Status
}

// Handler function for the status endpoint
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetStatus(w)
	default:
		http.Error(w, "REST method "+r.Method+" is not supported. Only "+http.MethodGet+" is supported.", http.StatusNotImplemented)
		return
	}
}

func handleGetStatus(w http.ResponseWriter) {
	status := as1.Status{
		GutendexAPI:  fetchAPIStatus(as1.GUTENDEX_STATUS),
		L2cAPI:       fetchAPIStatus(as1.L2C_STATUS),
		CountriesAPI: fetchAPIStatus(as1.COUNTRY_STATUS),
		Version:      as1.STATUS_VERSION,
		Uptime:       time.Since(startTime),
	}

	jsonBytes, err := json.Marshal(status)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in outputting JSON: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
