package davic 

import (
	"strings"
)

const (
	TYPE_BOOL   = "TYPE_BOOL"
	TYPE_INT    = "TYPE_INT"
	TYPE_FLOAT  = "TYPE_FLOAT" 
	TYPE_STRING = "TYPE_STRING"
	TYPE_OBJ    = "TYPE_OBJ" 
)

/*********
Validation Core 
*********/
/*
Type predicates
*/ 
func IsType (type_name string, value interface{}) bool {
	is_type := false
	if (strings.Compare(TYPE_BOOL, type_name) == 0) {
		_, is_type = value.(bool)
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
Validation function 
*/ 
func ValidateType (key []string, type_name string, value interface{}) ValidationResult {
	vResult := false
	vComments := []string{}

	if (ContainsString([]string{TYPE_BOOL, TYPE_FLOAT, TYPE_STRING, TYPE_OBJ}, type_name)) {
		vResult = IsType(type_name, value)
		if (!vResult) {
			vComments = append(vComments, MakeValidationComment(key, "Invalid type " + type_name))
		}
	} else {
		vResult = false 
		vComments = append(vComments, MakeValidationComment(key, "Unknown type " + type_name))
	}

	return ValidationResult{IsValid:vResult, Comments:vComments}
}