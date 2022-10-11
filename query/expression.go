package query

import (
	"fmt"

	"github.com/lfritz/toydb/types"
)

// An Expression is an expression composed of column references, constants, and operations on them.
// Each expression has a static type.
type Expression interface {
	Type() types.Type
	Check(schema *types.TableSchema) error
	Evaluate(r *types.Row) types.Value
}

type Constant struct {
	value types.Value
}

func NewConstant(value types.Value) *Constant {
	return &Constant{value: value}
}

func (c Constant) Type() types.Type {
	return c.value.Type()
}

func (c Constant) Check(schema *types.TableSchema) error {
	return nil
}

func (c Constant) Evaluate(r *types.Row) types.Value {
	return c.value
}

type ColumnReference struct {
	Index int
	T     types.Type
}

func NewColumnReference(index int, t types.Type) *ColumnReference {
	return &ColumnReference{
		Index: index,
		T:     t,
	}
}

func (c *ColumnReference) Type() types.Type {
	return c.T
}

func (c ColumnReference) Check(schema *types.TableSchema) error {
	return nil // TODO
}

func (c *ColumnReference) Evaluate(r *types.Row) types.Value {
	return r.Values[c.Index]
}

type BinaryOperation struct {
	Left     Expression
	Operator BinaryOperator
	Right    Expression
}

func NewBinaryOperation(left, right Expression, op BinaryOperator) (*BinaryOperation, error) {
	if left.Type() != right.Type() {
		return nil, fmt.Errorf("incompatible types: %v, %v", left.Type(), right.Type())
	}
	return &BinaryOperation{
		Left:     left,
		Operator: op,
		Right:    right,
	}, nil
}

func (o *BinaryOperation) Type() types.Type {
	return types.TypeBoolean
}

func (o BinaryOperation) Check(schema *types.TableSchema) error {
	return nil // TODO
}

func (o *BinaryOperation) Evaluate(r *types.Row) types.Value {
	left := o.Left.Evaluate(r)
	right := o.Right.Evaluate(r)
	var result bool
	switch left.Compare(right) {
	case types.ComparedLt:
		result = o.Operator == BinaryOperatorLt || o.Operator == BinaryOperatorLe
	case types.ComparedEq:
		result = o.Operator == BinaryOperatorLe || o.Operator == BinaryOperatorEq || o.Operator == BinaryOperatorGe
	case types.ComparedGt:
		result = o.Operator == BinaryOperatorGt || o.Operator == BinaryOperatorGe
	default: // ComparedInvalid
		panic("comparison returned ComparedInvalid")
	}
	return types.NewBoolean(result)
}

type BinaryOperator int

const (
	// comparison operators
	BinaryOperatorEq BinaryOperator = iota
	BinaryOperatorNe
	BinaryOperatorLt
	BinaryOperatorGt
	BinaryOperatorLe
	BinaryOperatorGe
)
