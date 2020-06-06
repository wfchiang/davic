package davic 

import (
// 	"fmt"
	"testing"
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
	if (eval_result != true) {
		t.Error("")
	}

	expr = []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1.0, 1.0}
	eval_result = EvalExpr(env, expr) 
	if (eval_result != true) {
		t.Error("")
	}

	expr = []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, 1.0, 2.0}
	eval_result = EvalExpr(env, expr) 
	if (eval_result != false) {
		t.Error("")
	}

	expr = []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, "wfchiang", "wfchiang"}
	eval_result = EvalExpr(env, expr)
	if (eval_result != true) {
		t.Error("")
	}

	expr = []interface{}{SYMBOL_OPT_MARK, OPT_RELATION_EQ, "wfchiang", "Jenny"}
	eval_result = EvalExpr(env, expr)
	if (eval_result != false) {
		t.Error("")
	}
}

func TestEvalExprOptArithmeticAdd (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	eval_result := EvalExpr(env, expr) 
	if (eval_result != 6.0) {
		t.Error("")
	}

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 2.0, 3.0, 4.0}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 3.0, 4.0, 5.0}
	expr = []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, expr1, expr2, expr3}
	eval_result = EvalExpr(env, expr) 
	if (eval_result != 27.0) {
		t.Error("")
	}
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
	if (eval_result != -1.0) {
		t.Error("")
	}

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, expr1, 6.0}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, 2.0, expr1}
	expr = []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, expr2, expr3}
	eval_result = EvalExpr(env, expr) 
	if (eval_result != 4.0) {
		t.Error("")
	}
}

func TestEvalExprOptArithmeticMul (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_MUL, -2.0, 2.0}
	eval_result := EvalExpr(env, expr) 
	if (eval_result != -4.0) {
		t.Error("")
	}

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_MUL, expr1, 6.0}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, 2.0, expr1}
	expr = []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_MUL, expr2, expr3}
	eval_result = EvalExpr(env, expr) 
	if (eval_result != -144.0) {
		t.Error("")
	}
}

func TestEvalExprOptArithmeticDiv (t *testing.T) {
	defer simpleRecover(t)

	env := CreateNewEnvironment()
	expr := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_DIV, -2.0, 2.0}
	eval_result := EvalExpr(env, expr) 
	if (eval_result != -1.0) {
		t.Error("")
	}

	expr1 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_ADD, 1.0, 2.0, 3.0}
	expr2 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_DIV, expr1, 6.0}
	expr3 := []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_SUB, 2.0, expr1}
	expr = []interface{}{SYMBOL_OPT_MARK, OPT_ARITHMETIC_DIV, expr2, expr3}
	eval_result = EvalExpr(env, expr) 
	if (eval_result != -0.25) {
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
	if (eval_result != "123") {
		t.Error("")
	}
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
	if simpleIsViolation(TYPE_BOOL, true, arr_rel[0]) {
		t.Error("")
	}
	if simpleIsViolation(TYPE_BOOL, false, arr_rel[1]) {
		t.Error("")
	}
	if simpleIsViolation(TYPE_BOOL, false, arr_rel[2]) {
		t.Error("")
	}
	if simpleIsViolation(TYPE_BOOL, true, arr_rel[3]) {
		t.Error("")
	}
	if simpleIsViolation(TYPE_BOOL, false, arr_rel[4]) {
		t.Error("")
	}
	if simpleIsViolation(TYPE_BOOL, false, arr_rel[5]) {
		t.Error("")
	}
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
	if simpleIsViolation(TYPE_NUMBER, 1.0, arr_rel[0]) {
		t.Error("")
	}
	if simpleIsViolation(TYPE_NUMBER, 1.0, arr_rel[1]) {
		t.Error("")
	}
}

func TestEvalExpr17 (t *testing.T) {
	defer simpleRecover(t)

	env0 := sampleEnvironment0()
	opt_store_read := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_READ, "keyB"}

	val0 := EvalExpr(env0, opt_store_read)	
	if (simpleIsViolation(TYPE_BOOL, false, val0)) {
		t.Error("")
	}

	opt_store_write := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_WRITE, "keyB", true}
	env1 := EvalExpr(env0, opt_store_write).(Environment)
	val0 = EvalExpr(env0, opt_store_read)
	val1 := EvalExpr(env1, opt_store_read)
	if (simpleIsViolation(TYPE_BOOL, false, val0)) {
		t.Error("")
	}
	if (simpleIsViolation(TYPE_BOOL, true, val1)) {
		t.Error("")
	}
}

func TestEvalExpr18 (t *testing.T) {
	defer simpleRecover(t)

	env0 := sampleEnvironment0()
	opt_obj_read := []interface{}{SYMBOL_OPT_MARK, OPT_OBJ_READ, env0.Store, []interface{}{"keyO", "keykeyB"}}
	val0 := EvalExpr(env0, opt_obj_read)
	if (simpleIsViolation(TYPE_BOOL, true, val0)) {
		t.Error("")
	}
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
	if (simpleIsViolation(TYPE_BOOL, false, val0)) {
		t.Error("")
	}
	if (simpleIsViolation(TYPE_BOOL, true, val1["keykeyB"])) {
		t.Error("")
	}

	opt_swrite_0 := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_WRITE, "keyB", true}
	opt_swrite_1 := []interface{}{SYMBOL_OPT_MARK, OPT_STORE_WRITE, "keyO", map[string]interface{}{"wife":"Jenny"}}
	env1 := Execute(env0, []interface{}{opt_swrite_0, opt_swrite_1})
	
	val0 = CastInterfaceToBool(EvalExpr(env0, opt_sread_0))
	val1 = CastInterfaceToObj(EvalExpr(env0, opt_sread_1))
	if (simpleIsViolation(TYPE_BOOL, false, val0)) {
		t.Error("")
	}
	if (simpleIsViolation(TYPE_BOOL, true, val1["keykeyB"])) {
		t.Error("")
	}

	val2 := CastInterfaceToBool(EvalExpr(env1, opt_sread_0))
	val3 := CastInterfaceToObj(EvalExpr(env1, opt_sread_1))
	if (simpleIsViolation(TYPE_BOOL, true, val2)) {
		t.Error("")
	}
	if (simpleIsViolation(TYPE_STRING, "Jenny", val3["wife"])) {
		t.Error("")
	}
}