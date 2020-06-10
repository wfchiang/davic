package davic 

import (
// 	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
)

/********
Tests of EvalExpr
********/
func TestEvalExpr0 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := 1

	eval_result := EvalExpr(env, expr)

	if (expr != eval_result) {
		t.Error("")
	}
}

func TestEvalExprOptRelationEq0 (t *testing.T) {
	defer simpleExpectPanic(t)

	env := CreateNewEnvironment()

	expr := []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1, 1}
	eval_result := EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_BOOL, true, eval_result)

	expr = []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1.0, 1.0}
	eval_result = EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_BOOL, true, eval_result)

	expr = []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1.0, 2.0}
	eval_result = EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_BOOL, false, eval_result)

	expr = []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, "wfchiang", "wfchiang"}
	eval_result = EvalExpr(env, expr)
	simpleTestingAssert(t, TYPE_BOOL, true, eval_result)

	expr = []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, "wfchiang", "Jenny"}
	eval_result = EvalExpr(env, expr)
	simpleTestingAssert(t, TYPE_BOOL, false, eval_result)
}

func TestEvalExprOptArithmeticAdd (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	eval_result := EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_NUMBER, 6.0, eval_result)

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 2.0, 3.0, 4.0}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 3.0, 4.0, 5.0}
	expr = []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, expr1, expr2, expr3}
	eval_result = EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_NUMBER, 27.0, eval_result)
}

func TestEvalExprOptArithmeticAddPanic0 (t *testing.T) {
	defer simpleExpectPanic(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0, "hello"}
	EvalExpr(env, expr) 
}

func TestEvalExprOptArithmeticSub (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, 1.0, 2.0}
	eval_result := EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_NUMBER, -1.0, eval_result)

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, expr1, 6.0}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, 2.0, expr1}
	expr = []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, expr2, expr3}
	eval_result = EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_NUMBER, 4.0, eval_result)
}

func TestEvalExprOptArithmeticMul (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_MUL, -2.0, 2.0}
	eval_result := EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_NUMBER, -4.0, eval_result)

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_MUL, expr1, 6.0}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, 2.0, expr1}
	expr = []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_MUL, expr2, expr3}
	eval_result = EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_NUMBER, -144.0, eval_result)
}

func TestEvalExprOptArithmeticDiv (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_DIV, -2.0, 2.0}
	eval_result := EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_NUMBER, -1.0, eval_result)

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_DIV, expr1, 6.0}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, 2.0, expr1}
	expr = []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_DIV, expr2, expr3}
	eval_result = EvalExpr(env, expr) 
	simpleTestingAssert(t, TYPE_NUMBER, -0.25, eval_result)
}

func TestEvalExpr7 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	env = env.PushStack("123")

	expr := []interface{}{SYMBOL_OPT_MARK, OPT_STACK_READ}
	eval_result := EvalExpr(env, expr)
	simpleTestingAssert(t, TYPE_STRING, "123", eval_result)
}

func TestEvalExpr8 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment() 

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARRAY_GET, []interface{}{1,2,3}, 1.0}
	eval_result := EvalExpr(env, expr1)
	if (eval_result != 2) {
		t.Error("")
	}

	obj2 := map[string]interface{}{"abc":1, "xyz":"123"}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_OBJ_READ, obj2, []interface{}{"abc"}}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_OBJ_READ, obj2, []interface{}{"xyz"}}
	
	eval_result = EvalExpr(env, expr2)
	if (eval_result != 1) {
		t.Error("")
	}

	eval_result = EvalExpr(env, expr3)
	simpleTestingAssert(t, TYPE_STRING, "123", eval_result)
}

func TestEvalExpr9 (t *testing.T) {
	defer simpleExpectPanic(t) 

	env := CreateNewEnvironment() 

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARRAY_GET, []int{1,2,3}, 1.0}
	eval_result := EvalExpr(env, expr1)
	if (eval_result != 2) {
		t.Error("")
	}
}

