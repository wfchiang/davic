package davic 

import (
	"fmt"
	"strconv"
	"bytes"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

func ContainsString (string_array []string, target string) bool {
	for _, s := range string_array {
		if (strings.Compare(s, target) == 0) {
			return true 
		}
	}
	return false 
}

func GetKeyString (key []string) string {
	return strings.Join(key, " ")
}

/*
Casting utils  
*/
func CastNumberToInt (value float64) int {
	int_value := int(value)
	if (float64(int_value) != value) {
		panic("Inprecise CastNumberToInt")
	}
	return int_value
}

func CastInterfaceToBool (value interface{}) bool {
	bool_value, ok := value.(bool)
	if (!ok) {
		panic("CastInterfaceToBool failed")
	}
	return bool_value
}

func CastInterfaceToNumber (value interface{}) float64 {
	float64_value, ok := value.(float64)
	if (!ok) {
		panic("CastInterfaceToNumber failed")
	}
	return float64_value
}

func CastInterfaceToInt (value interface{}) int {
	int_value, ok := value.(int)
	if (!ok) {
		panic("CastInterfaceToInt failed")
	}
	return int_value
}

func CastInterfaceToString (value interface{}) string {
	string_value, ok := value.(string)
	if (!ok) {
		panic(fmt.Sprintf("CastInterfaceToString failed: %v", value))
	}
	return string_value
}

func CastInterfaceToArray (value interface{}) []interface{} {
	array_value, ok := value.([]interface{})
	if (!ok) {
		panic("CastInterfaceToArray failed")
	}
	return array_value
}

func CastInterfaceToObj (value interface{}) map[string]interface{} {
	obj_value, ok := value.(map[string]interface{})
	if (!ok) {
		panic("CastInterfaceToObj failed")
	}
	return obj_value
}

/* 
Marshaling
*/ 
func MarshalInterfaceToBytes (value interface{}) []byte {
	barr, err := json.Marshal(value)
	if (err != nil) {
		panic("MarshalInterfaceToBytes failed")
	}
	return barr
}

/*
Http Utils
*/ 
func MakeHttpCall (http_client *http.Client, obj_request map[string]interface{}) (obj_response map[string]interface{}) {
	if (http_client == nil) {
		panic("MakeHttpCall cannot proceed with nil http_client")
	}
	if _, ok := IsHttpRequest(obj_request); !ok {
		panic("MakeHttpCall cannot proceed with a non-http-request")
	}

	http_method        := CastInterfaceToString(obj_request[KEY_HTTP_METHOD])
	http_url           := CastInterfaceToString(obj_request[KEY_HTTP_URL])
	// http_headers := http_request[KEY_HTTP_HEADERS]
	http_body_reader   := bytes.NewReader(MarshalInterfaceToBytes(obj_request[KEY_HTTP_BODY]))

	// Create the http.Request obj 
	http_request, err := http.NewRequest(http_method, http_url, http_body_reader)
	if (err != nil) {
		panic("MakeHttpCall failed on http.NewRequest")
	}

	// TODO: Insert the http headers 

	// Make the call 
	http_response, err := http_client.Do(http_request)
	if (err != nil) {
		panic("MakeHttpCall failed on http.Client.Do")
	}

	// Get the response status 
	obj_response = map[string]interface{}{}
	obj_response[KEY_HTTP_STATUS]  = strconv.Itoa(http_response.StatusCode)

	// TODO: Get the response headers 
	obj_response[KEY_HTTP_HEADERS] = map[string]interface{}{}

	// Get the response body 
	defer http_response.Body.Close()
	bytes_resp_body, err := ioutil.ReadAll(http_response.Body)
	if (err != nil) {
		panic("MakeHttpCall failed on ioutil.ReadAll")
	}	
	obj_response[KEY_HTTP_BODY] = CreateObjFromBytes(bytes_resp_body)

	return obj_response
}