package davic

import (
	"fmt"
	"strings"
	"encoding/json"
)

const (
	SYMBOL_OPT_MARK = "&"

	SYMBOL_REF_MARK = "#"
	SYMBOL_REF_SEPARATOR = "/"

	TYPE_NULL   = "TYPE_NULL"
	TYPE_BOOL   = "TYPE_BOOL"
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
	if (strings.Compare(TYPE_NULL, type_name) == 0) {
		is_type = (value == nil)
	} else if (strings.Compare(TYPE_BOOL, type_name) == 0) {
		_, is_type = value.(bool)
	} else if (strings.Compare(TYPE_NUMBER, type_name) == 0) {
		_, is_type = value.(float64)
	} else if (strings.Compare(TYPE_STRING, type_name) == 0) {
		_, is_type = value.(string)
	} else if (strings.Compare(TYPE_OBJ, type_name) == 0) {
		_, is_type = value.(map[string]interface{})
	} else {
		error_message := fmt.Sprintf("Unknown type (%v) of value %v", type_name, value)
		panic(error_message)
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

func GetObjValue (obj interface{}, key []string) interface{} {
	if (len(key) == 0) {
		return obj
	}

	kv, ok := obj.(map[string]interface{})
	if (!ok) {
		panic(fmt.Sprintf("Cannot get obj value from a non-obj value: %v", obj))
	}

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
func IsOperation (in_expr interface{}) ([]interface{}, bool) {
	expr, ok := in_expr.([]interface{})
	if (!ok) {
		return nil, false
	}

	if (len(expr) < 2) {
		return nil, false 
	}

	mark, ok := expr[0].(string)
	if (!ok) {
		return nil, false
	}
	if (strings.Compare(SYMBOL_OPT_MARK, mark) != 0) {
		return nil, false
	}

	_, ok = expr[1].(string)
	if (!ok) {
		return nil, false
	}

	return expr, true
}

func IsBinaryOperation (in_expr interface{}) ([]interface{}, bool) {
	operation, ok := IsOperation(in_expr)
	if (!ok) {
		return nil, false 
	}
	if (len(operation) != 4) {
		return nil, false
	}
	return operation, true
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

/*
Environment definition
*/
type Environment struct {
	Store interface{} 
}

func (env Environment) Deref (ref []string) interface{} {
	if (!IsRef(ref)) {
		panic(fmt.Sprintf("The given ref is not a given reference: %v", ref))
	}	
	return GetObjValue(env.Store, ref[1:])
}