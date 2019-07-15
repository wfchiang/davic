package davic 

import (
	"strings"
)

func ContainsString (string_array []string, target string) bool {
	for _, s := range string_array {
		if (strings.Compare(s, target) == 0) {
			return true 
		}
	}
	return false 
}

func GetKeyString (key []string) string {
	return strings.Join(key, " ")
}

/*
Casting utils  
*/
func CastNumberToInt (value float64) int {
	int_value := int(value)
	if (float64(int_value) != value) {
		panic("Inprecise CastNumberToInt")
	}
	return int_value
}

func CastInterfaceToBool (value interface{}) bool {
	bool_value, ok := value.(bool)
	if (!ok) {
		panic("CastInterfaceToBool failed")
	}
	return bool_value
}

func CastInterfaceToNumber (value interface{}) float64 {
	float64_value, ok := value.(float64)
	if (!ok) {
		panic("CastInterfaceToNumber failed")
	}
	return float64_value
}

func CastInterfaceToInt (value interface{}) int {
	int_value, ok := value.(int)
	if (!ok) {
		panic("CastInterfaceToInt failed")
	}
	return int_value
}

func CastInterfaceToString (value interface{}) string {
	string_value, ok := value.(string)
	if (!ok) {
		panic("CastInterfaceToString failed")
	}
	return string_value
}

func CastInterfaceToArray (value interface{}) []interface{} {
	array_value, ok := value.([]interface{})
	if (!ok) {
		panic("CastInterfaceToArray failed")
	}
	return array_value
}

func CastInterfaceToObj (value interface{}) map[string]interface{} {
	obj_value, ok := value.(map[string]interface{})
	if (!ok) {
		panic("CastInterfaceToObj failed")
	}
	return obj_value
}