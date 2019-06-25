package davic 

import (
//	"fmt"
	"testing"
)

/********* 
Sample Data
*********/
func sampleJsonBytes0 () []byte {
	return []byte("{\"keyN\":null,\"keyB\":false,\"keyI\":123,\"keyF\":1.23,\"keyS\":\"valS\",\"keyO\":{\"keykeyB\":true}}")
}

func sampleEnvironment0 () Environment {
	obj0 := CreateObjFromBytes(sampleJsonBytes0())
	var env0 = CreateNewEnvironment().PushStore(obj0)
	return env0
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

/********
Tests for syntax.go
********/
func TestIsType (t *testing.T) {
	defer simpleRecover(t) 

	if is_type := IsType(TYPE_BOOL, false); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_BOOL, 1); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_NUMBER, false); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_NUMBER, 1); is_type { // Only accept float as a number 
		t.Error("")
	}

	if is_type := IsType(TYPE_NUMBER, 1.1); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_STRING, 1); is_type {
		t.Error("")
	}
	
	if is_type := IsType(TYPE_STRING, "hello"); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_ARRAY, "hello"); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_ARRAY, []interface{}{123.0, false, "Hello"}); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_OBJ, 1); is_type {
		t.Error("")
	}

	obj0 := CreateObjFromBytes(sampleJsonBytes0())

	if is_type := IsType(TYPE_OBJ, obj0); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_NULL, obj0["keyN"]); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_NULL, obj0["no-such-key"]); !is_type {
		t.Error("")
	}
}

func TestParseRefString0 (t *testing.T) {
	defer simpleRecover(t)

	key0 := "abc"
	key1 := "xyz"
	key_string := SYMBOL_REF_MARK + SYMBOL_REF_SEPARATOR + key0 + SYMBOL_REF_SEPARATOR + key1

	keys := ParseRefString(key_string) 

	if (len(keys) != 3) {
		t.Error("")
	}

	if (simpleIsViolation(TYPE_STRING, SYMBOL_REF_MARK, keys[0])) {
		t.Error("")
	}

	if (simpleIsViolation(TYPE_STRING, key0, keys[1])) {
		t.Error("")
	}

	if (simpleIsViolation(TYPE_STRING, key1, keys[2])) {
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

func TestEnvironmentDeref0 (t *testing.T) {
	defer simpleRecover(t)

	obj0 := CreateObjFromBytes(sampleJsonBytes0())
	env0 := sampleEnvironment0()

	good_ref_0 := []string{SYMBOL_REF_MARK}
	good_ref_1 := []string{SYMBOL_REF_MARK, "keyB"}
	good_ref_2 := []string{SYMBOL_REF_MARK, "keyF"}
	
	if val := env0.Deref(good_ref_0) ; simpleIsViolation(TYPE_BOOL, true, IsType(TYPE_OBJ, val)) {
		t.Error("..")
	}
	
	if val := env0.Deref(good_ref_1) ; simpleIsViolation(TYPE_BOOL, obj0[good_ref_1[1]], val) {
		t.Error("")
	}

	if val := env0.Deref(good_ref_2) ; simpleIsViolation(TYPE_NUMBER, obj0[good_ref_2[1]], val) {
		t.Error("")
	}
}

func TestEnvironmentDeref1 (t *testing.T) {
	defer simpleRecover(t)

	obj0 := CreateObjFromBytes(sampleJsonBytes0())
	env0 := sampleEnvironment0()

	good_ref_0 := SYMBOL_REF_MARK
	good_ref_1 := SYMBOL_REF_MARK + SYMBOL_REF_SEPARATOR + "keyB"
	good_ref_2 := SYMBOL_REF_MARK + SYMBOL_REF_SEPARATOR + "keyF"

	if val := env0.Deref(good_ref_0) ; simpleIsViolation(TYPE_BOOL, true, IsType(TYPE_OBJ, val)) {
		t.Error("")
	}
	
	if val := env0.Deref(good_ref_1) ; simpleIsViolation(TYPE_BOOL, obj0["keyB"], val) {
		t.Error("")
	}

	if val := env0.Deref(good_ref_2) ; simpleIsViolation(TYPE_NUMBER, obj0["keyF"], val) {
		t.Error("")
	}
}

/********
Tests for semantics.go
*********/
func TestEvalExpr0 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := 1

	eval_result := EvalExpr(env, expr)

	if (expr != eval_result) {
		t.Error("")
	}
}

func TestEvalExpr1 (t *testing.T) {
	defer simpleExpectPanic(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1, 1}

	EvalExpr(env, expr) 
}

func TestEvalExpr2 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1.0, 1.0}
	eval_result := EvalExpr(env, expr) 
	if (eval_result != true) {
		t.Error("")
	}
}

func TestEvalExpr3 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1.0, 2.0}
	eval_result := EvalExpr(env, expr) 
	if (eval_result != false) {
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