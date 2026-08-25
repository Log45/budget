// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
	got := 2 + 3
	want := 5
	if got != want {
		t.Errorf("Add(2, 3) = %d; want %d", got, want)
	}
}
