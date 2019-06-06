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
func EvalExpr (in_expr interface{}) interface{} {
	expr, ok := in_expr.([]interface{})
	if (!ok) {
		return in_expr
	}
	
	if (!IsOperation(expr)) {
		return expr
	}

	opt, ok := expr[1].(string)
	if (!ok) {
		panic("Operator (" + fmt.Sprintf("%v", opt) + ")is not a string")
	}

	if (strings.Compare(OPT_RELATION_EQ, opt) == 0) {
		if (!IsBinaryExpr(expr)) {
			panic("Invalid operation: " + OPT_RELATION_EQ + " : " + fmt.Sprintf("%v", expr))
		}
		lhs := expr[2]
		rhs := expr[3]

		if (IsType(TYPE_BOOL, lhs) && IsType(TYPE_BOOL, rhs)) {
			return (lhs == rhs)
		} else if (IsType(TYPE_NUMBER, lhs) && IsType(TYPE_NUMBER, rhs)) {
			return (lhs == rhs)
		} else {
			panic("Unsupport operand type for operation " + OPT_RELATION_EQ)
		}
	} else {
		panic("Invalid operator: " + opt) 
	}
}