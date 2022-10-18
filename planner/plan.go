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
	plan, err := convertTableReference(stmt.From, db)
	if err != nil {
		return nil, err
	}

	if stmt.Where != nil {
		schema := plan.Schema()
		condition, err := ConvertExpression(stmt.Where, schema)
		if err != nil {
			return nil, err
		}
		plan, err = query.NewSelect(plan, condition)
		if err != nil {
			return nil, err
		}
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
			// TODO figure out the logic for setting the column name
			columns[i].Name = schema.Columns[i].Name
		}
		plan, err = query.NewProject(plan, columns)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("unexpected SelectList: %T", stmt.What))
	}

	return plan, nil
}

func convertTableReference(ref sql.TableReference, db *storage.Database) (query.Plan, error) {
	switch f := ref.(type) {
	case sql.TableName:
		table, err := db.Table(f.Name)
		if err != nil {
			return nil, err
		}
		return query.NewLoad(f.Name, table.Schema), nil
	case *sql.Join:
		joinType := convertJoinType(f.Type)
		left, err := convertTableReference(f.Left, db)
		if err != nil {
			return nil, err
		}
		right, err := convertTableReference(f.Right, db)
		if err != nil {
			return nil, err
		}
		schema := query.CombineSchemas(left.Schema(), right.Schema())
		condition, err := ConvertExpression(f.Condition, schema)
		if err != nil {
			return nil, err
		}
		return query.NewJoin(joinType, left, right, condition)
	}
	panic(fmt.Sprintf("unexpected TableReference: %T", ref))
}

func convertJoinType(input sql.JoinType) query.JoinType {
	switch input {
	case sql.JoinTypeInner:
		return query.JoinTypeInner
	case sql.JoinTypeLeftOuter:
		return query.JoinTypeLeftOuter
	case sql.JoinTypeRightOuter:
		return query.JoinTypeRightOuter
	}
	panic(fmt.Sprintf("unexpected JoinType: %d", input))
}
