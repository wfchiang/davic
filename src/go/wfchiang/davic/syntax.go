package davic

import (
	"fmt"
)

const (
	OPT_WEBCALL = "OPT_WEBCALL"
)

/*
Expression 
*/ 
type Expr struct {
	Operator string
	Operands []interface{} 
}

func (expr Expr) Eval() interface{} {
	if (expr.Operator == OPT_WEBCALL) {
		fmt.Println("This is a " + OPT_WEBCALL)
		fmt.Println(len(expr.Operands))
		return 0
	} else {
		panic("Unsupported expr.operator: " + expr.Operator)
	}

	panic("Incompleted expr.eval()")
}

