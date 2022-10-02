package sql

import "testing"

func TestParseDecimalValid(t *testing.T) {
	cases := []struct {
		input string
		want  Decimal
	}{
		{"0", Decimal{Value: 0}},
		{"000", Decimal{Value: 0}},
		{"123", Decimal{Value: 123}},
		{"000123", Decimal{Value: 123}},
		{"123.456", Decimal{Value: 123456, Digits: 3}},
		{"123.456000", Decimal{Value: 123456, Digits: 3}},
		{"-123", Decimal{Value: -123}},
		{"-123.456", Decimal{Value: -123456, Digits: 3}},
		{"123.", Decimal{Value: 123}},
		{".456", Decimal{Value: 456, Digits: 3}},
		{"123.000456", Decimal{Value: 123000456, Digits: 6}},
		{".000456", Decimal{Value: 456, Digits: 6}},
	}
	for _, c := range cases {
		got, err := ParseDecimal(c.input)
		if err != nil {
			t.Errorf("ParseDecimal(%q) returned error: %v", c.input, err)
			continue
		}
		if got != c.want {
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
