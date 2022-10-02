package sql

import (
	"fmt"
	"strings"
)

// A SelectList is the part of an SQL query between the "select" and "from" keywords.
type SelectList interface {
	String() string
}

// Star represents the star in "select * from ...".
type Star struct{}

func (s Star) String() string {
	return "Star"
}

// ExpressionList represents the comma-separated list of expressions in an SQL query.
type ExpressionList struct {
	Expressions []Expression
}

func (l ExpressionList) String() string {
	list := make([]string, len(l.Expressions))
	for i, e := range l.Expressions {
		list[i] = e.String()
	}
	return fmt.Sprintf("ExpressionList(%s)", strings.Join(list, ", "))
}

// An Expression is an SQL expression, for example the expression in a "where" clause.
type Expression interface {
	String() string
}

// A ColumnReference is the name of a column, optionally with the name of the relation.
type ColumnReference struct {
	Relation string
	Name     string
}

func (r ColumnReference) String() string {
	if r.Relation == "" {
		return fmt.Sprintf("Column(%s)", r.Name)
	}
	return fmt.Sprintf("Column(%s.%s)", r.Relation, r.Name)
}

// A String is an SQL string literal.
type String struct {
	Value string
}

func (s String) String() string {
	return fmt.Sprintf("String(%q)", s.Value)
}

// A Number is a decimal number value.
type Number struct {
	Value Decimal
}

func (n Number) String() string {
	return fmt.Sprintf("Number(%v)", n.Value)
}

// A BinaryOperation is an expression with a binary operator, for example "1 + 2" or "foo = 'bar'".
type BinaryOperation struct {
	Left     Expression
	Operator BinaryOperator
	Right    Expression
}

func (o *BinaryOperation) String() string {
	return fmt.Sprintf("BinaryOperation(Left: %s, Operator: %s, Right: %s)",
		o.Left.String(),
		o.Operator.String(),
		o.Right.String())
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

var binaryOperatorNames = map[BinaryOperator]string{
	BinaryOperatorEq: "Eq",
	BinaryOperatorNe: "Ne",
	BinaryOperatorLt: "Lt",
	BinaryOperatorGt: "Gt",
	BinaryOperatorLe: "Le",
	BinaryOperatorGe: "Ge",
}

func (o BinaryOperator) String() string {
	if name, ok := binaryOperatorNames[o]; ok {
		return name
	}
	return fmt.Sprintf("unexpected binary operator: %d", o)
}
