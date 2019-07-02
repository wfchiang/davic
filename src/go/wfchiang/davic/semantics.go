package davic 

import (
	"fmt"
	"strings"
)

/*********
Validation Core 
*********/


/********
Expression evaluation
********/ 
func EvalExpr (env Environment, in_expr interface{}) interface{} {
	// Try to evaluate a reference -- cannot recursively evaluate a reference 
	

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

	} else {
		panic(fmt.Sprintf("Invalid/Unsupported evaluation of expression: %v", in_expr))
	}
}