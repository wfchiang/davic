package main

import (
	"fmt"
	"log"
	"net/http"
)

// ====
// Handlers 
// ====
func homepageHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	fmt.Println("Homepage is hit")
	fmt.Fprintf(http_resp, "Rest-mocker here...")
}

// ====
// Main
// ====
func main () {
	fmt.Println("Starting rest-mocker...")
	http.HandleFunc("/", homepageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}