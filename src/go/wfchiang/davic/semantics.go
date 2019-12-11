package davic 

import (
	"fmt"
	"strings"
)

/********
Expression evaluation
********/ 
func EvalExpr (env Environment, in_expr interface{}) interface{} {
	// No need to evaluate lambda 
	if _, is_lambda := IsLambdaOperation(in_expr); is_lambda {
		return in_expr
	}

	// Try to evaluate an operation -- an operation has its form...
	operation, ok := IsOperation(in_expr)
	if (!ok) {
		return in_expr
	} 

	opt, ok := operation[1].(string)
	if (!ok) {
		panic("Operator (" + fmt.Sprintf("%v", opt) + ")is not a string")
	}

	if (strings.Compare(OPT_STACK_READ, opt) == 0) {
		// should NOT evaluate the stack value 
		return env.ReadStack()

	} else if (strings.Compare(OPT_STORE_READ, opt) == 0) {
		operation, ok = IsUnaryOperation(in_expr)
		if (!ok) {
			panic(fmt.Sprintf("Invalid operation: %v", operation))
		}
		opd := EvalExpr(env, operation[2])
		if (!IsType(TYPE_STRING, opd)) {
			panic(fmt.Sprintf("Invalid type of the store-read operation: %v", operation))
		}
		str_opd := CastInterfaceToString(opd)
		return env.ReadStore(str_opd)

	} else if (strings.Compare(OPT_STORE_WRITE, opt) == 0) {
		operation, ok = IsBinaryOperation(in_expr)
		if (!ok) {
			panic(fmt.Sprintf("Invalid operation: %v", operation))
		}
		vkey := EvalExpr(env, operation[2])
		sval := EvalExpr(env, operation[3])
		if (!IsType(TYPE_STRING, vkey)) {
			panic(fmt.Sprintf("Invalid type of the store-write operation (the key must be a string)"))
		}
		str_vkey := CastInterfaceToString(vkey)
		// Here we will return an Environment -- since a Store-Write operation will change the environment... 
		return env.WriteStore(str_vkey, sval)
	
	} else if (strings.Compare(OPT_RELATION_EQ, opt) == 0) {
		operation, ok = IsBinaryOperation(in_expr)
		if (!ok) {
			panic(fmt.Sprintf("Invalid operation: %v", operation))
		}
		lhs := EvalExpr(env, operation[2])
		rhs := EvalExpr(env, operation[3])

		if (IsType(TYPE_BOOL, lhs) && IsType(TYPE_BOOL, rhs)) {
			return (lhs == rhs)
		} else if (IsType(TYPE_NUMBER, lhs) && IsType(TYPE_NUMBER, rhs)) {
			return (lhs == rhs)
		} else {
			panic(fmt.Sprintf("Invalid operation: %v", operation))
		}

	} else if (strings.Compare(OPT_ARITHMETIC_ADD, opt) == 0) {
		if (len(operation) < 3) {
			panic(fmt.Sprintf("Invalid operation: %v", operation))
		}
		
		add_result := CastInterfaceToNumber(EvalExpr(env, operation[2]))
		for _, v := range operation[3:] {
			add_result = add_result + CastInterfaceToNumber(EvalExpr(env, v))
		}
		
		return add_result

	} else if (strings.Compare(OPT_FIELD_READ, opt) == 0) {
		if (len(operation) != 4) {
			panic(fmt.Sprintf("Invalid field-read operation: %v : %s", operation, "mal-form"))
		}

		container := EvalExpr(env, operation[2])
		key := EvalExpr(env, operation[3])

		if (IsType(TYPE_OBJ, container) && IsType(TYPE_STRING, key)) { // If the container is an object ...
			typed_container := CastInterfaceToObj(container)
			typed_key := CastInterfaceToString(key)
			return typed_container[typed_key]

		} else if (IsType(TYPE_ARRAY, container) && IsType(TYPE_NUMBER, key)) { // If the container is an array ... 
			typed_container := CastInterfaceToArray(container)

			float64_key := CastInterfaceToNumber(key)
			typed_key := int(float64_key)
			if (float64_key != float64(typed_key)) {
				panic(fmt.Sprintf("Invalid operation: %v : %s", operation, "cannot cast the array index to a number"))
			}

			return typed_container[typed_key]

		} else {
			panic(fmt.Sprintf("Invalid field-read operation: %v", operation))
		}

	} else if (strings.Compare(OPT_FIELD_UPDATE, opt) == 0) {
		if (len(operation) != 5) {
			panic(fmt.Sprintf("Invalid field-update operation: %v : %s", operation, "mal-form"))
		}

		new_container := CopyValue(EvalExpr(env, operation[2]))
		key := EvalExpr(env, operation[3])
		new_val := EvalExpr(env, operation[4])

		if (IsType(TYPE_OBJ, new_container) && IsType(TYPE_STRING, key)) { // If the container is an object ...
			typed_container := CastInterfaceToObj(new_container)
			typed_key := CastInterfaceToString(key)
			typed_container[typed_key] = new_val
			return typed_container

		} else if (IsType(TYPE_ARRAY, new_container) && IsType(TYPE_NUMBER, key)) { // If the container is an array ... 
			typed_container := CastInterfaceToArray(new_container)

			float64_key := CastInterfaceToNumber(key)
			typed_key := int(float64_key)
			if (float64_key != float64(typed_key)) {
				panic(fmt.Sprintf("Invalid operation: %v : %s", operation, "cannot cast the array index to a number"))
			}

			typed_container[typed_key] = new_val

			return typed_container

		} else {
			panic(fmt.Sprintf("Invalid field-update operation: %v", operation))
		}

	} else if (strings.Compare(SYMBOL_REF_MARK, opt) == 0) {
		ref_key := []string{}

		for _, key_part := range operation[1:] {
			// DO NOT evaluate key_part --> if key_part itself is not a string --> this is an error! 
			string_key_part := CastInterfaceToString(key_part)
			ref_key = append(ref_key, string_key_part)
		}

		return env.Deref(ref_key)

	} else if (strings.Compare(OPT_FUNC_CALL, opt) == 0) {
		if (len(operation) != 4) {
			panic(fmt.Sprintf("Invalid operation: %v", operation))
		}

		opt_lambda := EvalExpr(env, operation[2])
		lambda, is_lambda := IsLambdaOperation(opt_lambda)

		param := EvalExpr(env, operation[3])
		
		if (!is_lambda) {
			panic(fmt.Sprintf("Invalid lambda given in function call: %v", operation))
		}

		new_env := env.PushStack(param) 
		return EvalExpr(new_env, lambda[2]) 

	} else if (strings.Compare(OPT_ARRAY_MAP, opt) == 0) {
		if (len(operation) != 4) {
			panic(fmt.Sprintf("Invalid array-map operation: %v : %s", operation, "Exact 4 parameters are required"))
		} 

		typed_array := CastInterfaceToArray(EvalExpr(env, operation[2]))
		lambda := EvalExpr(env, operation[3])

		if _, is_lambda := IsLambdaOperation(lambda); !is_lambda {
			panic(fmt.Sprintf("Invalid lambda for array-map operation: %v : %v", operation, lambda))
		}

		arr_result := []interface{}{}
		for _, uneval_arr_element := range typed_array {
			arr_element := EvalExpr(env, uneval_arr_element)
			opt_fcall := []interface{}{SYMBOL_OPT_MARK, OPT_FUNC_CALL, lambda, arr_element}
			eval_result := EvalExpr(env, opt_fcall)
			arr_result = append(arr_result, eval_result) 
		}

		return arr_result

	} else if (strings.Compare(OPT_ARRAY_FILTER, opt) == 0) {
		if (len(operation) != 4) {
			panic(fmt.Sprintf("Invalid array-filter operation: %v : %s", operation, "Exact 4 parameters are required"))
		} 

		typed_array := CastInterfaceToArray(EvalExpr(env, operation[2]))
		lambda := EvalExpr(env, operation[3])

		if _, is_lambda := IsLambdaOperation(lambda); !is_lambda {
			panic(fmt.Sprintf("Invalid lambda for array-filter operation: %v : %v", operation, lambda))
		}

		opt_array_map := []interface{}{SYMBOL_OPT_MARK, OPT_ARRAY_MAP, typed_array, lambda}
		arr_tests := CastInterfaceToArray(EvalExpr(env, opt_array_map))

		if (len(typed_array) != len(arr_tests)) {
			panic(fmt.Sprintf("Unexpected error when evaluating the predicate for array-filter: %v", operation))
		}

		arr_result := []interface{}{}
		for p_id,pred := range arr_tests {
			if (!IsType(TYPE_BOOL, pred)) {
				panic(fmt.Sprintf("Predicate result is not a boolean for element: %v", p_id))
			}
			if (CastInterfaceToBool(pred)) {
				arr_result = append(arr_result, typed_array[p_id])
			}
		}

		return arr_result

	} else {
		panic(fmt.Sprintf("Invalid/Unsupported evaluation of expression: %v", in_expr))
	}
}

/********
State changes 
********/
func Execute (env Environment, opt_list []interface{}) Environment {
	curr_env := env
	for _, opt := range opt_list {
		eval_result := EvalExpr(curr_env, opt)
		new_env, ok := eval_result.(Environment)
		if (ok) { // If the evaluation result is not an Environment, just discard the result 
			curr_env = new_env 
		}
	}
	return curr_env
}