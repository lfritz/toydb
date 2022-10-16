package types

import (
	"reflect"
	"testing"
)

func TestParseDecimalValid(t *testing.T) {
	cases := []struct {
		input string
		want  Decimal
	}{
		{"0", DecimalZero()},
		{"000", DecimalZero()},
		{"123", Decimal{digits: []uint8{1, 2, 3}, n: 3}},
		{"000123", Decimal{digits: []uint8{1, 2, 3}, n: 3}},
		{"123.456", Decimal{digits: []uint8{1, 2, 3, 4, 5, 6}, n: 3}},
		{"123.456000", Decimal{digits: []uint8{1, 2, 3, 4, 5, 6}, n: 3}},
		{"-123", Decimal{negative: true, digits: []uint8{1, 2, 3}, n: 3}},
		{"-123.456", Decimal{negative: true, digits: []uint8{1, 2, 3, 4, 5, 6}, n: 3}},
		{"123.", Decimal{digits: []uint8{1, 2, 3}, n: 3}},
		{".456", Decimal{digits: []uint8{4, 5, 6}, n: 0}},
		{"123.000456", Decimal{digits: []uint8{1, 2, 3, 0, 0, 0, 4, 5, 6}, n: 3}},
		{".000456", Decimal{digits: []uint8{0, 0, 0, 4, 5, 6}, n: 0}},
	}
	for _, c := range cases {
		got, err := ParseDecimal(c.input)
		if err != nil {
			t.Errorf("ParseDecimal(%q) returned error: %v", c.input, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ParseDecimal(%q) returned %#v, want %#v", c.input, got, c.want)
		}
	}
}

func TestParseDecimalInvalid(t *testing.T) {
	cases := []string{
		"",
		"abc",
		"12.34.56",
	}
	for _, input := range cases {
		_, err := ParseDecimal(input)
		if err == nil {
			t.Errorf("ParseDecimal(%q) did not return error", input)
		}
	}
}

func TestDecimalCompare(t *testing.T) {
	cases := []struct {
		a, b string
		want Compared
	}{
		{"0", "0", ComparedEq},
		{"123", "123", ComparedEq},
		{"-123", "-123", ComparedEq},
		{"0.00123", "0.00123", ComparedEq},
		{"123.45", "123.45", ComparedEq},

		{"123", "124", ComparedLt},
		{"-123", "123", ComparedLt},
		{"-124", "-123", ComparedLt},
		{"123.45", "123.46", ComparedLt},
		{"0.123", "0.124", ComparedLt},
		{"0.123", "0.1234", ComparedLt},

		{"124", "123", ComparedGt},
		{"123", "-123", ComparedGt},
		{"-123", "-124", ComparedGt},
		{"123.46", "123.45", ComparedGt},
		{"0.124", "0.123", ComparedGt},
		{"0.1234", "0.123", ComparedGt},
	}

	for _, c := range cases {
		a := NewDecimal(c.a)
		b := NewDecimal(c.b)
		got := a.Compare(b)
		if got != c.want {
			t.Errorf("(%v).Compare(%v) == %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestDecimalString(t *testing.T) {
	cases := []string{
		"0",
		"123",
		"123.45",
		"0.123",
		"0.000123",
		"-123",
		"-123.45",
		"-0.123",
		"-0.000123",
	}
	for _, c := range cases {
		decimal, err := ParseDecimal(c)
		if err != nil {
			t.Fatalf("ParseDecimal(%q) returned error: %v", c, err)
		}
		got := decimal.String()
		if got != c {
			t.Errorf("Decimal %s formatted as %q", c, got)
		}
	}
}

func TestDecimalNormalize(t *testing.T) {
	cases := []struct {
		negative bool
		digits   []uint8
		n        int
		want     Decimal
	}{
		{
			false,
			[]uint8{},
			0,
			Decimal{false, nil, 0},
		},
		{
			false,
			[]uint8{1, 2, 3},
			3,
			Decimal{false, []uint8{1, 2, 3}, 3},
		},
		{
			false,
			[]uint8{1, 2, 3, 4, 5},
			5,
			Decimal{false, []uint8{1, 2, 3, 4, 5}, 5},
		},
		{
			false,
			[]uint8{0, 0, 1, 2, 3},
			5,
			Decimal{false, []uint8{1, 2, 3}, 3},
		},
		{
			false,
			[]uint8{1, 2, 3, 0, 0},
			3,
			Decimal{false, []uint8{1, 2, 3}, 3},
		},
		{
			false,
			[]uint8{1, 2, 3, 4, 0},
			3,
			Decimal{false, []uint8{1, 2, 3, 4}, 3},
		},
		{
			false,
			[]uint8{1, 2, 3, 0, 0, 0, 0},
			5,
			Decimal{false, []uint8{1, 2, 3, 0, 0}, 5},
		},
		{
			false,
			[]uint8{1, 2, 3},
			0,
			Decimal{false, []uint8{1, 2, 3}, 0},
		},
		{
			false,
			[]uint8{0, 0, 1, 2, 3},
			0,
			Decimal{false, []uint8{0, 0, 1, 2, 3}, 0},
		},
	}
	for _, c := range cases {
		got := normalize(c.negative, c.digits, c.n)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("normalize(%v, %v, %v) == %v, want %v", c.negative, c.digits, c.n, got, c.want)
		}
	}
}
