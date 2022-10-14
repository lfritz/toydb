package planner

import (
	"errors"

	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/storage"
)

var NotImplemented = errors.New("not implemented")

// Plan creates a query plan for the query.
func Plan(stmt *sql.SelectStatement, db *storage.Database) (query.Plan, error) {
	var plan query.Plan

	switch f := stmt.From.(type) {
	case sql.TableName:
		table, err := db.Table(f.Name)
		if err != nil {
			return nil, err
		}
		plan = query.NewLoad(f.Name, table.Schema)
	default:
		return nil, NotImplemented
	}

	switch stmt.What.(type) {
	case sql.Star:
		// ok
	default:
		return nil, NotImplemented
	}

	if stmt.Where != nil {
		return nil, NotImplemented
	}

	return plan, nil
}
