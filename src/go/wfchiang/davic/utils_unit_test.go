package davic 

import (
	"fmt"
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

func TestGetValue0 (t *testing.T) {
	defer simpleRecover(t) 

	obj := CreateObjFromBytes(sampleJsonBytes0())

	if val := GetObjValue(obj,[]string{"keyN"}); simpleIsViolation(TYPE_NULL, nil, val) {
		t.Error("")
	}

	if val := GetObjValue(obj,[]string{"no-such-key"}); simpleIsViolation(TYPE_NULL, nil, val) {
		t.Error("")
	}

	if val := GetObjValue(obj,[]string{"keyB"}); simpleIsViolation(TYPE_BOOL, false, val) {
		t.Error("") 
	}

	if val := GetObjValue(obj,[]string{"keyI"}); simpleIsViolation(TYPE_BOOL, true, IsType(TYPE_NUMBER, val)) {
		t.Error("")
	}
	if val := GetObjValue(obj,[]string{"keyI"}); simpleIsViolation(TYPE_NUMBER, 123, val) {
		t.Error("")
	}
	
	if val := GetObjValue(obj,[]string{"keyF"}); simpleIsViolation(TYPE_NUMBER, 1.23, val) {
		t.Error("")
	}

	if val := GetObjValue(obj,[]string{"keyS"}); simpleIsViolation(TYPE_STRING, "valS", val) {
		t.Error("")
	}
	
	if val := GetObjValue(obj,[]string{"keyO"}); (!IsType(TYPE_OBJ, val)) {
		t.Error("")
	}

	if val := GetObjValue(obj,[]string{"keyO", "keykeyB"}); simpleIsViolation(TYPE_BOOL, true, val) {
		t.Error("")
	}
}

func TestTestGetValue1 (t *testing.T) {
	defer simpleExpectPanic(t)

	obj := CreateObjFromBytes(sampleJsonBytes0())
	GetObjValue(obj,[]string{"no-such-key","no-more-such-key"})
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
	obj_resp                   = MakeHttpCall(mock_http_client, obj_reqt)
	if (obj_resp[KEY_HTTP_STATUS] != "200") {
		t.Error("")
	}

	// TestMakeHttpCall/1 
	obj_reqt[KEY_HTTP_URL]     = mock_http_server.URL + "/TestMakeHttpCall/1"
	obj_resp                   = MakeHttpCall(mock_http_client, obj_reqt) 
	obj_resp_headers          := CastInterfaceToObj(obj_resp[KEY_HTTP_HEADERS])
	if hv, ok := ReadHttpHeader(obj_resp_headers, "header1"); (obj_resp[KEY_HTTP_STATUS] != "200" || !ok || hv != "value1") {
		t.Error("")
	}
	if hv, ok := ReadHttpHeader(obj_resp_headers, "bad-header1"); (obj_resp[KEY_HTTP_STATUS] != "200" || ok) {
		t.Error(fmt.Sprintf("Impossible value of bad-header1: %v", hv))
	}
}