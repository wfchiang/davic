package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
)

// ====
// Handlers 
// ====
func homepageHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	log.Println("Homepage is Hit!")
	fmt.Fprintf(http_resp, "Rest-mocker is Here!")
}

func echoHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body)
	if err != nil {
		log.Println("[Echo] Failed to read the request body...")
	}

	reqt_body := string(bytes_reqt_body)
	log.Println("Echo is Hit! Body: " + reqt_body)

	fmt.Fprintf(http_resp, reqt_body)
}

// ====
// Main
// ====
func main () {
	log.Println("Starting rest-mocker...")
	mux_router := mux.NewRouter()

	mux_router.HandleFunc("/", homepageHandler)
	mux_router.HandleFunc("/echo", echoHandler)
	
	http.Handle("/", mux_router)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}