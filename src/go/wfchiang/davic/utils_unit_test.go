package davic 

import (
//	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContainsString (t *testing.T) {
	string_array := []string{"abc", "xyz"}

	if (!ContainsString(string_array, "abc")) {
		t.Error("")
	}

	if (ContainsString(string_array, "cba")) {
		t.Error("")
	}
}

func mockTestingServerHandler (http_resp_writer http.ResponseWriter, http_request *http.Request) {
	reqt_method := http_request.Method 
	reqt_path   := http_request.URL.Path

	if (reqt_method == SYMBOL_HTTP_METHOD_GET && reqt_path == "/TestMakeHttpCall/0") {
		http_resp_writer.WriteHeader(200)
	} else {
		http_resp_writer.WriteHeader(404)
	}
}

func TestMakeHttpCall (t *testing.T) {
	defer simpleRecover(t)

	mock_http_server := httptest.NewServer(http.HandlerFunc(mockTestingServerHandler))
	if (mock_http_server == nil) {
		panic("Error: httptest.NewServer failed")
	}

	mock_http_client := mock_http_server.Client() 
	if (mock_http_client == nil) {
		panic("Error: failed to get the Client of the mock Server")
	}

	// Bad endpoint 
	obj_reqt := map[string]interface{}{}
	obj_reqt[KEY_HTTP_METHOD]  = SYMBOL_HTTP_METHOD_GET
	obj_reqt[KEY_HTTP_URL]     = mock_http_server.URL + "/BadEndpoint"
	obj_reqt[KEY_HTTP_HEADERS] = map[string]interface{}{}
	obj_reqt[KEY_HTTP_BODY]    = nil
	obj_resp := MakeHttpCall(mock_http_client, obj_reqt)
	if (obj_resp[KEY_HTTP_STATUS] != "404") {
		t.Error("")
	}

	// TestMakeHttpCall/0
	obj_reqt[KEY_HTTP_URL]     = mock_http_server.URL + "/TestMakeHttpCall/0"
	obj_resp = MakeHttpCall(mock_http_client, obj_reqt)
	if (obj_resp[KEY_HTTP_STATUS] != "200") {
		t.Error("")
	}
}