package main

import (
	"asciiartweb"
	"fmt"
	"net/http"
)

/**
* Route registering and server initialization
 */
func main() {
	// initilizing mux for matching the incoming requests to the correct to the registered routes and to the handler functions
	mux := http.NewServeMux()

	mux.HandleFunc("/", asciiartweb.HomeHandler)
	mux.HandleFunc("/ascii-art", asciiartweb.ASCIIArtHandler)
	mux.HandleFunc("/export-png", asciiartweb.ExportPNGHandler)
	mux.HandleFunc("/export-txt", asciiartweb.ExportTXTHandler)
	mux.HandleFunc("/export-html", asciiartweb.ExportHTMLHandler)

	mux.Handle(
		"/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))),
	)

	fmt.Println("Server is running on http://localhost:8081")

	err := http.ListenAndServe(":8081", mux)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
