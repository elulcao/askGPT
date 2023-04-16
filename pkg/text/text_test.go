package text

import (
	"reflect"
	"testing"
)

func TestWrapText(t *testing.T) {
	input := "This is a long string that needs to be wrapped."
	expected := "This is a long\nstring that needs\nto be wrapped."
	result := wrapText(input, 15)

	// Compare the result to the expected output slices
	if reflect.DeepEqual(result, []string{expected}) {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
