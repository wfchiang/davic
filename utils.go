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

func CastInterfaceToStringArray (value interface{}) []string {
	arr := CastInterfaceToArray(value)
	var str_arr	[]string 
	for _, s := range arr {
		str_arr = append(str_arr, CastInterfaceToString(s)) 
	}
	return str_arr
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
Copy Value
*/ 
func CopyValue (value interface{}) interface{} {
	if (IsType(TYPE_NULL, value)) {
		return nil
	} else if (IsType(TYPE_BOOL, value) || IsType(TYPE_NUMBER, value) || IsType(TYPE_STRING, value)) {
		return value
	} else if (IsType(TYPE_ARRAY, value)) {
		origin_array, _ := value.([]interface{})
		copy_array := []interface{}{}
		for _, v := range origin_array {
			copy_array = append(copy_array, CopyValue(v))
		}
		return copy_array
	} else if (IsType(TYPE_OBJ, value)) {
		origin_obj, _ := value.(map[string]interface{})
		copy_obj := map[string]interface{}{}
		for k, v := range origin_obj {
			copy_obj[k] = CopyValue(v)
		}
		return copy_obj
	} else {
		error_message := fmt.Sprintf("Cannot copy value with unknown type. Value: %v", value)
		panic(error_message)
	}
}

/* 
Object Utils 
*/ 
func GetObjKeys (obj map[string]interface{}) []string {
	var keys []string
	for k, _ := range obj {
		keys = append(keys, k)
	}
	return keys
}

func ReadObjValue (obj map[string]interface{}, key []string) interface{} {
	if (len(key) == 0) {
		return obj
	}

	var retv interface{} = nil
	
	for i, k := range key {
		next_obj, ok := obj[k] 
		if (!ok) {
			panic(fmt.Sprintf("ReadObjValue failed because the key %v does not exist", key))
		}
		if (i == (len(key)-1)) {
			retv = obj[k]
			break
		} else {
			obj = CastInterfaceToObj(next_obj)
		}
	}

	return retv
}

func UpdateObjValue (obj map[string]interface{}, key []string, val interface{}) map[string]interface{} {
	if (len(key) == 0) {
		return obj
	}

	new_val := CopyValue(val)
	new_obj := CastInterfaceToObj(CopyValue(obj))
	traversal := new_obj
	
	for i, k := range key {
		next_obj, ok := traversal[k]
		if (!ok) {
			panic(fmt.Sprintf("UpdateObjValue failed because the key %v does not exist", key))
		} 
		if (i == (len(key)-1)) {
			traversal[k] = new_val
			break
		} else {
			traversal = CastInterfaceToObj(next_obj)
		}
	}

	return new_obj
}

func WriteObjValue (obj map[string]interface{}, key []string, val interface{}) map[string]interface{} {
	if (len(key) == 0) {
		return obj
	}

	new_val := CopyValue(val)
	new_obj := CastInterfaceToObj(CopyValue(obj))
	traversal := new_obj
	
	for i, k := range key {
		next_obj, ok := traversal[k] 
		if (!ok) {
			traversal[k] = map[string]interface{}{}
			next_obj, ok = traversal[k] 
			if (!ok) {
				panic(fmt.Sprintf("WriteObjValue failed because the key %v cannot be automatically generated", key))
			}
		} 
		if (i == (len(key)-1)) {
			traversal[k] = new_val
			break
		} else {
			traversal = CastInterfaceToObj(next_obj)
		}
	}

	return new_obj
}

/*
Http Utils
*/ 
func ReadHttpHeader (obj_resp_header map[string]interface{}, header_key string) (interface{}, bool) {
	header_keys := GetObjKeys(obj_resp_header)
	for _, k := range header_keys {
		if (strings.ToLower(k) == strings.ToLower(header_key)) {
			return obj_resp_header[k], true
		}
	}
	return nil, false
}

func MakeHttpCall (http_client *http.Client, obj_request map[string]interface{}) (obj_response map[string]interface{}) {
	if (http_client == nil) {
		panic("MakeHttpCall cannot proceed with nil http_client")
	}
	if _, ok := IsHttpRequest(obj_request); !ok {
		panic("MakeHttpCall cannot proceed with a non-http-request")
	}

	http_method      := CastInterfaceToString(obj_request[KEY_HTTP_METHOD])
	http_url         := CastInterfaceToString(obj_request[KEY_HTTP_URL])
	http_headers     := CastInterfaceToObj(obj_request[KEY_HTTP_HEADERS])
	http_body_reader := bytes.NewReader(MarshalInterfaceToBytes(obj_request[KEY_HTTP_BODY]))

	// Create the http.Request obj 
	http_request, err := http.NewRequest(http_method, http_url, http_body_reader)
	if (err != nil) {
		panic("MakeHttpCall failed on http.NewRequest")
	}

	// Insert the http headers 
	for hkey, kval := range http_headers {
		if _, ok := http_request.Header[hkey]; ok {
			http_request.Header[hkey] = append(http_request.Header[hkey], CastInterfaceToString(kval))
		} else {
			http_request.Header[hkey] = []string {CastInterfaceToString(kval)}
		}
	}

	// Make the call 
	http_response, err := http_client.Do(http_request)
	if (err != nil) {
		panic("MakeHttpCall failed on http.Client.Do")
	}

	// Get the response status 
	obj_response = map[string]interface{}{}
	obj_response[KEY_HTTP_STATUS] = strconv.Itoa(http_response.StatusCode)

	// Get the response headers 
	obj_headers := map[string]interface{}{}
	for hkey, hval := range http_response.Header {
		str_hval := strings.Join(hval,";")
		if _, ok := obj_headers[hkey]; ok {
			obj_headers[hkey] = CastInterfaceToString(obj_headers[hkey]) + str_hval
		} else {
			obj_headers[hkey] = str_hval
		}
	}
	obj_response[KEY_HTTP_HEADERS] = obj_headers

	// Get the response body 
	defer http_response.Body.Close()
	bytes_resp_body, err := ioutil.ReadAll(http_response.Body)
	if (err != nil) {
		panic("MakeHttpCall failed on ioutil.ReadAll")
	}	
	obj_response[KEY_HTTP_BODY] = CreateObjFromBytes(bytes_resp_body)

	return obj_response
}