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

// Handler function for the status endpoint
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handleGetStatus(w)
	default:
		http.Error(w, "REST method "+r.Method+" is not supported. Only "+http.MethodGet+" is supported.",
			http.StatusNotImplemented)
		return
	}
}

func handleGetStatus(w http.ResponseWriter) {

	status := as1.Status{}

	gutendexRes, err := http.Get(as1.GUTENDEX_STATUS)
	if err != nil {
		status.GutendexAPI = http.StatusText(http.StatusServiceUnavailable)
	} else {
		status.GutendexAPI = gutendexRes.Status
		defer gutendexRes.Body.Close()
	}

	l2cRes, err := http.Get(as1.L2C_STATUS)
	if err != nil {
		status.L2cAPI = http.StatusText(http.StatusServiceUnavailable)
	} else {
		status.L2cAPI = l2cRes.Status
		defer l2cRes.Body.Close()
	}

	countryRes, err := http.Get(as1.COUNTRY_STATUS)
	if err != nil {
		status.CountriesAPI = http.StatusText(http.StatusServiceUnavailable)
	} else {
		status.CountriesAPI = countryRes.Status
		defer countryRes.Body.Close()
	}

	status.Version = as1.STATUS_VERSION
	status.Uptime = time.Since(startTime)

	// Encode struct to JSON
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in outputting JSON: %s", err), http.StatusInternalServerError)
		return
	}

	// Write JSON to ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
