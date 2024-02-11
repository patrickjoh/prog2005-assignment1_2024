package handlers

import "net/http"

func ReadershipHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handleGetReadership(w)
	default:
		http.Error(w, "REST method "+r.Method+" is not supported. Only "+http.MethodGet+" is supported.",
			http.StatusNotImplemented)
		return
	}
}

func handleGetReadership(w http.ResponseWriter) {

}
