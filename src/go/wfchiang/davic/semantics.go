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

/* 
Validation Core 
*/ 
func IsType (type_name string, value interface{}) bool {
	if (strings.Compare(TYPE_BOOL, type_name) == 0) {
		_, is_type := value.(bool)
		if (is_type) {
			return true
		}
	} else {
		panic("Unknown type: " + type_name)
	}

	return false
}