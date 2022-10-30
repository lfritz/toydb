package sql

import (
	"reflect"
	"testing"

	"github.com/lfritz/toydb/types"
)

func TestParse(t *testing.T) {
	_, err := Parse("select foo from bar;")
	if err != nil {
		t.Errorf("Parse returned error for valid statement: %v", err)
	}

	_, err = Parse("select foo from bar; hello")
	if err == nil {
		t.Error("Parse did not return error for statement with extra text at the end")
	}

}

func checkParser[T any](t *testing.T, name string, parse Parser[T], input string, want T) {
	t.Helper()

	ts, err := Tokenize(input)
	if err != nil {
		t.Errorf("Tokenize(%q) returned error: %v", input, err)
		return
	}

	tokens := &TokenList{input, ts}
	got, remaining, err := parse(tokens)
	if err != nil {
		t.Errorf("%s returned error for %q: %v", name, input, err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s for %q returned\n%v, want\n%v", name, input, got, want)
	}
	if remaining.Len() != 0 {
		t.Errorf("%s for %q did not consume all tokens", name, input)
	}
}

func checkParserInvalid[T any](t *testing.T, name string, parse Parser[T], input string) {
	t.Helper()

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

var condition1 = &BinaryOperation{
	Left:     ColumnReference{Relation: "foo", Name: "x"},
	Operator: BinaryOperatorEq,
	Right:    ColumnReference{Relation: "bar", Name: "x"},
}

func TestParseSelectStatement(t *testing.T) {
	cases := []struct {
		input string
		want  *SelectStatement
	}{
		{
			"select x, y from foo",
			&SelectStatement{
				What: ExpressionList{
					[]Expression{ColumnReference{Name: "x"}, ColumnReference{Name: "y"}},
				},
				From: TableName{Name: "foo"},
			},
		},
		{
			"select * from foo join bar on foo.x = bar.x",
			&SelectStatement{
				What: Star{},
				From: &Join{
					Type:      JoinTypeInner,
					Left:      TableName{"foo"},
					Right:     TableName{"bar"},
					Condition: condition1,
				},
			},
		},
		{
			"select * from foo where foo.x = 0",
			&SelectStatement{
				What: Star{},
				From: TableName{Name: "foo"},
				Where: &BinaryOperation{
					Left:     ColumnReference{Relation: "foo", Name: "x"},
					Operator: BinaryOperatorEq,
					Right:    Number{Value: types.DecimalZero()},
				},
			},
		},
	}
	for _, c := range cases {
		checkParser(t, "ParseSelectStatement", ParseSelectStatement, c.input, c.want)
	}

	invalid := []string{
		"",
		"select x, * from foo",
		"select x, y from",
	}
	for _, input := range invalid {
		checkParserInvalid(t, "ParseSelectStatement", ParseSelectStatement, input)
	}
}

func TestParseTableReference(t *testing.T) {
	cases := []struct {
		input string
		want  TableReference
	}{
		{"foo", TableName{"foo"}},
		{
			"foo join bar on foo.x = bar.x",
			&Join{
				Type:      JoinTypeInner,
				Left:      TableName{"foo"},
				Right:     TableName{"bar"},
				Condition: condition1,
			},
		},
		{
			"foo left outer join bar on foo.x = bar.x",
			&Join{
				Type:      JoinTypeLeftOuter,
				Left:      TableName{"foo"},
				Right:     TableName{"bar"},
				Condition: condition1,
			},
		},
	}
	for _, c := range cases {
		checkParser(t, "ParseTableReference", ParseTableReference, c.input, c.want)
	}

	invalid := []string{
		"",
		"123",
		"foo join bar",
		"foo join bar on",
		"foo join on foo.x = bar.x",
	}
	for _, input := range invalid {
		checkParserInvalid(t, "ParseTableReference", ParseTableReference, input)
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
				Left:     Number{types.NewDecimal("12.3")},
				Operator: BinaryOperatorLt,
				Right:    Number{types.NewDecimal("45.6")},
			},
		},
		{
			"foo is null",
			&UnaryOperation{
				Operand:  ColumnReference{Name: "foo"},
				Operator: UnaryOperatorIsNull,
			},
		},
		{
			"foo is not null",
			&UnaryOperation{
				Operand:  ColumnReference{Name: "foo"},
				Operator: UnaryOperatorIsNotNull,
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
		"4 = is null",
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
		{"12.34", Number{Value: types.NewDecimal("12.34")}},
		{"true", Boolean{Value: true}},
		{"foo", ColumnReference{Name: "foo"}},
		{"foo.bar", ColumnReference{Relation: "foo", Name: "bar"}},
		{"date '1999-12-31'", Date{Value: types.NewDate(1999, 12, 31)}},
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
