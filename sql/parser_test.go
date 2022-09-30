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
				What: []Expression{new(Star)},
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

func TestParserInvalid(t *testing.T) {
	cases := []string{
		"",
		"select % from foo",
		"hello x from foo",
		"select x, from foo",
	}
	for _, input := range cases {
		_, err := Parse(input)
		if err == nil {
			t.Errorf("Parse(%q) did not return error", input)
		}
	}
}
