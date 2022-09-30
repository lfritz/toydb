package sql

type Statement interface{}

type SelectStatement struct {
	What  []Expression
	From  []FromExpression
	Where *Condition
}

type Expression interface{}

type Star struct{}

type Column struct {
	Name string
}

type FromExpression interface{}

type TableName string

// TODO add type for joins

type Condition interface{}
