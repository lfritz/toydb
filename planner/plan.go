package planner

import (
	"errors"
	"fmt"

	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/storage"
)

var NotImplemented = errors.New("not implemented")

// Plan creates a query plan for the query.
func Plan(stmt *sql.SelectStatement, db *storage.Database) (query.Plan, error) {
	var err error
	var plan query.Plan

	switch f := stmt.From.(type) {
	case sql.TableName:
		table, err := db.Table(f.Name)
		if err != nil {
			return nil, err
		}
		plan = query.NewLoad(f.Name, table.Schema)
	case *sql.Join:
		// TODO
		return nil, NotImplemented
	default:
		panic(fmt.Sprintf("unexpected TableReference: %T", stmt.From))
	}

	switch what := stmt.What.(type) {
	case sql.Star:
		// ok
	case sql.ExpressionList:
		schema := plan.Schema()
		columns := make([]query.OutputColumn, len(what.Expressions))
		for i, e := range what.Expressions {
			converted, err := ConvertExpression(e, schema)
			if err != nil {
				return nil, err
			}
			columns[i].Expression = converted
			columns[i].Name = schema.Columns[i].Name
		}
		plan, err = query.NewProject(plan, columns)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("unexpected SelectList: %T", stmt.What))
	}

	if stmt.Where != nil {
		// TODO
		return nil, NotImplemented
	}

	return plan, nil
}
