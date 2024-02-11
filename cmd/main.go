package main

import (
	as1 "example/assignment1_2024"
	"example/assignment1_2024/handlers"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: " + as1.DEFAULT_PORT)
		port = as1.DEFAULT_PORT
	}

	// Handler endpoints
	http.HandleFunc(as1.DEFAULT_PATH, handlers.DefaultHandler)
	http.HandleFunc(as1.COUNT_PATH, handlers.CountHandler)
	http.HandleFunc(as1.READERSHIP_PATH, handlers.ReadershipHandler)
	http.HandleFunc(as1.STATUS_PATH, handlers.StatusHandler)

	log.Println("Starting server on " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
