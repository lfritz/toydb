package planner

import (
	"fmt"
	"strings"

	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/types"
)

func ConvertExpression(input sql.Expression, schema types.TableSchema) (query.Expression, string, error) {
	switch e := input.(type) {
	case sql.ColumnReference:
		return convertColumnReference(e, schema)
	case sql.String:
		return query.NewConstant(types.NewValue(types.NewText(e.Value))), "", nil
	case sql.Boolean:
		return query.NewConstant(types.NewValue(types.NewBoolean(e.Value))), "", nil
	case sql.Number:
		return query.NewConstant(types.NewValue(e.Value)), "", nil
	case sql.Date:
		return query.NewConstant(types.NewValue(e.Value)), "", nil
	case *sql.BinaryOperation:
		return convertBinaryOperation(e, schema)
	}
	panic(fmt.Sprintf("unexpected sql.Expression: %T", input))
}

func convertColumnReference(r sql.ColumnReference, schema types.TableSchema) (*query.ColumnReference, string, error) {
	if r.Relation == "" {
		index, name, err := FindColumn(r.Name, schema)
		if err != nil {
			return nil, "", err
		}
		return query.NewColumnReference(index, schema.Columns[index].Type), name, nil
	} else {
		name := fmt.Sprintf("%s.%s", r.Relation, r.Name)
		index, t, ok := schema.Column(name)
		if !ok {
			return nil, "", fmt.Errorf("Column not found: %s", name)
		}
		return query.NewColumnReference(index, t), name, nil
	}
}

func convertBinaryOperation(o *sql.BinaryOperation, schema types.TableSchema) (*query.BinaryOperation, string, error) {
	left, _, err := ConvertExpression(o.Left, schema)
	if err != nil {
		return nil, "", err
	}
	right, _, err := ConvertExpression(o.Right, schema)
	if err != nil {
		return nil, "", err
	}
	operator := convertBinaryOperator(o.Operator)
	expression, err := query.NewBinaryOperation(left, operator, right)
	return expression, "", err
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

func FindColumn(input string, schema types.TableSchema) (index int, name string, err error) {
	suffix := fmt.Sprintf(".%s", input)
	var found bool
	for i, col := range schema.Columns {
		if strings.HasSuffix(col.Name, suffix) {
			if found {
				err = fmt.Errorf("ambiguous column reference: %s", input)
				return
			}
			index = i
			name = col.Name
			found = true
		}
	}
	if !found {
		err = fmt.Errorf("column not found: %s", input)
		return
	}
	return
}
