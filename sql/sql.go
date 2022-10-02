package sql

import "fmt"

// An Expression is an SQL expression, for example the expression in a "where" clause.
type Expression interface {
	PrintExpression() string
}

// A String is an SQL string literal.
type String struct {
	Value string
}

func (s String) PrintExpression() string {
	return fmt.Sprintf("String(%q)", s.Value)
}

// A Number is a decimal number value.
type Number struct {
	Value Decimal
}

func (n Number) PrintExpression() string {
	return fmt.Sprintf("Number(%v)", n.Value)
}
