package sql

import (
	"reflect"
	"testing"
)

func TestLexerValid(t *testing.T) {
	cases := []struct {
		input string
		want  []Token
	}{
		{
			"",
			nil,
		},
		{
			"select foo, bar, baz from qux;",
			[]Token{
				Token{TokenTypeSelect, "select"},
				Token{TokenTypeIdentifier, "foo"},
				Token{TokenTypeComma, ","},
				Token{TokenTypeIdentifier, "bar"},
				Token{TokenTypeComma, ","},
				Token{TokenTypeIdentifier, "baz"},
				Token{TokenTypeFrom, "from"},
				Token{TokenTypeIdentifier, "qux"},
				Token{TokenTypeSemicolon, ";"},
			},
		},
		{
			"select\tfoo  ,bar\t,baz from \t \n qux;   ",
			[]Token{
				Token{TokenTypeSelect, "select"},
				Token{TokenTypeIdentifier, "foo"},
				Token{TokenTypeComma, ","},
				Token{TokenTypeIdentifier, "bar"},
				Token{TokenTypeComma, ","},
				Token{TokenTypeIdentifier, "baz"},
				Token{TokenTypeFrom, "from"},
				Token{TokenTypeIdentifier, "qux"},
				Token{TokenTypeSemicolon, ";"},
			},
		},
		{
			"SELECT _foo, bar123, a_b$c FROM qux",
			[]Token{
				Token{TokenTypeSelect, "SELECT"},
				Token{TokenTypeIdentifier, "_foo"},
				Token{TokenTypeComma, ","},
				Token{TokenTypeIdentifier, "bar123"},
				Token{TokenTypeComma, ","},
				Token{TokenTypeIdentifier, "a_b$c"},
				Token{TokenTypeFrom, "FROM"},
				Token{TokenTypeIdentifier, "qux"},
			},
		},
		{
			"select foo from bar where (x = 123.45 or y < 0) and z >= .4",
			[]Token{
				Token{TokenTypeSelect, "select"},
				Token{TokenTypeIdentifier, "foo"},
				Token{TokenTypeFrom, "from"},
				Token{TokenTypeIdentifier, "bar"},
				Token{TokenTypeWhere, "where"},
				Token{TokenTypeOpenParen, "("},
				Token{TokenTypeIdentifier, "x"},
				Token{TokenTypeEq, "="},
				Token{TokenTypeNumber, "123.45"},
				Token{TokenTypeOr, "or"},
				Token{TokenTypeIdentifier, "y"},
				Token{TokenTypeLt, "<"},
				Token{TokenTypeNumber, "0"},
				Token{TokenTypeCloseParen, ")"},
				Token{TokenTypeAnd, "and"},
				Token{TokenTypeIdentifier, "z"},
				Token{TokenTypeGe, ">="},
				Token{TokenTypeNumber, ".4"},
			},
		},
		{
			"select * from foo where x is not null",
			[]Token{
				Token{TokenTypeSelect, "select"},
				Token{TokenTypeStar, "*"},
				Token{TokenTypeFrom, "from"},
				Token{TokenTypeIdentifier, "foo"},
				Token{TokenTypeWhere, "where"},
				Token{TokenTypeIdentifier, "x"},
				Token{TokenTypeIs, "is"},
				Token{TokenTypeNot, "not"},
				Token{TokenTypeNull, "null"},
			},
		},
		{
			"select foo from bar where x != 'hello' or y <> 'ciao'",
			[]Token{
				Token{TokenTypeSelect, "select"},
				Token{TokenTypeIdentifier, "foo"},
				Token{TokenTypeFrom, "from"},
				Token{TokenTypeIdentifier, "bar"},
				Token{TokenTypeWhere, "where"},
				Token{TokenTypeIdentifier, "x"},
				Token{TokenTypeNe, "!="},
				Token{TokenTypeString, "hello"},
				Token{TokenTypeOr, "or"},
				Token{TokenTypeIdentifier, "y"},
				Token{TokenTypeNe, "<>"},
				Token{TokenTypeString, "ciao"},
			},
		},
		{
			"select foo.x, bar.y from foo left outer join bar",
			[]Token{
				Token{TokenTypeSelect, "select"},
				Token{TokenTypeIdentifier, "foo"},
				Token{TokenTypeDot, "."},
				Token{TokenTypeIdentifier, "x"},
				Token{TokenTypeComma, ","},
				Token{TokenTypeIdentifier, "bar"},
				Token{TokenTypeDot, "."},
				Token{TokenTypeIdentifier, "y"},
				Token{TokenTypeFrom, "from"},
				Token{TokenTypeIdentifier, "foo"},
				Token{TokenTypeLeft, "left"},
				Token{TokenTypeOuter, "outer"},
				Token{TokenTypeJoin, "join"},
				Token{TokenTypeIdentifier, "bar"},
			},
		},
	}
	for _, c := range cases {
		got, err := Tokenize(c.input)
		if err != nil {
			t.Errorf("Tokenize(%q) returned error: %v", c.input, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Tokenize(%q) returned\n%v, want\n%v", c.input, got, c.want)
		}
	}
}

func TestLexerInvalid(t *testing.T) {
	cases := []string{
		"select foo,, bar from baz",
		"select % from foo",
	}
	for _, input := range cases {
		_, err := Tokenize(input)
		if err == nil {
			t.Errorf("Tokenize(%q) did not return error", input)
		}
	}
}
