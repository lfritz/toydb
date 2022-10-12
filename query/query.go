package query

import (
	"fmt"

	"github.com/lfritz/toydb/storage"
	"github.com/lfritz/toydb/types"
)

// A Query implements the steps to run a query on a database.
type Query interface {
	Schema() types.TableSchema
	Run(db *storage.Database) *types.Relation
}

// A Load step loads a table from the database.
type Load struct {
	TableName   string
	TableSchema types.TableSchema
}

func NewLoad(name string, schema types.TableSchema) *Load {
	return &Load{
		TableName:   name,
		TableSchema: schema,
	}
}

func (l *Load) Schema() types.TableSchema {
	return l.TableSchema
}

func (l *Load) Run(db *storage.Database) *types.Relation {
	t, err := db.Table(l.TableName)
	if err != nil {
		panic(fmt.Sprintf("error loading table %s: %v", l.TableName, err))
	}
	return t
}

// A Select step selects the rows matching an expression.
type Select struct {
	From      Query
	Condition Expression
}

func NewSelect(from Query, condition Expression) (*Select, error) {
	if condition.Type() != types.TypeBoolean {
		return nil, fmt.Errorf("invalid condition for select step: %v", condition)
	}
	if err := condition.Check(from.Schema()); err != nil {
		return nil, err
	}
	return &Select{
		From:      from,
		Condition: condition,
	}, nil
}

func (s *Select) Schema() types.TableSchema {
	return s.From.Schema()
}

func (s *Select) Run(db *storage.Database) *types.Relation {
	from := s.From.Run(db)

	var rows [][]types.Value
	for i := range from.Rows {
		row := from.Row(i)
		got := s.Condition.Evaluate(row).(types.Boolean)
		if got.Bool() {
			rows = append(rows, row.Values)
		}
	}

	return &types.Relation{
		Schema: from.Schema,
		Rows:   rows,
	}
}

// TODO Project type

// TODO Join type
