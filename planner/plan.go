package planner

import (
	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/storage"
)

// Plan creates a query plan for the query.
func Plan(query *sql.SelectStatement, db *storage.Database) (*query.Plan, error) {
	// TODO
	return nil, nil
}
