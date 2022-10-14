package planner

import (
	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/types"
)

func ConvertExpression(input sql.Expression, schema types.TableSchema) (query.Expression, error) {
	switch e := input.(type) {
	case sql.String:
		return query.NewConstant(types.NewText(e.Value)), nil
	case sql.Boolean:
		return query.NewConstant(types.NewBoolean(e.Value)), nil
	case sql.Number:
		return query.NewConstant(e.Value), nil
	}
	return nil, NotImplemented
}
