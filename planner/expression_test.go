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
	}{
		{
			sql.String{"hello"},
			query.NewConstant(types.NewText("hello")),
		},
		{
			sql.Boolean{true},
			query.NewConstant(types.NewBoolean(true)),
		},
		{
			sql.Number{types.NewDecimal("123")},
			query.NewConstant(types.NewDecimal("123")),
		},
		{
			sql.ColumnReference{"films", "name"},
			query.NewColumnReference(1, types.TypeText),
		},
		{
			&sql.BinaryOperation{
				sql.Number{types.NewDecimal("4")},
				sql.BinaryOperatorEq,
				sql.ColumnReference{"films", "id"},
			},
			&query.BinaryOperation{
				query.NewConstant(types.NewDecimal("4")),
				query.BinaryOperatorEq,
				query.NewColumnReference(0, types.TypeDecimal),
			},
		},
	}

	for _, c := range cases {
		got, err := ConvertExpression(c.input, schema)
		if err != nil {
			t.Errorf("ConvertExpression returned error: %v", err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ConvertExpression(%v, schema) returned %v, want %v", c.input, got, c.want)
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
		sql.ColumnReference{"", "id"},
		&sql.BinaryOperation{sql.ColumnReference{"foo", "id"}, op, four},
		&sql.BinaryOperation{sql.ColumnReference{"films", "name"}, op, four},
	}

	for _, c := range cases {
		_, err := ConvertExpression(c, schema)
		if err == nil {
			t.Errorf("ConvertExpression did not return error for %v", c)
		}
	}
}
