package query

import (
	"fmt"

	"github.com/lfritz/toydb/storage"
	"github.com/lfritz/toydb/types"
)

// A Plan implements the steps to run a query on a database.
type Plan interface {
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
	From      Plan
	Condition Expression
}

func NewSelect(from Plan, condition Expression) (*Select, error) {
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

type OutputColumn struct {
	Name       string
	Expression Expression
}

func SimpleColumn(name string, index int, t types.Type) OutputColumn {
	return OutputColumn{
		Name:       name,
		Expression: NewColumnReference(index, t),
	}
}

func ComputedColumn(name string, expression Expression) OutputColumn {
	return OutputColumn{
		Name:       name,
		Expression: expression,
	}
}

func (c OutputColumn) Schema() types.ColumnSchema {
	return types.ColumnSchema{
		Name: c.Name,
		Type: c.Expression.Type(),
	}
}

// A Project step produces a new set of columns.
type Project struct {
	From    Plan
	Columns []OutputColumn
}

func NewProject(from Plan, columns []OutputColumn) (*Project, error) {
	// check for duplicate column names
	names := make(map[string]bool)
	for _, c := range columns {
		if names[c.Name] {
			return nil, fmt.Errorf("duplicate column: %s", c.Name)
		}
		names[c.Name] = true
	}

	// check expressions
	for _, c := range columns {
		if err := c.Expression.Check(from.Schema()); err != nil {
			return nil, err
		}
	}

	project := &Project{
		From:    from,
		Columns: columns,
	}
	return project, nil
}

func (p *Project) Schema() types.TableSchema {
	columns := make([]types.ColumnSchema, len(p.Columns))
	for i, c := range p.Columns {
		columns[i] = c.Schema()
	}
	return types.TableSchema{
		Columns: columns,
	}
}

func (p *Project) Run(db *storage.Database) *types.Relation {
	from := p.From.Run(db)
	rows := make([][]types.Value, len(from.Rows))
	for i := range from.Rows {
		row := make([]types.Value, len(p.Columns))
		for j := range p.Columns {
			row[j] = p.Columns[j].Expression.Evaluate(from.Row(i))
		}
		rows[i] = row
	}
	return &types.Relation{
		Schema: p.Schema(),
		Rows:   rows,
	}
}

// TODO Join type
