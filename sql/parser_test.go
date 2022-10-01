package sql

import (
	"reflect"
	"testing"
)

func TestParserValid(t *testing.T) {
	cases := []struct {
		input string
		want  Statement
	}{
		{
			"select x, y from foo",
			&SelectStatement{
				What: []Expression{
					&Column{Name: "x"},
					&Column{Name: "y"},
				},
				From: []FromExpression{
					TableName("foo"),
				},
			},
		},
		{
			"select * from foo",
			&SelectStatement{
				What: new(Star),
				From: []FromExpression{
					TableName("foo"),
				},
			},
		},
	}
	for _, c := range cases {
		got, err := Parse(c.input)
		if err != nil {
			t.Errorf("Parse(%q) returned error: %v", c.input, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Parse(%q) returned\n%v, want\n%v", c.input, got, c.want)
		}
	}
}

func TestParserError(t *testing.T) {
	cases := []struct {
		input string
		want  SyntaxError
	}{
		{
			"",
			SyntaxError{Position: 0, Msg: `unexpected end of input`},
		},
		{
			"select % from foo",
			SyntaxError{Position: 7, Msg: `unexpected character: '%'`},
		},
		{
			"hello x from foo",
			SyntaxError{Position: 0, Msg: `got "hello", expected select`},
		},
		{
			"select x, from foo",
			SyntaxError{Position: 10, Msg: `got "from", expected star or identifier`},
		},
	}
	for _, c := range cases {
		_, err := Parse(c.input)
		if err != c.want {
			t.Errorf("Parse(%q) returned \n%#v, want\n%#v", c.input, err, c.want)
		}
	}
}
