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

	if (strings.Compare(OPT_RELATION_EQ, opt) == 0) {
		operation, ok = IsBinaryOperation(in_expr)
		if (!ok) {
			panic("Invalid operation: " + OPT_RELATION_EQ + " : " + fmt.Sprintf("%v", operation))
		}
		lhs := EvalExpr(env, operation[2])
		rhs := EvalExpr(env, operation[3])

		if (IsType(TYPE_BOOL, lhs) && IsType(TYPE_BOOL, rhs)) {
			return (lhs == rhs)
		} else if (IsType(TYPE_NUMBER, lhs) && IsType(TYPE_NUMBER, rhs)) {
			return (lhs == rhs)
		} else {
			panic("Unsupport operand type for operation " + OPT_RELATION_EQ)
		}

	} else if (strings.Compare(OPT_ARITHMETIC_ADD, opt) == 0) {
		if (len(operation) < 3) {
			panic("Invalid operation: " + OPT_ARITHMETIC_ADD + " : " + fmt.Sprintf("%v", operation))
		}
		
		add_result := AsNumber(EvalExpr(env, operation[2]))
		for _, v := range operation[3:] {
			add_result = add_result + AsNumber(EvalExpr(env, v))
		}
		
		return add_result

	} else {
		panic(fmt.Sprintf("Invalid/Unsupported evaluation of expression: %v", in_expr))
	}
}