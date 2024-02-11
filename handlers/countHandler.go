package handlers

import "net/http"

func CountHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handleGetCount(w)
	default:
		http.Error(w, "REST method "+r.Method+" is not supported. Only "+http.MethodGet+" is supported.",
			http.StatusNotImplemented)
		return
	}
}

func handleGetCount(w http.ResponseWriter) {

}