func TestEvalExpr10 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment() 

	lambda := []interface{}{SYMBOL_OPT_MARK, OPT_LAMBDA, 123}
	eval_result := EvalExpr(env, lambda) 
	original_lambda, is_lambda := IsLambdaOperation(eval_result)
	if !is_lambda {
		t.Error("")
	}
	if lambda[2] != original_lambda[2] {
		t.Error("")
	}
}

func TestEvalExpr11 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment() 

	stack_read := []interface{}{SYMBOL_OPT_MARK, OPT_STACK_READ}
	lambda := []interface{}{SYMBOL_OPT_MARK, OPT_LAMBDA, stack_read}

	func_call_0 := []interface{}{SYMBOL_OPT_MARK, OPT_FUNC_CALL, lambda, 123}
	if eval := EvalExpr(env, func_call_0); eval != 123 {
		t.Error("")
	}

	func_call_1 := []interface{}{SYMBOL_OPT_MARK, OPT_FUNC_CALL, lambda, "abc"}
	if eval := EvalExpr(env, func_call_1); eval != "abc" {
		t.Error("")
	}
}

func TestEvalExpr12 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment() 

	stack_read := []interface{}{SYMBOL_OPT_MARK, OPT_STACK_READ}
	stack_field_0 := []interface{}{SYMBOL_OPT_MARK, OPT_ARRAY_GET, stack_read, 0.0}
	stack_field_1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARRAY_GET, stack_read, 1.0}

	opt_add_stack := []interface{}{
		SYMBOL_OPT_MARK,
		OPT_ARITHMETIC_ADD, 
		stack_field_0, 
		stack_field_1}
	lambda := []interface{}{SYMBOL_OPT_MARK, OPT_LAMBDA, opt_add_stack}

	func_call_0 := []interface{}{
		SYMBOL_OPT_MARK, 
		OPT_FUNC_CALL, 
		lambda,
		[]interface{}{1.0, 2.0}}
	if eval := EvalExpr(env, func_call_0); eval != 3.0 {
		t.Error("")
	}

	func_call_1 := []interface{}{
		SYMBOL_OPT_MARK, 
		OPT_FUNC_CALL, 
		lambda, 
		[]interface{}{2.0, 2.0}}
	if eval := EvalExpr(env, func_call_1); eval != 4.0 {
		t.Error("")
	}
}

func TestEvalExpr13 (t *testing.T) {
	defer simpleExpectPanic(t)

	arr := []interface{}{1, 2, 3}
	env := sampleEnvironment0()
	opt_arr_get := []interface{}{SYMBOL_OPT_MARK, OPT_ARRAY_GET, arr, 1.2}
	EvalExpr(env, opt_arr_get)
}

func TestEvalExpr14 (t *testing.T) {
	defer simpleRecover(t)

	env := sampleEnvironment0()

	opt_deref := []interface{}{SYMBOL_OPT_MARK, SYMBOL_REF_MARK, "keyO", "keykeyB"}

	if val := CastInterfaceToBool(EvalExpr(env, opt_deref)); val != true {
		t.Error("")
	}
}

func TestEvalExpr15 (t *testing.T) {
	defer simpleRecover(t)

	env := sampleEnvironment0()
	arr_test := []interface{}{1.0, 2.0, 3.0, 1.0, 2.0, 3.0}

	stack_read := []interface{}{SYMBOL_OPT_MARK, OPT_STACK_READ}
	opt_eq     := []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1.0, stack_read} 
	opt_lambda := []interface{}{SYMBOL_OPT_MARK, OPT_LAMBDA, opt_eq}
	opt_map    := []interface{}{SYMBOL_OPT_MARK, OPT_ARRAY_MAP, arr_test, opt_lambda}

	rel := EvalExpr(env, opt_map)
	arr_rel := rel.([]interface{})
	if (len(arr_rel) != len(arr_test)) {
		t.Error("")
	}
	simpleTestingAssert(t, TYPE_BOOL, true,  arr_rel[0])
	simpleTestingAssert(t, TYPE_BOOL, false, arr_rel[1])
	simpleTestingAssert(t, TYPE_BOOL, false, arr_rel[2])
	simpleTestingAssert(t, TYPE_BOOL, true,  arr_rel[3])
	simpleTestingAssert(t, TYPE_BOOL, false, arr_rel[4])
	simpleTestingAssert(t, TYPE_BOOL, false, arr_rel[5])
}

