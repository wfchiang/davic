package davic

import (
	"fmt"
	"strings"
	"encoding/json"
)

const (
	STACK_DEPTH = 5

	SYMBOL_OPT_MARK = "&"

	SYMBOL_REF_MARK = "#"
	SYMBOL_REF_SEPARATOR = "/"

	TYPE_NULL   = "null"
	TYPE_BOOL   = "bool"
	TYPE_NUMBER = "number"
	TYPE_STRING = "string"
	TYPE_ARRAY  = "array"
	TYPE_OBJ    = "object"

	OPT_RELATION_EQ = "="
	OPT_ARITHMETIC_ADD = "+"
	OPT_WEBCALL = "OPT_WEBCALL"
	OPT_OBJ_FIELD_READ = "OPT_OBJ_FIELD_READ"
)

/*
Type predicates and utils
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
	} else if (strings.Compare(TYPE_ARRAY, type_name) == 0) {
		_, is_type = value.([]interface{})
	} else if (strings.Compare(TYPE_OBJ, type_name) == 0) {
		_, is_type = value.(map[string]interface{})
	} else {
		error_message := fmt.Sprintf("Unknown type (%v) of value %v", type_name, value)
		panic(error_message)
	}

	return is_type
}

func AsNumber (value interface{}) float64 {
	var is_num bool 
	var num float64
	if num, is_num = value.(float64); !is_num {
		panic(fmt.Sprintf("AsNumber: not a number: %v", value))
	}
	return num
}

/*
Copy Value
*/ 
func CopyValue (value interface{}) interface{} {
	if (IsType(TYPE_NULL, value)) {
		return nil
	} else if (IsType(TYPE_BOOL, value) || IsType(TYPE_NUMBER, value) || IsType(TYPE_STRING, value)) {
		return value
	} else if (IsType(TYPE_ARRAY, value)) {
		origin_array, _ := value.([]interface{})
		copy_array := []interface{}{}
		for _, v := range origin_array {
			copy_array = append(copy_array, CopyValue(v))
		}
		return copy_array
	} else if (IsType(TYPE_OBJ, value)) {
		origin_obj, _ := value.(map[string]interface{})
		copy_obj := map[string]interface{}{}
		for k, v := range origin_obj {
			copy_obj[k] = CopyValue(v)
		}
		return copy_obj
	} else {
		error_message := fmt.Sprintf("Cannot copy value with unknown type. Value: %v", value)
		panic(error_message)
	}
}

/*
Json/Object -- as a map[string]interface{}
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
	
	return tokens
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
	StoreStackIndex int
	StoreStack [STACK_DEPTH]interface{} 
}

func CreateNewEnvironment () Environment {
	newEnv := Environment{StoreStackIndex:-1, StoreStack:[STACK_DEPTH]interface{}{}}
	return newEnv
}

func (env Environment) Clone () Environment {
	new_env := CreateNewEnvironment() 
	
	new_env.StoreStackIndex = env.StoreStackIndex
	
	for i, store := range env.StoreStack {
		new_env.StoreStack[i] = CopyValue(store)
	}

	return new_env
}

func (env Environment) GetStore () interface{} {
	if (env.StoreStackIndex < 0 || env.StoreStackIndex >= STACK_DEPTH) {
		panic(fmt.Sprintf("StoreStackIndex out of bound: %v", env.StoreStackIndex))
	}
	return env.StoreStack[env.StoreStackIndex]
}

func (env Environment) PushStore (newStore interface{}) Environment {
	new_env := env.Clone()
	new_env.StoreStackIndex = new_env.StoreStackIndex + 1
	if (new_env.StoreStackIndex < 0 || env.StoreStackIndex >= STACK_DEPTH) {
		panic(fmt.Sprintf("StoreStackIndex out of bound: %v", env.StoreStackIndex))
	}
	new_env.StoreStack[new_env.StoreStackIndex] = CopyValue(newStore)
	return new_env
}

func (env Environment) PopStore () Environment {
	new_env := env.Clone() 
	new_env.StoreStackIndex = new_env.StoreStackIndex - 1
	if (new_env.StoreStackIndex < 0 || env.StoreStackIndex >= STACK_DEPTH) {
		panic(fmt.Sprintf("StoreStackIndex out of bound: %v", env.StoreStackIndex))
	}
	return new_env 
}

func (env Environment) Deref (in_ref interface{}) interface{} {
	var ref []string 

	ref_string, ok := in_ref.(string) 
	if (ok) {
		ref = ParseRefString(ref_string)
	} else {
		ref, ok = in_ref.([]string)
		if (!ok) {
			panic(fmt.Sprintf("Environment deref failed due to invalid reference: %v", in_ref))
		}
	}

	if (!IsRef(ref)) {
		panic(fmt.Sprintf("The given ref is not a reference: %v", ref))
	}	
	return GetObjValue(env.GetStore(), ref[1:])
}