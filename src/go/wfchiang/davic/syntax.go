package davic

import (
	"fmt"
	"strconv"
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

	OPT_STORE_READ = "-store-read-"
	OPT_STORE_WRITE = "-store-write-"
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

	KEY_HTTP_STATUS  = "status" 
	KEY_HTTP_METHOD  = "method"
	KEY_HTTP_URL     = "url"
	KEY_HTTP_HEADERS = "headers"
	KEY_HTTP_BODY    = "body"
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
	operation, ok := IsUnaryOperation(in_expr)
	if (!ok || operation[1] != OPT_HTTP_CALL) {
		return nil, false
	}
	if (!IsType(TYPE_OBJ, operation[2])) {
		return nil, false 
	}
	http_request := CastInterfaceToObj(operation[2])
	
	// check method 
	http_method, ok := http_request[KEY_HTTP_METHOD]
	if (!ok || !IsType(TYPE_STRING, http_method)) {
		return nil, false
	}
	if (!ContainsString([]string{SYMBOL_HTTP_METHOD_GET,SYMBOL_HTTP_METHOD_POST}, CastInterfaceToString(http_method))) {
		return nil, false
	}

	// check headers 
	http_headers, ok := http_request[KEY_HTTP_HEADERS]
	if (!ok || !IsType(TYPE_OBJ, http_headers)) {
		return nil, false
	}

	// check body 
	http_body, ok := http_request[KEY_HTTP_BODY]
	if (!ok) {
		return nil, false
	}
	if (
		!IsType(TYPE_NULL, http_body) && 
		!IsType(TYPE_BOOL, http_body) && 
		!IsType(TYPE_NUMBER, http_body) && 
		!IsType(TYPE_STRING, http_body) && 
		!IsType(TYPE_ARRAY, http_body) && 
		!IsType(TYPE_OBJ, http_body)) {
		return nil, false
	}
			
	return operation, true
}

func IsHttpHeaders (in_expr interface{}) (map[string]interface{}, bool) {
	if (!IsType(TYPE_OBJ, in_expr)) {
		return nil, false 
	}	

	return CastInterfaceToObj(in_expr), true
} 

func IsHttpBody (in_expr interface{}) (interface{}, bool) {
	if (
		IsType(TYPE_NULL, in_expr) || 
		IsType(TYPE_BOOL, in_expr) || 
		IsType(TYPE_NUMBER, in_expr) || 
		IsType(TYPE_STRING, in_expr) || 
		IsType(TYPE_ARRAY, in_expr) || 
		IsType(TYPE_OBJ, in_expr)) {
		return in_expr, true
	}
	return nil, false 
}

func IsHttpRequest (in_expr interface{}) (map[string]interface{}, bool) {
	if (!IsType(TYPE_OBJ, in_expr)) {
		return nil, false 
	}
	http_request := CastInterfaceToObj(in_expr)

	// check method 
	http_method, ok := http_request[KEY_HTTP_METHOD]
	if (!ok || !IsType(TYPE_STRING, http_method)) {
		return nil, false
	}
	if (!ContainsString([]string{SYMBOL_HTTP_METHOD_GET,SYMBOL_HTTP_METHOD_POST}, CastInterfaceToString(http_method))) {
		return nil, false
	}

	// check URL 
	val_url, ok := http_request[KEY_HTTP_METHOD]
	if (!ok) {
		return nil, false 
	}
	if (!IsType(TYPE_STRING, val_url)) {
		return nil, false 
	}

	// check headers 
	val_headers, ok := http_request[KEY_HTTP_HEADERS]
	if (!ok) {
		return nil, false 
	}
	_, ok = IsHttpHeaders(val_headers)
	if (!ok) {
		return nil, false
	}

	// check body 
	val_body, ok := http_request[KEY_HTTP_BODY]
	if (!ok) {
		return nil, false 
	}
	_, ok = IsHttpBody(val_body)
	if (!ok) {
		return nil, false
	}

	return http_request, true
}


func IsHttpResponse (in_expr interface{}) (map[string]interface{}, bool) {
	if (!IsType(TYPE_OBJ, in_expr)) {
		return nil, false 
	}
	http_res := CastInterfaceToObj(in_expr)

	// check http status 
	http_status, ok := http_res[KEY_HTTP_STATUS]
	if (!ok || !IsType(TYPE_STRING, http_status)) {
		return nil, false
	}
	_, err := strconv.Atoi(CastInterfaceToString(http_status))
	if (err != nil) {
		return nil, false 
	}

	// check http headers 
	http_headers, ok := http_res[KEY_HTTP_HEADERS]
	if (!ok) {
		return nil, false 
	}
	if _, ok := IsHttpHeaders(http_headers); !ok {
		return nil, false 
	}

	// check http body 
	http_body, ok := http_res[KEY_HTTP_BODY]
	if (!ok) {
		return nil, false 
	}
	if _, ok := IsHttpBody(http_body); !ok {
		return nil, false 
	}

	return http_res, true
}

/*
Environment definition
*/
type Environment struct {
	Store interface{}
	Stack *list.List
}

func CreateNewEnvironment () Environment {
	newEnv := Environment{Store:map[string]interface{}{},Stack:list.New()}
	return newEnv
}

func (env Environment) Clone () Environment {
	new_env := CreateNewEnvironment() 
	
	new_env.Store = CopyValue(env.Store)

	if (env.Stack == nil) {
		new_env.Stack = nil
	} else {
		new_env.Stack = list.New()
		stack_element := env.Stack.Front()
		for i := 0 ; i < STACK_DEPTH ; i++ {
			if (stack_element == nil) {
				break
			}
			new_env.Stack.PushBack(CopyValue(stack_element.Value))
			stack_element = stack_element.Next()
		}
	}

	return new_env
}

func (env Environment) WriteStore (vkey string, val interface{}) Environment {
	if (!IsType(TYPE_OBJ, env.Store)) {
		panic("Environment is corruptted -- the Store is not an object")
	}
	new_env := env.Clone() 
	if (!IsType(TYPE_OBJ, new_env.Store)) {
		panic("Environment (the clone) is corruptted -- the Store is not an object")
	}
	new_env_store := CastInterfaceToObj(new_env.Store) 
	new_env_store[vkey] = val
	return new_env
}

func (env Environment) ReadStore (vkey string) interface{} {
	if (!IsType(TYPE_OBJ, env.Store)) {
		panic("Environment is corruptted -- the Store is not an object")
	}
	return GetObjValue(env.Store, []string{vkey})
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