package types

import "testing"

func TestCompareBoolean(t *testing.T) {
	cases := []struct {
		a, b bool
		want Compared
	}{
		{false, false, ComparedEq},
		{true, true, ComparedEq},
		{false, true, ComparedLt},
		{true, false, ComparedGt},
	}
	for _, c := range cases {
		a := NewBoolean(c.a)
		b := NewBoolean(c.b)
		got := a.Compare(b)
		if got != c.want {
			t.Errorf("(%v).Compare(%v) == %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
