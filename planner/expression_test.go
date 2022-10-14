package planner

import (
	"reflect"
	"testing"

	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/storage"
	"github.com/lfritz/toydb/types"
)

func TestConvertExpression(t *testing.T) {
	sampleData := storage.GetSampleData()
	schema := sampleData.Films.Schema

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
		// TODO ColumnReference
		// TODO BinaryOperation
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
