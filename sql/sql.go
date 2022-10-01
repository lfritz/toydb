package sql

type Statement interface{}

type SelectStatement struct {
	What  SelectList
	From  []FromExpression
	Where *Condition
}

type SelectList interface{}

type Star struct{}

type ExpressionList struct {
	Expressions []Expression
}

type Expression interface{}

type Column struct {
	Name string
}

type FromExpression interface{}

type TableName string

// TODO add type for joins

type Condition interface{}
