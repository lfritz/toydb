package types

import "testing"

func TestTextCompare(t *testing.T) {
	cases := []struct {
		a, b string
		want Compared
	}{
		{"", "", ComparedEq},
		{"a", "a", ComparedEq},

		{"a", "b", ComparedLt},
		{"a", "aa", ComparedLt},

		{"b", "a", ComparedGt},
		{"aa", "a", ComparedGt},
	}
	for _, c := range cases {
		a := NewText(c.a)
		b := NewText(c.b)
		got := a.Compare(b)
		if got != c.want {
			t.Errorf("(%v).Compare(%v) == %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
