package davic

import (
	"fmt"
	"strings"
	"encoding/json"
	"container/list"
)

const (
	STACK_DEPTH = 5

	SYMBOL_OPT_MARK = "&"
	SYMBOL_REF_MARK = "#"

	SYMBOL_HTTP_METHOD_GET  = "GET"
	SYMBOL_HTTP_METHOD_POST = "POST"

	TYPE_NULL   = "null"
	TYPE_BOOL   = "bool"
	TYPE_NUMBER = "number"
	TYPE_STRING = "string"
	TYPE_ARRAY  = "array"
	TYPE_OBJ    = "object"

	OPT_STACK_READ = "~"
	OPT_LAMBDA = "^"
	OPT_FUNC_CALL = "!"

	OPT_RELATION_EQ = "="

	OPT_ARITHMETIC_ADD = "+"

	OPT_ARRAY_MAP    = "-a-map-"
	OPT_ARRAY_FILTER = "-a-fitler-"
	
	OPT_HTTP_CALL = "-http-call-"
	
	OPT_FIELD_READ = "-fr-"
	OPT_FIELD_UPDATE = "-fu-"
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
func IsRef (ref []string) bool {
	if (len(ref) == 0) {
		return false 
	}
	return (strings.Compare(SYMBOL_REF_MARK, ref[0]) == 0)
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

func IsUnaryOperation (in_expr interface{}) ([]interface{}, bool) {
	operation, ok := IsOperation(in_expr)
	if (!ok || len(operation) != 3) {
		return nil, false
	}
	return operation, true
}

func IsBinaryOperation (in_expr interface{}) ([]interface{}, bool) {
	operation, ok := IsOperation(in_expr)
	if (!ok || len(operation) != 4) {
		return nil, false 
	}
	return operation, true
}

func IsLambdaOperation (in_expr interface{}) ([]interface{}, bool) {
	operation, ok := IsUnaryOperation(in_expr)
	if (!ok || operation[1] != OPT_LAMBDA) {
		return nil, false
	}
	return operation, true 
}

func IsHttpOperation (in_expr interface{}) ([]interface{}, bool) {
	operation, ok := IsOperation(in_expr)
	if (!ok || len(operation) != 5) {
		return nil, false
	}

	// check method 
	if (!IsType(TYPE_STRING, operation[2])) {
		return nil, false
	}
	str_http_method := CastInterfaceToString(operation[2])
	if (!ContainsString([]string{SYMBOL_HTTP_METHOD_GET,SYMBOL_HTTP_METHOD_POST}, str_http_method)) {
		return nil, false
	}

	// check headers 
	if (!IsType(TYPE_OBJ, operation[3])) {
		return nil, false 
	}

	// check body 
	if (
		!IsType(TYPE_NULL, operation[4]) && 
		!IsType(TYPE_BOOL, operation[4]) && 
		!IsType(TYPE_NUMBER, operation[4]) && 
		!IsType(TYPE_STRING, operation[4]) && 
		!IsType(TYPE_ARRAY, operation[4]) && 
		!IsType(TYPE_OBJ, operation[4])) {
		return nil, false
	}
			
	return operation, true
}

/*
Environment definition
*/
type Environment struct {
	Store interface{}
	Stack *list.List
}

func CreateNewEnvironment () Environment {
	newEnv := Environment{Store:nil,Stack:list.New()}
	return newEnv
}

func (env Environment) Clone () Environment {
	new_env := CreateNewEnvironment() 
	
	new_env.Store = CopyValue(env.Store)

	new_env.Stack = list.New()
	stack_element := env.Stack.Front()
	for i := 0 ; i < STACK_DEPTH ; i++ {
		if (stack_element == nil) {
			break
		}
		new_env.Stack.PushBack(CopyValue(stack_element.Value))
		stack_element = stack_element.Next()
	}

	return new_env
}

func (env Environment) Deref (in_ref interface{}) interface{} {
	ref, ok := in_ref.([]string)
	if (!ok) {
		panic(fmt.Sprintf("Environment deref failed due to invalid reference: %v", in_ref))
	}
	
	if (!IsRef(ref)) {
		panic(fmt.Sprintf("The given ref is not a reference: %v", ref))
	}	
	return GetObjValue(env.Store, ref[1:])
}

func (env Environment) PushStack (stack_value interface{}) Environment {
	if (env.Stack.Len() >= STACK_DEPTH) {
		panic("Cannot push stack: stack overflow")
	}
	new_env := env.Clone() 
	new_env.Stack.PushBack(stack_value)
	return new_env 
}

func (env Environment) PopStack () Environment {
	if (env.Stack.Len() == 0) {
		panic("Cannot pop stack: stack depth = 0")
	}
	new_env := env.Clone() 
	new_env.Stack.Remove(new_env.Stack.Back())
	return new_env
}

func (env Environment) ReadStack () interface{} {
	if (env.Stack.Len() == 0) {
		panic("Cannot read from an empty stack")
	}
	if (env.Stack.Len() >= STACK_DEPTH) {
		panic("Stack overflow found when read...")
	}
	return env.Stack.Back().Value
}