func TestEvalExpr16 (t *testing.T) {
	defer simpleRecover(t)

	env := sampleEnvironment0()
	arr_test := []interface{}{1.0, 2.0, 3.0, 1.0, 2.0, 3.0}

	stack_read := []interface{}{SYMBOL_OPT_MARK, OPT_STACK_READ}
	opt_eq     := []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1.0, stack_read} 
	opt_lambda := []interface{}{SYMBOL_OPT_MARK, OPT_LAMBDA, opt_eq}
	opt_filter := []interface{}{SYMBOL_OPT_MARK, OPT_ARRAY_FILTER, arr_test, opt_lambda}

	rel := EvalExpr(env, opt_filter)
	arr_rel := rel.([]interface{})
	if (len(arr_rel) != 2) {
		t.Error("")
	}
	simpleTestingAssert(t, TYPE_NUMBER, 1.0, arr_rel[0])
	simpleTestingAssert(t, TYPE_NUMBER, 1.0, arr_rel[1])
}

func TestEvalExpr17 (t *testing.T) {
	defer simpleRecover(t)

	env0 := sampleEnvironment0()
	opt_store_read := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_READ, "keyB"}

	val0 := EvalExpr(env0, opt_store_read)	
	simpleTestingAssert(t, TYPE_BOOL, false, val0)

	opt_store_write := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_WRITE, "keyB", true}
	env1 := EvalExpr(env0, opt_store_write).(Environment)
	val0 = EvalExpr(env0, opt_store_read)
	val1 := EvalExpr(env1, opt_store_read)
	simpleTestingAssert(t, TYPE_BOOL, false, val0)
	simpleTestingAssert(t, TYPE_BOOL, true, val1)
}

func TestEvalExpr18 (t *testing.T) {
	defer simpleRecover(t)

	env0 := sampleEnvironment0()
	opt_obj_read := []interface{}{SYMBOL_OPT_MARK, OPT_OBJ_READ, env0.Store, []interface{}{"keyO", "keykeyB"}}
	val0 := EvalExpr(env0, opt_obj_read)
	simpleTestingAssert(t, TYPE_BOOL, true, val0)
}

