package davic 

import (
//	"fmt"
	"testing"
)

func TestContainsString (t *testing.T) {
	string_array := []string{"abc", "xyz"}

	if (!ContainsString(string_array, "abc")) {
		t.Error("")
	}

	if (ContainsString(string_array, "cba")) {
		t.Error("")
	}
}