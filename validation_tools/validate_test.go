package validation_tools

import "testing"

func TestName(t *testing.T) {
	res := ValidateEmail("im@gmail.com")
	if !res {
		t.Errorf("Email validation works wrong")
	}
}