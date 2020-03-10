package main

import (
	"fmt"
	"log"
	// "strings"
	"net/http"
	// "encoding/json"
	"github.com/gorilla/mux"
	// "io/ioutil"
	// "wfchiang/davic"
)

const (
	KEY_STORE_HTTP_REQUEST  = "http-reqt"
	KEY_STORE_HTTP_RESPONSE = "http-resp"
)

// ====
// Utils 
// ====

// ==== 
// Recovery function 
// ====
func recoverFromPanic (http_resp http.ResponseWriter, id_service string) {
	if r := recover(); r != nil {
		err_message := fmt.Sprintf("%v", r)
		log.Println(fmt.Sprintf("[%s] %s", id_service, err_message))
		fmt.Fprintf(http_resp, err_message)
	}
}

// ====
// Handlers 
// ====
func homepageHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	log.Println("Homepage is Hit!")
	fmt.Fprintf(http_resp, "Web-Offline-DB is Here!")
}

func davicHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	log.Println("Davic is Hit!")
	fmt.Fprintf(http_resp, "Yo")
}

// ====
// Main
// ====
func main () {
	log.Println("Init File Server...")
	file_server := http.FileServer(http.Dir("./static"))

	log.Println("Starting Web-Offline-DB...")
	mux_router := mux.NewRouter()

	mux_router.PathPrefix("/pages/").Handler(http.StripPrefix("/pages/", file_server))
	mux_router.HandleFunc("/davic", davicHandler).Methods("POST")
	mux_router.HandleFunc("/", homepageHandler)

	http.Handle("/", mux_router)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}