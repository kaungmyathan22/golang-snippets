package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Starting server on :4000")
	errCh := make(chan error, 1)

	// Start the server in a goroutine and check for errors
	go func() {
		err := http.ListenAndServe(":4000", mux)
		if err != nil {
			errCh <- err
		}
	}()

	// Wait for the error channel to receive an error or close
	err := <-errCh
	log.Fatal(err) // Log and exit gracefully
}
