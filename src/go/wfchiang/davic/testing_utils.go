package davic 

import (
	"testing"
)

func simpleIsViolation (type_string string, expected interface{}, actual interface{}) bool {
	if (!IsType(type_string, expected)) {
		return false
	}
	if (!IsType(type_string, actual)) {
		return true
	}
	return (expected != actual); 
}

func simpleExpectPanic (t *testing.T) {
	r := recover() 
	if (r == nil) {
		t.Error("No expected panic occurred...")
	}
}

func simpleRecover (t *testing.T) {
	if r := recover(); r != nil {
		t.Error("There was a panic... ", r)
	}
}