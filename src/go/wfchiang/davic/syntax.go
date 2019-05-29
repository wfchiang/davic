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
Json -- as a map[string]interface{}
*/
func GetObjKeys (obj map[string]interface{}) []string {
	var keys []string
	for k, _ := range obj {
		keys = append(keys, k)
	}
	return keys
}

func GetObjValue (obj map[string]interface{}, key []string) interface{} {
	var kv map[string]interface{} = obj
	var retv interface{} = nil
	
	for i, k := range key {
		if (i == (len(key)-1)) {
			retv = kv[k]
			break
		}
		kv = kv[k].(map[string]interface{})
	}

	return retv
}

func CreateObjFromBytes (byte_array []byte) map[string]interface{} {
	var new_jnode map[string]interface{}
	json.Unmarshal(byte_array, &new_jnode)
	return new_jnode
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

/*
Validation Result
*/ 
type ValidationResult struct {
	IsValid bool
	Comments []string 
}

func MergeValidationResults (vr0 ValidationResult, vr1 ValidationResult) ValidationResult {
	var final_result ValidationResult
	final_result.IsValid = (vr0.IsValid && vr1.IsValid)
	final_result.Comments = append(vr0.Comments, vr1.Comments...)
	return final_result
}

func MakeValidationComment (key []string, comment string) string {
	return "On field " + GetKeyString(key) + " : " + comment 
}

