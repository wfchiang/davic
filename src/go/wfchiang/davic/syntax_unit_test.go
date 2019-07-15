package davic 

import (
	"testing"
)

func TestIsType0 (t *testing.T) {
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

	if is_type := IsType(TYPE_ARRAY, []int{1, 2, 3}); is_type { // because the Go language unmarshaller should not give us this type: []int
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

	if is_type := IsType(TYPE_ARRAY, obj0["keyA"]); !is_type {
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

func TestEnvironmentStack0 (t *testing.T) {
	defer simpleRecover(t)

	env0 := sampleEnvironment0()
	stack_value := "123"
	env1 := env0.PushStack(stack_value)

	if env0.Stack.Len() != 0 {
		t.Error("")
	}

	if env1.Stack.Len() != 1 {
		t.Error("")
	}

	env2 := env1.PopStack() 
	
	if env1.Stack.Len() != 1 {
		t.Error("")
	}

	if env2.Stack.Len() != 0 {
		t.Error("")
	}

	env1_stack_value := env1.ReadStack()
	if env1_stack_value != "123" {
		t.Error("")
	}
}

func TestEnvironmentStack1 (t *testing.T) {
	defer simpleExpectPanic(t)

	env0 := sampleEnvironment0()
	env0.ReadStack()
}
