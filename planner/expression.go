package planner

import (
	"fmt"

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
	case sql.ColumnReference:
		name := fmt.Sprintf("%s.%s", e.Relation, e.Name)
		i, t, ok := schema.Column(name)
		if !ok {
			return nil, fmt.Errorf("Column not found: %s", name)
		}
		return query.NewColumnReference(i, t), nil
	case *sql.BinaryOperation:
		left, err := ConvertExpression(e.Left, schema)
		if err != nil {
			return nil, err
		}
		right, err := ConvertExpression(e.Right, schema)
		if err != nil {
			return nil, err
		}
		operator := convertBinaryOperator(e.Operator)
		return query.NewBinaryOperation(left, right, operator)
	}
	return nil, NotImplemented
}

func convertBinaryOperator(o sql.BinaryOperator) query.BinaryOperator {
	switch o {
	case sql.BinaryOperatorEq:
		return query.BinaryOperatorEq
	case sql.BinaryOperatorNe:
		return query.BinaryOperatorNe
	case sql.BinaryOperatorLt:
		return query.BinaryOperatorLt
	case sql.BinaryOperatorGt:
		return query.BinaryOperatorGt
	case sql.BinaryOperatorLe:
		return query.BinaryOperatorLe
	case sql.BinaryOperatorGe:
		return query.BinaryOperatorGe
	}
	panic(fmt.Sprintf("unexpected value for BinaryOperator: %v", o))
}
