package sql

import (
	"reflect"
	"testing"
)

func TestParseExpressionValid(t *testing.T) {
	cases := []struct {
		input string
		want  Expression
	}{
		{"'hello'", String{Value: "hello"}},
	}

	for _, c := range cases {
		ts, err := Tokenize(c.input)
		if err != nil {
			t.Fatalf("Tokenize(%q) returned error: %v", c.input, err)
		}

		tokens := &TokenList{c.input, ts}
		got, remaining, err := ParseExpression(tokens)
		if err != nil {
			t.Fatalf("ParseExpression returned error: %v", err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ParseExpression for %q returned\n%s, want\n%s",
				c.input, got.PrintExpression(), c.want.PrintExpression())
		}
		if remaining.Len() != 0 {
			t.Errorf("ParseExpression for %q did not consume all tokens", c.input)
		}
	}
}
