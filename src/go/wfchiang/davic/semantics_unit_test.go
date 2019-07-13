package davic 

import (
//	"fmt"
	"testing"
)

/********* 
Sample Data
*********/
func sampleJsonBytes0 () []byte {
	return []byte("{\"keyN\":null,\"keyB\":false,\"keyI\":123,\"keyF\":1.23,\"keyS\":\"valS\",\"keyO\":{\"keykeyB\":true},\"keyA\":[1, 2, 3]}")
}

func sampleEnvironment0 () Environment {
	obj0 := CreateObjFromBytes(sampleJsonBytes0())
	env0 := CreateNewEnvironment()
	env0.Store = obj0
	return env0
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

func TestEvalExpr4 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	eval_result := EvalExpr(env, expr) 
	if (eval_result != 6.0) {
		t.Error("")
	}
}

func TestEvalExpr5 (t *testing.T) {
	defer simpleExpectPanic(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0, "hello"}
	EvalExpr(env, expr) 
}

func TestEvalExpr6 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 2.0, 3.0, 4.0}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 3.0, 4.0, 5.0}
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, expr1, expr2, expr3}
	eval_result := EvalExpr(env, expr) 
	if (eval_result != 27.0) {
		t.Error("")
	}
}

func TestEvalExpr7 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	env = env.PushStack("123")

	expr := []interface{}{SYMBOL_OPT_MARK, OPT_STACK_READ}
	eval_result := EvalExpr(env, expr)
	if (eval_result != "123") {
		t.Error("")
	}
}

func TestEvalExpr8 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment() 

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_FIELD_READ, []interface{}{1,2,3}, 1.0}
	eval_result := EvalExpr(env, expr1)
	if (eval_result != 2) {
		t.Error("")
	}

	obj2 := map[string]interface{}{"abc":1, "xyz":"123"}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_FIELD_READ, obj2, "abc"}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_FIELD_READ, obj2, "xyz"}
	
	eval_result = EvalExpr(env, expr2)
	if (eval_result != 1) {
		t.Error("")
	}

	eval_result = EvalExpr(env, expr3)
	if (eval_result != "123") {
		t.Error("")
	}
}

func TestEvalExpr9 (t *testing.T) {
	defer simpleExpectPanic(t) 

	env := CreateNewEnvironment() 

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_FIELD_READ, []int{1,2,3}, 1.0}
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

	func_call_0 := []interface{}{SYMBOL_OPT_MARK, OPT_FUNC_CALL, 123, lambda}
	if eval := EvalExpr(env, func_call_0); eval != 123 {
		t.Error("")
	}

	func_call_1 := []interface{}{SYMBOL_OPT_MARK, OPT_FUNC_CALL, "abc", lambda}
	if eval := EvalExpr(env, func_call_1); eval != "abc" {
		t.Error("")
	}
}

func TestEvalExpr12 (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment() 

	stack_read := []interface{}{SYMBOL_OPT_MARK, OPT_STACK_READ}
	stack_field_0 := []interface{}{SYMBOL_OPT_MARK, OPT_FIELD_READ, stack_read, 0.0}
	stack_field_1 := []interface{}{SYMBOL_OPT_MARK, OPT_FIELD_READ, stack_read, 1.0}

	opt_add_stack := []interface{}{
		SYMBOL_OPT_MARK,
		OPT_ARITHMETIC_ADD, 
		stack_field_0, 
		stack_field_1}
	lambda := []interface{}{SYMBOL_OPT_MARK, OPT_LAMBDA, opt_add_stack}

	func_call_0 := []interface{}{
		SYMBOL_OPT_MARK, 
		OPT_FUNC_CALL, 
		[]interface{}{1.0, 2.0}, 
		lambda}
	if eval := EvalExpr(env, func_call_0); eval != 3.0 {
		t.Error("")
	}

	func_call_1 := []interface{}{
		SYMBOL_OPT_MARK, 
		OPT_FUNC_CALL, 
		[]interface{}{2.0, 2.0}, 
		lambda}
	if eval := EvalExpr(env, func_call_1); eval != 4.0 {
		t.Error("")
	}
}