package davic

import (
	"fmt"
	"encoding/json"
)

const (
	OPT_WEBCALL = "OPT_WEBCALL"
	OPT_OBJ_FIELD_READ = "OPT_OBJ_FIELD_READ"
)

/*
JNode
*/
type JNode struct {
	key_value map[string]interface{}
}

func (jnode *JNode) InterpretBytes (in_bytes []byte) {
	json.Unmarshal(in_bytes, &jnode.key_value)
}

func (jnode *JNode) GetKeys () []string {
	var keys []string
	for k, _ := range jnode.key_value {
		keys = append(keys, k)
	}
	return keys
}

func (jnode *JNode) GetValue (key []string) interface{} {
	var kv map[string]interface{} = jnode.key_value
	var retv interface{} = jnode.key_value

	for i, k := range key {
		retv = kv[k]
		if (i == (len(key)-1)) {
			break
		}
		kv = retv.(map[string]interface{})
	}

	return retv
}


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

