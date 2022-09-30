package sql

import (
	"reflect"
	"testing"
)

func TestLexerValid(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		{
			"",
			"",
		},
		{
			"select foo, bar, baz from qux;",
			`select (identifier "foo") comma (identifier "bar") comma (identifier "baz") from (identifier "qux") semicolon`,
		},
		{
			"select\tfoo  ,bar\t,baz from \t \n qux;   ",
			`select (identifier "foo") comma (identifier "bar") comma (identifier "baz") from (identifier "qux") semicolon`,
		},
		{
			"SELECT _foo, bar123, a_b$c FROM qux",
			`select (identifier "_foo") comma (identifier "bar123") comma (identifier "a_b$c") from (identifier "qux")`,
		},
		{
			"select foo from bar where (x = 123.45 or y < 0) and z >= .4",
			`select (identifier "foo") from (identifier "bar") where openparen (identifier "x") eq (number "123.45") or (identifier "y") lt (number "0") closeparen and (identifier "z") ge (number ".4")`,
		},
		{
			"select * from foo where x is not null",
			`select star from (identifier "foo") where (identifier "x") is not null`,
		},
		{
			"select foo from bar where x != 'hello' or y <> 'ciao'",
			`select (identifier "foo") from (identifier "bar") where (identifier "x") ne (string "hello") or (identifier "y") ne (string "ciao")`,
		},
		{
			"select foo.x, bar.y from foo left outer join bar",
			`select (identifier "foo") dot (identifier "x") comma (identifier "bar") dot (identifier "y") from (identifier "foo") left outer join (identifier "bar")`,
		},
	}
	for _, c := range cases {
		tokens, err := Tokenize(c.input)
		if err != nil {
			t.Errorf("Tokenize(%q) returned error: %v", c.input, err)
			continue
		}
		got := PrintTokens(tokens)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Tokenize(%q) returned\n%s, want\n%s", c.input, got, c.want)
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
