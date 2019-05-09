package davic 

import (
	"testing"
)

func simpleRecover (t *testing.T) {
	if r := recover(); r != nil {
		t.Error("There was a panic... ", r)
	}
}

func TestIsType (t *testing.T) {

	defer simpleRecover(t) 

	if is_type := IsType(TYPE_BOOL, false); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_BOOL, 1); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_INT, false); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_INT, 1); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_FLOAT, 1); is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_FLOAT, 1.1); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_STRING, 1); is_type {
		t.Error("")
	}
	
	if is_type := IsType(TYPE_STRING, "hello"); !is_type {
		t.Error("")
	}

	if is_type := IsType(TYPE_OBJ, 1); is_type {
		t.Error("")
	}
}