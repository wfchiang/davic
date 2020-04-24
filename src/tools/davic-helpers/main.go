package main

import (
	"fmt"
	"log"
	// "strings"
	"net/http"
	"encoding/json"
	"html/template"
	"github.com/gorilla/mux"
	"io/ioutil"
	"wfchiang/davic"
)

const (
	KEY_STORE_HTTP_REQUEST  = "http-reqt"
	KEY_STORE_HTTP_RESPONSE = "http-resp" 
)

type OptData struct {
	Name string 
	Symbol string 
	OpdNames []string
}

type OptList struct {
	SymbolOptMark string
	Operations []OptData 
}

// ====
// Globals 
// ====
var OptListData = OptList {
	SymbolOptMark: davic.SYMBOL_OPT_MARK, 
	Operations: []OptData {
		OptData {
			Name: "Lambda", 
			Symbol: davic.OPT_LAMBDA, 
			OpdNames: []string {"Function (an expression)"}},
		OptData {
			Name: "Function Call", 
			Symbol: davic.OPT_FUNC_CALL, 
			OpdNames: []string {"Lambda", "Value (an expression)"}},  
		OptData {
			Name: "Stack Read", 
			Symbol: davic.OPT_STACK_READ, 
			OpdNames: []string {}}, 
		OptData {
			Name: "Store Read", 
			Symbol: davic.OPT_STORE_READ, 
			OpdNames: []string {"Key (a string)"}}, 
		OptData {
			Name: "Store Write", 
			Symbol: davic.OPT_STORE_WRITE, 
			OpdNames: []string {"Key (a string)", "Value (an expression)"}},
		OptData {
			Name: "Relation Equal", 
			Symbol: davic.OPT_RELATION_EQ, 
			OpdNames: []string {"LHS Value (an expression)", "RHS Value (an expression)"}}, 
		OptData {
			Name: "Arithmetic Add", 
			Symbol: davic.OPT_ARITHMETIC_ADD, 
			OpdNames: []string {"LHS Value (an expression)", "RHS Value (an expression)"}}, 
		OptData {
			Name: "Arithmetic Subtract", 
			Symbol: davic.OPT_ARITHMETIC_SUB, 
			OpdNames: []string {"LHS Value (an expression)", "RHS Value (an expression)"}}, 
		OptData {
			Name: "Arithmetic Multiply", 
			Symbol: davic.OPT_ARITHMETIC_MUL, 
			OpdNames: []string {"LHS Value (an expression)", "RHS Value (an expression)"}}, 
		OptData {
			Name: "Arithmetic Division", 
			Symbol: davic.OPT_ARITHMETIC_DIV, 
			OpdNames: []string {"LHS Value (an expression)", "RHS Value (an expression)"}}, 
		OptData {
			Name: "Array Get/Read", 
			Symbol: davic.OPT_ARRAY_GET, 
			OpdNames: []string {"Array (an expression)", "Index (a number)"}}, 
		OptData {
			Name: "Object Read", 
			Symbol: davic.OPT_OBJ_READ,
			OpdNames: []string {"Object (an expression)", "Key (a string-array)"}} }}

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
	log.Println("Davic-helpers is Hit!")
	fmt.Fprintf(http_resp, "Davic Helpers are Here!")	
}

func optDataHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "opt-data")

	resp_body, err := json.Marshal(OptListData) 
	if err != nil {
		panic(fmt.Sprintf("Response marshalling failed: %v", err))
	} 
	
	fmt.Fprintf(http_resp, string(resp_body))
}

func optMakerHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "opt-maker")

	log.Println("Opt-maker is Hit!")
	
	template_fname := "opt-maker.html"
	tmpl, err := template.New(template_fname).Delims("<<", ">>").ParseFiles(template_fname)
	if (err != nil) {
		panic(fmt.Sprintf("Template load failed: %v", err))
	}

	tmpl.Execute(http_resp, nil)
	
	log.Println("Opt-maker responded")
}

func runDavicHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "run-davic")

	log.Println("Run-davic is Hit!")

	template_fname := "run-davic.html"
	tmpl, err := template.New(template_fname).Delims("<<", ">>").ParseFiles(template_fname)
	if (err != nil) {
		panic(fmt.Sprintf("Template load failed: %v", err))
	}

	tmpl.Execute(http_resp, nil)

	log.Println("Run-davic responded")
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
	file_server := http.FileServer(http.Dir("./static/"))

	log.Println("Starting Davic-helpers...")
	mux_router := mux.NewRouter()

	mux_router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", file_server))
	mux_router.HandleFunc("/opt-data", optDataHandler).Methods("GET")
	mux_router.HandleFunc("/opt-maker", optMakerHandler).Methods("GET")
	mux_router.HandleFunc("/run-davic", runDavicHandler).Methods("GET")
	mux_router.HandleFunc("/davic", davicHandler).Methods("POST")
	mux_router.HandleFunc("/", homepageHandler)

	http.Handle("/", mux_router)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}