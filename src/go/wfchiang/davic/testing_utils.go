package davic 

import (
	"testing"
	"net/http"
)

/********
Util functions 
********/
func simpleIsViolation (type_string string, expected interface{}, actual interface{}) bool {
	if (!IsType(type_string, expected)) {
		return false
	}
	if (!IsType(type_string, actual)) {
		return true
	}
	return (expected != actual); 
}

func simpleExpectPanic (t *testing.T) {
	r := recover() 
	if (r == nil) {
		t.Error("No expected panic occurred...")
	}
}

func simpleRecover (t *testing.T) {
	if r := recover(); r != nil {
		t.Error("There was a panic... ", r)
	}
}

/********* 
Sample Data
*********/
func sampleJsonBytes0 () []byte {
	return []byte("{\"keyN\":null,\"keyB\":false,\"keyI\":123,\"keyF\":1.23,\"keyS\":\"valS\",\"keyO\":{\"keykeyB\":true},\"keyA\":[1, 2, 3, 1, 2, 3]}")
}

func sampleEnvironment0 () Environment {
	obj0 := CreateObjFromBytes(sampleJsonBytes0())
	env0 := CreateNewEnvironment()
	env0.Store = obj0
	return env0
}

/********
Mock Testing Handler
********/ 
func mockTestingServerHandler (http_resp_writer http.ResponseWriter, http_request *http.Request) {
	reqt_method := http_request.Method 
	reqt_path   := http_request.URL.Path

	if (reqt_method == SYMBOL_HTTP_METHOD_GET && reqt_path == "/TestMakeHttpCall/0") {
		http_resp_writer.WriteHeader(200)
	} else if (reqt_method == SYMBOL_HTTP_METHOD_GET && reqt_path == "/TestMakeHttpCall/1") {
		http_resp_writer.Header().Set("header1", "value1")
		http_resp_writer.WriteHeader(200)
	} else if (reqt_method == SYMBOL_HTTP_METHOD_GET && reqt_path == "/TestMakeHttpCall/2") {
		if hv := http_request.Header.Get("Header2"); hv == "value2" {
			http_resp_writer.WriteHeader(200)	
		} else {
			http_resp_writer.WriteHeader(400)
		}		
	} else {
		http_resp_writer.WriteHeader(404)
	}
}