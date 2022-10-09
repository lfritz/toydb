package types

import "testing"

func TestInvalidComparisons(t *testing.T) {
	values := []Value{
		NewBoolean(false),
		NewText("hello"),
		NewDecimal("123"),
		NewDate(1999, 12, 31),
	}
	for i, a := range values {
		for j, b := range values {
			got := a.Compare(b)
			gotInvalid := got == ComparedInvalid
			wantInvalid := i != j
			if gotInvalid && !wantInvalid {
				t.Errorf("(%v).Compare(%v) == ComparedInvalid", a, b)
			}
			if !gotInvalid && wantInvalid {
				t.Errorf("(%v).Compare(%v) != ComparedInvalid", a, b)
			}
		}
	}
}
