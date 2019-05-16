package davic 

import (
	// "fmt"
	"testing"
)

func sampleJsonBytes0 () []byte {
	return []byte("{\"keyB\":false,\"keyI\":123,\"keyF\":1.23,\"keyS\":\"valS\",\"keyO\":{\"keykeyB\":true}}")
	// return []byte("{\"keyB\":false,\"keyI\":123,\"keyF\":1.23,\"keyS\":\"valS\",\"keyO\":{\"keykeyI\":456,\"keykeyO\":{\"keykeykeyI\",789}}}")
}

func simpleIsViolation (type_string string, expected interface{}, actual interface{}) bool {
	if (!IsType(type_string, expected)) {
		return false
	}
	if (!IsType(type_string, actual)) {
		return true
	}
	return (expected != actual); 
}

func simpleRecover (t *testing.T) {
	if r := recover(); r != nil {
		t.Error("There was a panic... ", r)
	}
}

func TestIsType (t *testing.T) {
	defer simpleRecover(t) 

	if is_type := IsType(TYPE_BOOL, false); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_BOOL, 1); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_INT, false); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_INT, 1); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_FLOAT, 1); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_FLOAT, 1.1); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_STRING, 1); is_type {
		t.Error("")
	}
	
	if is_type := IsType(TYPE_STRING, "hello"); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_OBJ, 1); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_OBJ, CreateObjFromBytes(sampleJsonBytes0())); !is_type {
		t.Error("")
	}
}

func TestGetValue0 (t *testing.T) {
	defer simpleRecover(t) 

	obj := CreateObjFromBytes(sampleJsonBytes0())

	if val := GetObjValue(obj,[]string{"keyB"}); simpleIsViolation(TYPE_BOOL, false, val) {
		t.Error("") 
	}

	if val := GetObjValue(obj,[]string{"keyI"}); simpleIsViolation(TYPE_BOOL, true, IsType(TYPE_FLOAT, val)) {
		t.Error("")
	}
	if val := GetObjValue(obj,[]string{"keyI"}); simpleIsViolation(TYPE_FLOAT, 123, val) {
		t.Error("")
	}
	
	if val := GetObjValue(obj,[]string{"keyF"}); simpleIsViolation(TYPE_FLOAT, 1.23, val) {
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