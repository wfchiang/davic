package davic 

import (
	// "fmt"
	"testing"
)

/********* 
Sample Data
*********/
func sampleJsonBytes0 () []byte {
	return []byte("{\"keyB\":false,\"keyI\":123,\"keyF\":1.23,\"keyS\":\"valS\",\"keyO\":{\"keykeyB\":true}}")
	// return []byte("{\"keyB\":false,\"keyI\":123,\"keyF\":1.23,\"keyS\":\"valS\",\"keyO\":{\"keykeyI\":456,\"keykeyO\":{\"keykeykeyI\",789}}}")
}

/********
Testing utilities
*********/
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

/********
Tests for semantics.go
*********/
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

func TestEvalExpr0 (t *testing.T) {
	expr := 1

	eval_result := EvalExpr(expr)

	if (expr != eval_result) {
		t.Error("")
	}
}

func TestEvalExpr1 (t *testing.T) {
	expr_0 := []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1, 1}
	eval_result_0 := EvalExpr(expr_0) 
	if (eval_result_0 != true) {
		t.Error("")
	}

	expr_1 := []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1, 2}
	eval_result_1 := EvalExpr(expr_1) 
	if (eval_result_1 != false) {
		t.Error("")
	}
}

/********
Tests for syntax.go
********/
func TestParseRefString0 (t *testing.T) {
	defer simpleRecover(t)

	key0 := "abc"
	key1 := "xyz"
	key_string := SYMBOL_REF_MARK + SYMBOL_REF_SEPARATOR + key0 + SYMBOL_REF_SEPARATOR + key1

	keys := ParseRefString(key_string) 

	if (len(keys) != 2) {
		t.Error("")
	}

	if (simpleIsViolation(TYPE_STRING, key0, keys[0])) {
		t.Error("")
	}

	if (simpleIsViolation(TYPE_STRING, key1, keys[1])) {
		t.Error("")
	}
}

func TestIsRef0 (t *testing.T) {
	good_ref_0 := []string{SYMBOL_REF_MARK, "abc"}
	good_ref_1 := []string{SYMBOL_REF_MARK}

	bad_ref_0 := []string{"abc"}

	if (simpleIsViolation(TYPE_BOOL, true, IsRef(good_ref_0))) {
		t.Error("")
	}

	if (simpleIsViolation(TYPE_BOOL, true, IsRef(good_ref_1))) {
		t.Error("")
	}

	if (simpleIsViolation(TYPE_BOOL, false, IsRef(bad_ref_0))) {
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

func TestMergeValidationResults (t *testing.T) {
	defer simpleRecover(t)

	var result0 ValidationResult
	result0.IsValid = true
	result0.Comments = []string{}

	var result1 ValidationResult 
	result1.IsValid = true 
	result1.Comments = []string{"Reason 1.0"}

	var result2 ValidationResult
	result2.IsValid = false 
	result2.Comments = []string{"Reason 2.0", "Reason 2.1"}

	result01 := MergeValidationResults(result0, result1)
	result12 := MergeValidationResults(result1, result2)

	if (!result01.IsValid) {
		t.Error("")
	}

	if (len(result01.Comments) != 1) {
		t.Error("")
	}

	if (result01.Comments[0] != "Reason 1.0") {
		t.Error("")
	}

	if (result12.IsValid) {
		t.Error("")
	}

	if (len(result12.Comments) != 3) {
		t.Error("")
	}

	if (result12.Comments[0] != "Reason 1.0" || result12.Comments[1] != "Reason 2.0" || result12.Comments[2] != "Reason 2.1") {
		t.Error("")
	}
}

/********
Tests for utils.go
********/
func TestContainsString (t *testing.T) {
	string_array := []string{"abc", "xyz"}

	if (!ContainsString(string_array, "abc")) {
		t.Error("")
	}

	if (ContainsString(string_array, "cba")) {
		t.Error("")
	}
}