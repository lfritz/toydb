package sql

import (
	"fmt"
	"strings"
)

// A Statement is an SQL statement, e.g. a select query or an insert statement.
type Statement interface {
	String() string
}

// A SelectStatement is a "select ... from ..." query.
type SelectStatement struct {
	What  SelectList
	From  TableReference
	Where Expression
}

func (q SelectStatement) String() string {
	where := ""
	if q.Where != nil {
		where = fmt.Sprintf(", Where: %s", q.Where.String())
	}
	return fmt.Sprintf("SelectStatement(What: %s, From: %s%s)",
		q.What.String(),
		q.From.String(),
		where)
}

// A table reference defines a single table or multiple joined tables.
type TableReference interface {
	String() string
}

// A TableName is a TableReference that specifies a single table.
type TableName struct {
	Name string
}

func (t TableName) String() string {
	return fmt.Sprintf("Table(%s)", t.Name)
}

// A TableName is a TableReference that specifies a join.
type Join struct {
	Type      JoinType
	Left      TableReference
	Right     TableReference
	Condition Expression
}

func (j *Join) String() string {
	return fmt.Sprintf("Join(%s, %s, %s, %s)",
		j.Type.String(),
		j.Left.String(),
		j.Right.String(),
		j.Condition.String())
}

type JoinType int

const (
	JoinTypeInner JoinType = iota
	JoinTypeLeftOuter
	JoinTypeRightOuter
)

func (t JoinType) String() string {
	switch t {
	case JoinTypeInner:
		return "inner"
	case JoinTypeLeftOuter:
		return "left outer"
	case JoinTypeRightOuter:
		return "right outer"
	}
	return fmt.Sprintf("<unexpected join: %d>", t)
}

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

// A BinaryOperation is an operator with two operands.
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
