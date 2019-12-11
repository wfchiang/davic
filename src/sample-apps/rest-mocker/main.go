package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"wfchiang/davic"
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

func hi2youHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body)
	if err != nil {
		log.Println("[Echo] Failed to read the request body...")
	}

	reqt_body := string(bytes_reqt_body)
	log.Println("Hi2You is Hit! Body: " + reqt_body)

	// Initialize the Davic environment 
	reqt_obj := davic.CreateObjFromBytes(bytes_reqt_body)
	env := davic.CreateNewEnvironment()
	env.Store = map[string]interface{}{"reqt":reqt_obj}

	// Prepare the operation 
	opt_store_read := []interface{}{davic.SYMBOL_OPT_MARK, davic.OPT_STORE_READ, "reqt"}
	opt_obj_read := []interface{}{davic.SYMBOL_OPT_MARK, davic.OPT_OBJ_READ, opt_store_read, []interface{}{"name"}}
	opt_store_write := []interface{}{davic.SYMBOL_OPT_MARK, davic.OPT_STORE_WRITE, "name", opt_obj_read}
	env = davic.Execute(env, []interface{}{opt_store_write})
	log.Println(fmt.Sprintf("[Hi2You] name is %v", env.Store["name"]))

	obj_resp := map[string]interface{}{"name": env.Store["name"], "message": "Hi!"}
	resp_body, err := json.Marshal(obj_resp) 
	if err != nil {
		log.Println(fmt.Sprintf("[Hi2You] response marshalling failed: %v", err))
		fmt.Fprintf(http_resp, "Sorry... I failed")
	} else {
		fmt.Fprintf(http_resp, string(resp_body))
	}
}

// ====
// Main
// ====
func main () {
	log.Println("Starting rest-mocker...")
	mux_router := mux.NewRouter()

	mux_router.HandleFunc("/", homepageHandler)
	mux_router.HandleFunc("/echo", echoHandler)
	mux_router.HandleFunc("/hi2you", hi2youHandler)
	
	http.Handle("/", mux_router)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}