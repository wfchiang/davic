package davic 

import (
	"fmt"
	"testing"
	"net/http"
	"runtime/debug"
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
	return (expected != actual)
}

func simpleExpectPanic (t *testing.T) {
	r := recover() 
	if (r == nil) {
		t.Error("No expected panic occurred...")
	}
}

func simpleRecover (t *testing.T) {
	if r := recover(); r != nil {
		t.Error(fmt.Sprintf("Panic: %v.\nStack Trace: %v", r, string(debug.Stack())))
	}
}

func simpleTestingAssert(t *testing.T, type_string string, expected interface{}, actual interface{}) {
	if (!IsType(type_string, expected)) {
		return
	}
	if (!IsType(type_string, actual)) {
		panic(fmt.Sprintf("Invalid type of the actual value. Type %v is expected.", type_string))
	}
	if (expected != actual) {
		panic(fmt.Sprintf("Testing Assertion Failed. Expected vs Actual is [%v] vs [%v]", expected, actual))
	}
}

/********* 
Sample Data
*********/
func sampleJsonBytes0 () []byte {
	return []byte("{\"keyN\":null,\"keyB\":false,\"keyI\":123,\"keyF\":1.23,\"keyS\":\"valS\",\"keyO\":{\"keykeyB\":true},\"keyA\":[1, 2, 3, 1, 2, 3]}")
}

func sampleObj0 () map[string]interface{} {
	return CreateObjFromBytes(sampleJsonBytes0())
}

func sampleEnvironment0 () Environment {
	obj0 := sampleObj0()
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
	} else if (reqt_method == SYMBOL_HTTP_METHOD_POST && reqt_path == "/TestMakeHttpCall/post/0") {
		http_resp_writer.WriteHeader(200)
		fmt.Fprintf(http_resp_writer, string(sampleJsonBytes0()))
	} else {
		http_resp_writer.WriteHeader(404)
	}
}