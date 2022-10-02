package sql

import "fmt"

// An Expression is an SQL expression, for example the expression in a "where" clause.
type Expression interface {
	String() string
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
