package sql

import (
	"reflect"
	"testing"
)

func checkParser[T any](t *testing.T, name string, parse Parser[T], input string, want T) {
	ts, err := Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize(%q) returned error: %v", input, err)
	}

	tokens := &TokenList{input, ts}
	got, remaining, err := parse(tokens)
	if err != nil {
		t.Fatalf("%s returned error: %v", name, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s for %q returned\n%v, want\n%v", name, input, got, want)
	}
	if remaining.Len() != 0 {
		t.Errorf("%s for %q did not consume all tokens", name, input)
	}
}

func checkParserInvalid[T any](t *testing.T, name string, parse Parser[T], input string) {
	ts, err := Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize(%q) returned error: %v", input, err)
	}

	tokens := &TokenList{input, ts}
	_, _, err = parse(tokens)
	if err == nil {
		t.Errorf("%s did not return error for %q", name, input)
	}
}

func TestParseSelectList(t *testing.T) {
	cases := []struct {
		input string
		want  SelectList
	}{
		{"*", Star{}},
		{"", ExpressionList{}},
		{"foo", ExpressionList{Expressions: []Expression{ColumnReference{Name: "foo"}}}},
		{
			"'a', 'b', 'c'",
			ExpressionList{
				Expressions: []Expression{String{Value: "a"}, String{Value: "b"}, String{Value: "c"}},
			},
		},
	}
	for _, c := range cases {
		checkParser(t, "ParseSelectList", ParseSelectList, c.input, c.want)
	}

	invalid := []string{
		"foo,",
	}
	for _, input := range invalid {
		checkParserInvalid(t, "ParseSelectList", ParseSelectList, input)
	}
}

func TestParseExpression(t *testing.T) {
	cases := []struct {
		input string
		want  Expression
	}{
		{
			"'hello'",
			String{Value: "hello"},
		},
		{
			"'hello' = 'ciao'",
			&BinaryOperation{
				Left:     String{Value: "hello"},
				Operator: BinaryOperatorEq,
				Right:    String{Value: "ciao"},
			},
		},
		{
			"12.3 < 45.6",
			&BinaryOperation{
				Left:     Number{Decimal{Value: 123, Digits: 1}},
				Operator: BinaryOperatorLt,
				Right:    Number{Decimal{Value: 456, Digits: 1}},
			},
		},
	}
	for _, c := range cases {
		checkParser(t, "ParseExpression", ParseExpression, c.input, c.want)
	}

	invalid := []string{
		"",
		"'hello' = ",
		" = 'hello'",
	}
	for _, input := range invalid {
		checkParserInvalid(t, "ParseExpression", ParseExpression, input)
	}
}

func TestParseValue(t *testing.T) {
	cases := []struct {
		input string
		want  Expression
	}{
		{"'hello'", String{Value: "hello"}},
		{"12.34", Number{Value: Decimal{Value: 1234, Digits: 2}}},
		{"foo", ColumnReference{Name: "foo"}},
		{"foo.bar", ColumnReference{Relation: "foo", Name: "bar"}},
	}
	for _, c := range cases {
		checkParser(t, "ParseValue", ParseValue, c.input, c.want)
	}

	invalid := []string{
		"",
		",",
	}
	for _, input := range invalid {
		checkParserInvalid(t, "ParseValue", ParseValue, input)
	}
}
