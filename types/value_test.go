package types

import "testing"

func TestValueCompare(t *testing.T) {
	cases := []struct {
		a, b Value
		want Compared
	}{
		{NewValue(NewDecimal("4")), NewValue(NewDecimal("5")), ComparedLt},

		{NewNull(TypeDecimal), NewValue(NewDecimal("5")), ComparedNull},
		{NewValue(NewDecimal("5")), NewNull(TypeDecimal), ComparedNull},
		{NewNull(TypeDecimal), NewNull(TypeDecimal), ComparedNull},

		{NewValue(NewDecimal("4")), NewValue(NewText("hello")), ComparedInvalid},
		{NewValue(NewDecimal("4")), NewNull(TypeText), ComparedInvalid},
		{NewNull(TypeDecimal), NewValue(NewText("hello")), ComparedInvalid},
		{NewNull(TypeDecimal), NewNull(TypeText), ComparedInvalid},
	}
	for _, c := range cases {
		got := c.a.Compare(c.b)
		if got != c.want {
			t.Errorf("(%v).Compare(%v) == %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestInvalidComparisons(t *testing.T) {
	values := []BasicValue{
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
