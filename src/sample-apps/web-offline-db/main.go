package main

import (
	"fmt"
	"log"
	// "strings"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"wfchiang/davic"
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
	defer recoverFromPanic(http_resp, "davic")

	log.Println("Davic is Hit!")

	// Read the request body 
	bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body)
	if err != nil {
		panic("Failed to read the request body")
	}

	reqt_body := string(bytes_reqt_body)
	log.Println("Davic/Go is Hit! Body: " + reqt_body)

	// Convert the string type request body to object
	reqt_obj := davic.CreateObjFromBytes(bytes_reqt_body)
	
	_, ok := reqt_obj["data"]
	if (!ok) {
		panic("Data field is missed in the request object")
	}
	opt_obj, ok := reqt_obj["opt"]
	if (!ok) {
		panic("Opt field is missed in the request object")
	}

	// Setup the Davic environment 
	env := davic.CreateNewEnvironment()
	env.Store = davic.CastInterfaceToObj(reqt_obj)

	// Execute the operation 
	rel_obj := davic.EvalExpr(env, opt_obj)

	// Marshal the response 
	resp_body, err := json.Marshal(rel_obj) 
	if err != nil {
		panic(fmt.Sprintf("Response marshalling failed: %v", err))
	} 
	
	fmt.Fprintf(http_resp, string(resp_body))
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