/********
Named Tests of EvalExpr 
********/
func TestEvalExprHttpCall0 (t *testing.T) {
	defer simpleRecover(t)

	// Initialize the mock server 
	mock_http_server := httptest.NewServer(http.HandlerFunc(mockTestingServerHandler))
	if (mock_http_server == nil) {
		panic("Error: httptest.NewServer failed")
	}
	mock_http_client := mock_http_server.Client() 
	if (mock_http_client == nil) {
		panic("Error: failed to get the Client of the mock Server")
	}

	// Init the http request, the http call operation, and the environment
	http_reqt := map[string]interface{}{}
	http_call_opt := []interface{}{SYMBOL_OPT_MARK, OPT_HTTP_CALL, http_reqt}
	env := sampleEnvironment0()

	// Bad endpoint 
	http_reqt[KEY_HTTP_METHOD]  = SYMBOL_HTTP_METHOD_GET
	http_reqt[KEY_HTTP_URL]     = mock_http_server.URL + "/BadEndpoint"
	http_reqt[KEY_HTTP_HEADERS] = map[string]interface{}{}
	http_reqt[KEY_HTTP_BODY]    = nil
	http_resp := CastInterfaceToObj(EvalExpr(env, http_call_opt))
	simpleTestingAssert(t, TYPE_STRING, "404", http_resp[KEY_HTTP_STATUS])

	// TestMakeHttpCall/0
	http_reqt[KEY_HTTP_URL] = mock_http_server.URL + "/TestMakeHttpCall/0"
	http_resp = CastInterfaceToObj(EvalExpr(env, http_call_opt))
	simpleTestingAssert(t, TYPE_STRING, "200", http_resp[KEY_HTTP_STATUS])

	// TestMakeHttpCall/1 
	http_reqt[KEY_HTTP_URL]  = mock_http_server.URL + "/TestMakeHttpCall/1"
	http_resp = CastInterfaceToObj(EvalExpr(env, http_call_opt))
	
	simpleTestingAssert(t, TYPE_STRING, "200", http_resp[KEY_HTTP_STATUS])
	
	hv, ok := ReadHttpHeader(CastInterfaceToObj(http_resp[KEY_HTTP_HEADERS]), "header1")
	if (!ok || simpleIsViolation(TYPE_STRING, "value1", hv)) {
		t.Error("")
	}

	// TestMakeHttpCall/post/0 
	http_reqt[KEY_HTTP_METHOD] = SYMBOL_HTTP_METHOD_POST
	http_reqt[KEY_HTTP_URL] = mock_http_server.URL + "/TestMakeHttpCall/post/0"
	http_resp = CastInterfaceToObj(EvalExpr(env, http_call_opt))
	
	simpleTestingAssert(t, TYPE_STRING, "200", http_resp[KEY_HTTP_STATUS])
	
	http_resp_body := CastInterfaceToObj(http_resp[KEY_HTTP_BODY])
	simpleTestingAssert(t, TYPE_STRING, "valS", http_resp_body["keyS"]) 
}

/********
Tests of Execution 
********/ 
func TestExecution0 (t *testing.T) {
	defer simpleRecover(t)

	env0 := sampleEnvironment0()
	opt_sread_0 := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_READ, "keyB"}
	opt_sread_1 := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_READ, "keyO"}

	val0 := CastInterfaceToBool(EvalExpr(env0, opt_sread_0))
	val1 := CastInterfaceToObj(EvalExpr(env0, opt_sread_1))
	simpleTestingAssert(t, TYPE_BOOL, false, val0)
	simpleTestingAssert(t, TYPE_BOOL, true, val1["keykeyB"])

	opt_swrite_0 := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_WRITE, "keyB", true}
	opt_swrite_1 := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_WRITE, "keyO", map[string]interface{}{"wife":"Jenny"}}
	env1 := Execute(env0, []interface{}{opt_swrite_0, opt_swrite_1})
	
	val0 = CastInterfaceToBool(EvalExpr(env0, opt_sread_0))
	val1 = CastInterfaceToObj(EvalExpr(env0, opt_sread_1))
	simpleTestingAssert(t, TYPE_BOOL, false, val0)
	simpleTestingAssert(t, TYPE_BOOL, true, val1["keykeyB"])

	val2 := CastInterfaceToBool(EvalExpr(env1, opt_sread_0))
	val3 := CastInterfaceToObj(EvalExpr(env1, opt_sread_1))
	simpleTestingAssert(t, TYPE_BOOL, true, val2)
	simpleTestingAssert(t, TYPE_STRING, "Jenny", val3["wife"])
}

func TestExecution1 (t *testing.T) {
	defer simpleRecover(t)

	my_store := map[string]interface{}{}
	my_store["hero"] = "iron-man"
	my_store["power"] = "rich"

	env := CreateNewEnvironment()
	env.Store = my_store

	opt := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_WRITE, "hero", "bat-man"}

	env = Execute(env, []interface{}{opt})
	simpleTestingAssert(t, TYPE_STRING, "bat-man", env.Store["hero"])
}