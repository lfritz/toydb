package planner

import (
	"reflect"
	"testing"

	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/storage"
	"github.com/lfritz/toydb/types"
)

func TestConvertExpressionValid(t *testing.T) {
	sampleData := storage.GetSampleData()
	schema := sampleData.Films.Schema.Prefix("films")

	cases := []struct {
		input sql.Expression
		want  query.Expression
		name  string
	}{
		{
			sql.String{"hello"},
			query.NewConstant(types.Txt("hello")),
			"",
		},
		{
			sql.Boolean{true},
			query.NewConstant(types.Boo(true)),
			"",
		},
		{
			sql.Number{types.NewDecimal("123")},
			query.NewConstant(types.Dec("123")),
			"",
		},
		{
			sql.ColumnReference{"films", "name"},
			query.NewColumnReference(1, types.TypeText),
			"films.name",
		},
		{
			&sql.BinaryOperation{
				sql.Number{types.NewDecimal("4")},
				sql.BinaryOperatorEq,
				sql.ColumnReference{"films", "id"},
			},
			&query.BinaryOperation{
				query.NewConstant(types.Dec("4")),
				query.BinaryOperatorEq,
				query.NewColumnReference(0, types.TypeDecimal),
			},
			"",
		},
	}

	for _, c := range cases {
		got, name, err := ConvertExpression(c.input, schema)
		if err != nil {
			t.Errorf("ConvertExpression returned error: %v", err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ConvertExpression(%v, schema) returned %v, want %v", c.input, got, c.want)
		}
		if name != c.name {
			t.Errorf("CovnertExpression(%v, schema) returned name %q, want %q", c.input, name, c.name)
		}
	}
}

func TestConvertExpressionInvalid(t *testing.T) {
	sampleData := storage.GetSampleData()
	schema := sampleData.Films.Schema.Prefix("films")

	op := sql.BinaryOperatorEq
	four := sql.Number{types.NewDecimal("0")}
	cases := []sql.Expression{
		sql.ColumnReference{"foo", "id"},
		sql.ColumnReference{"films", "foo"},
		&sql.BinaryOperation{sql.ColumnReference{"foo", "id"}, op, four},
		&sql.BinaryOperation{sql.ColumnReference{"films", "name"}, op, four},
	}

	for _, c := range cases {
		_, _, err := ConvertExpression(c, schema)
		if err == nil {
			t.Errorf("ConvertExpression did not return error for %v", c)
		}
	}
}

func TestFindColumn(t *testing.T) {
	schema := types.TableSchema{
		Columns: []types.ColumnSchema{
			types.ColumnSchema{"films.name", types.TypeText, false},
			types.ColumnSchema{"films.release_date", types.TypeDate, false},
			types.ColumnSchema{"people.name", types.TypeText, false},
		},
	}

	name := "release_date"
	wantIndex := 1
	wantName := "films.release_date"
	index, name, err := FindColumn(name, schema)
	if err != nil {
		t.Errorf("FindColumn(%q) returned error: %v", name, err)
	}
	if index != wantIndex {
		t.Errorf("FindColumn(%q) returned index %v, want %v", name, index, wantIndex)
	}
	if name != wantName {
		t.Errorf("FindColumn(%q) returned name %v, want %v", name, name, wantName)
	}

	name = "name"
	_, _, err = FindColumn(name, schema)
	if err == nil {
		t.Errorf("FindColumn(%q) did not return error", name)
	}
}
