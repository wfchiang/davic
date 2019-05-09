package davic 

import (
	"testing"
)

func TestIsTypeBool0 (t *testing.T) {
	is_type := IsType(TYPE_BOOL, false)
	if (!is_type) {
		t.Error("IsType fails to check bool")
	}

	is_type = IsType(TYPE_BOOL, 1) 
	if (is_type) {
		t.Error("IsType fails to check bool")
	}
}