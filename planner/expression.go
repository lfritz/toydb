package planner

import (
	"fmt"
	"strings"

	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/types"
)

func ConvertExpression(input sql.Expression, schema types.TableSchema) (query.Expression, error) {
	switch e := input.(type) {
	case sql.ColumnReference:
		return convertColumnReference(e, schema)
	case sql.String:
		return query.NewConstant(types.NewText(e.Value)), nil
	case sql.Boolean:
		return query.NewConstant(types.NewBoolean(e.Value)), nil
	case sql.Number:
		return query.NewConstant(e.Value), nil
	case *sql.BinaryOperation:
		return convertBinaryOperation(e, schema)
	}
	panic(fmt.Sprintf("unexpected sql.Expression: %T", input))
}

func convertColumnReference(r sql.ColumnReference, schema types.TableSchema) (*query.ColumnReference, error) {
	if r.Relation == "" {
		index, err := FindColumn(r.Name, schema)
		if err != nil {
			return nil, err
		}
		return query.NewColumnReference(index, schema.Columns[index].Type), nil
	} else {
		name := fmt.Sprintf("%s.%s", r.Relation, r.Name)
		index, t, ok := schema.Column(name)
		if !ok {
			return nil, fmt.Errorf("Column not found: %s", name)
		}
		return query.NewColumnReference(index, t), nil
	}
}

func convertBinaryOperation(o *sql.BinaryOperation, schema types.TableSchema) (*query.BinaryOperation, error) {
	left, err := ConvertExpression(o.Left, schema)
	if err != nil {
		return nil, err
	}
	right, err := ConvertExpression(o.Right, schema)
	if err != nil {
		return nil, err
	}
	operator := convertBinaryOperator(o.Operator)
	return query.NewBinaryOperation(left, right, operator)
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

func FindColumn(name string, schema types.TableSchema) (int, error) {
	suffix := fmt.Sprintf(".%s", name)
	var index int
	var found bool
	for i, col := range schema.Columns {
		if strings.HasSuffix(col.Name, suffix) {
			if found {
				return 0, fmt.Errorf("ambiguous column reference: %s", name)
			}
			index = i
			found = true
		}
	}
	if !found {
		return 0, fmt.Errorf("column not found: %s", name)
	}
	return index, nil
}
