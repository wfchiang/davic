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

func TestIsHttpOperation (t *testing.T) {
	http_request := map[string]interface{}{
		KEY_HTTP_METHOD:SYMBOL_HTTP_METHOD_GET,
		KEY_HTTP_URL:"http://127.0.0.1", 
		KEY_HTTP_HEADERS:map[string]interface{}{},
		KEY_HTTP_BODY:map[string]interface{}{},
	}
	http_opt := []interface{}{SYMBOL_OPT_MARK, OPT_HTTP_CALL, http_request}
	if _, is_http_opt := IsHttpOperation(http_opt); simpleIsViolation(TYPE_BOOL, true, is_http_opt) {
		t.Error("")
	}

	http_request[KEY_HTTP_METHOD] = SYMBOL_HTTP_METHOD_POST
	http_opt = []interface{}{SYMBOL_OPT_MARK, OPT_HTTP_CALL, http_request}
	if _, is_http_opt := IsHttpOperation(http_opt); simpleIsViolation(TYPE_BOOL, true, is_http_opt) {
		t.Error("")
	}

	http_request[KEY_HTTP_METHOD] = "not a valid method"
	http_opt = []interface{}{SYMBOL_OPT_MARK, OPT_HTTP_CALL, http_request}
	if _, is_http_opt := IsHttpOperation(http_opt); simpleIsViolation(TYPE_BOOL, false, is_http_opt) {
		t.Error("")
	}

	http_request[KEY_HTTP_METHOD] = SYMBOL_HTTP_METHOD_GET
	http_request[KEY_HTTP_HEADERS] = "not a valid headers"
	http_opt = []interface{}{SYMBOL_OPT_MARK, OPT_HTTP_CALL, http_request}
	if _, is_http_opt := IsHttpOperation(http_opt); simpleIsViolation(TYPE_BOOL, false, is_http_opt) {
		t.Error("")
	}

	http_request[KEY_HTTP_HEADERS] = map[string]interface{}{}
	http_request[KEY_HTTP_BODY] = "this is still correct for now"
	http_opt = []interface{}{SYMBOL_OPT_MARK, OPT_HTTP_CALL, http_request}
	if _, is_http_opt := IsHttpOperation(http_opt); simpleIsViolation(TYPE_BOOL, true, is_http_opt) {
		t.Error("")
	}

	delete(http_request, KEY_HTTP_BODY)
	http_opt = []interface{}{SYMBOL_OPT_MARK, OPT_HTTP_CALL, http_request}
	if _, is_http_opt := IsHttpOperation(http_opt); simpleIsViolation(TYPE_BOOL, false, is_http_opt) {
		t.Error("")
	}
}

func TestIsHttpResponse (t *testing.T) {
	defer simpleRecover(t) 

	http_res := map[string]interface{}{
		KEY_HTTP_STATUS:"200", 
		KEY_HTTP_HEADERS:map[string]interface{}{}, 
		KEY_HTTP_BODY:map[string]interface{}{}}
	if _, is_http_res := IsHttpResponse(http_res); !is_http_res {
		t.Error("")
	}

	http_res = map[string]interface{}{
		KEY_HTTP_STATUS:"not a valid status", 
		KEY_HTTP_HEADERS:map[string]interface{}{}, 
		KEY_HTTP_BODY:map[string]interface{}{}}
	if _, is_http_res := IsHttpResponse(http_res); is_http_res {
		t.Error("")
	}

	http_res = map[string]interface{}{
		KEY_HTTP_STATUS:"200", 
		KEY_HTTP_HEADERS:"not a valid headers",  
		KEY_HTTP_BODY:map[string]interface{}{}}
	if _, is_http_res := IsHttpResponse(http_res); is_http_res {
		t.Error("")
	}

	http_res = map[string]interface{}{
		KEY_HTTP_STATUS:"200", 
		KEY_HTTP_HEADERS:map[string]interface{}{}, 
		KEY_HTTP_BODY:"still a valid body"}
	if _, is_http_res := IsHttpResponse(http_res); !is_http_res {
		t.Error("")
	}

	http_res = map[string]interface{}{
		KEY_HTTP_STATUS:200, 
		KEY_HTTP_HEADERS:map[string]interface{}{}, 
		KEY_HTTP_BODY:map[string]interface{}{}}
	if _, is_http_res := IsHttpResponse(http_res); is_http_res {
		t.Error("")
	}

	http_res = map[string]interface{}{
		KEY_HTTP_STATUS:200.0, 
		KEY_HTTP_HEADERS:map[string]interface{}{}, 
		KEY_HTTP_BODY:map[string]interface{}{}}
	if _, is_http_res := IsHttpResponse(http_res); is_http_res {
		t.Error("")
	}

	http_res = map[string]interface{}{
		KEY_HTTP_HEADERS:map[string]interface{}{}, 
		KEY_HTTP_BODY:map[string]interface{}{}}
	if _, is_http_res := IsHttpResponse(http_res); is_http_res {
		t.Error("")
	}

	http_res = map[string]interface{}{
		KEY_HTTP_STATUS:"200",  
		KEY_HTTP_BODY:map[string]interface{}{}}
	if _, is_http_res := IsHttpResponse(http_res); is_http_res {
		t.Error("")
	}

	http_res = map[string]interface{}{
		KEY_HTTP_STATUS:"200", 
		KEY_HTTP_HEADERS:map[string]interface{}{}}
	if _, is_http_res := IsHttpResponse(http_res); is_http_res {
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

func TestEnvironmentClone (t *testing.T) {
	defer simpleRecover(t)

	// Scenario of empty stack 
	env0 := sampleEnvironment0()
	if env0.Stack == nil || env0.Stack.Len() != 0 {
		t.Error("")
	}
	clone_env0 := env0.Clone()
	if clone_env0.Stack == nil || env0.Stack.Len() != 0 {
		t.Error("")
	}

	// Non-empty stack scenario 
	env0 = env0.PushStack(env0.Store)
	if env0.Stack == nil || env0.Stack.Len() != 1 {
		t.Error("")
	}
	clone_env0 = env0.Clone()
	if clone_env0.Stack == nil || env0.Stack.Len() != 1 {
		t.Error("")
	}

	// Nil stack scenario -->> even though this scenario is not expected ... 
	env0.Stack = nil
	if env0.Stack != nil {
		t.Error("")
	}
	clone_env0 = env0.Clone()
	if clone_env0.Stack != nil {
		t.Error("")
	}
}
