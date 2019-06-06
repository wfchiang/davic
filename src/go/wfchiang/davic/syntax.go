package davic

import (
//	"fmt"
	"strings"
	"encoding/json"
)

const (
	SYMBOL_OPT_MARK = "&"

	SYMBOL_REF_MARK = "#"
	SYMBOL_REF_SEPARATOR = "/"

	TYPE_BOOL   = "TYPE_BOOL"
	TYPE_INT    = "TYPE_INT"
	TYPE_FLOAT  = "TYPE_FLOAT" 
	TYPE_NUMBER = "TYPE_NUMBER"
	TYPE_STRING = "TYPE_STRING"
	TYPE_OBJ    = "TYPE_OBJ"

	OPT_RELATION_EQ = "OPT_RELATION_EQ"
	OPT_WEBCALL = "OPT_WEBCALL"
	OPT_OBJ_FIELD_READ = "OPT_OBJ_FIELD_READ"
)

/*
Type predicates
*/ 
func IsType (type_name string, value interface{}) bool {
	is_type := false
	if (strings.Compare(TYPE_BOOL, type_name) == 0) {
		_, is_type = value.(bool)
	} else if (strings.Compare(TYPE_NUMBER, type_name) == 0) {
		return (IsType(TYPE_INT, value) || IsType(TYPE_FLOAT, value))
	} else if (strings.Compare(TYPE_INT, type_name) == 0) {
		_, is_type = value.(int)
	} else if (strings.Compare(TYPE_FLOAT, type_name) == 0) {
		_, is_type = value.(float64)
	} else if (strings.Compare(TYPE_STRING, type_name) == 0) {
		_, is_type = value.(string)
	} else if (strings.Compare(TYPE_OBJ, type_name) == 0) {
		_, is_type = value.(map[string]interface{})
	} else {
		panic("Unknown type: " + type_name)
	}

	return is_type
}

/*
Json -- as a map[string]interface{}
*/
func IsRefString (ref_string string) bool { 
	return strings.HasPrefix(ref_string, SYMBOL_REF_MARK)
}

func IsRef (ref []string) bool {
	if (len(ref) == 0) {
		return false 
	}
	return (strings.Compare(SYMBOL_REF_MARK, ref[0]) == 0)
}

func ParseRefString (ref_string string) []string {
	tokens := strings.Split(ref_string, SYMBOL_REF_SEPARATOR) 

	if (!IsRef(tokens)) {
		panic("Invalid reference string: " + ref_string) 
	}
	
	return tokens[1:len(tokens)]
}

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
func IsOperation (expr []interface{}) bool {
	if (len(expr) < 2) {
		return false 
	}

	mark, ok := expr[0].(string)
	if (!ok) {
		return false
	}
	if (strings.Compare(SYMBOL_OPT_MARK, mark) != 0) {
		return false
	}

	_, ok = expr[1].(string)
	if (!ok) {
		return false
	}

	return true
}

func IsBinaryExpr (expr []interface{}) bool {
	if (!IsOperation(expr)) {
		return false
	}
	return (len(expr) == 4)
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

