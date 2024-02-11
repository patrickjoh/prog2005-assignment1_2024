package handlers

import (
	prog2005_assignment1_2024 "example/assignment1_2024"
	"fmt"
	"net/http"
)

/*
Empty handler as default handler
*/
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	// Ensure request from client is handled as HTML
	w.Header().Set("content-type", "text/html")

	output := fmt.Sprintf(
		`This service does not provide functionality at this path. 
			    Use <a href="%s">%s</a> or <a href="%s">%s</a>.
			    For diagnostic information about the service, visit: <a href="%s">%s</a>`,
		prog2005_assignment1_2024.COUNT_PATH, prog2005_assignment1_2024.COUNT_PATH,
		prog2005_assignment1_2024.READERSHIP_PATH, prog2005_assignment1_2024.READERSHIP_PATH,
		prog2005_assignment1_2024.STATUS_PATH, prog2005_assignment1_2024.STATUS_PATH)

	// Make the output visible to the client
	_, err := fmt.Fprintf(w, "%v", output)

	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}

}
