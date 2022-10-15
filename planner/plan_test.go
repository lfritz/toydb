package planner

import (
	"reflect"
	"testing"

	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/storage"
)

func parse(t *testing.T, input string) *sql.SelectStatement {
	statement, err := sql.Parse(input)
	if err != nil {
		t.Fatalf("sql.Parse returned error for %q: %v", input, err)
	}
	selectStatement, ok := statement.(*sql.SelectStatement)
	if !ok {
		t.Fatalf("Not a select statement: %q", input)
	}
	return selectStatement
}

func TestPlan(t *testing.T) {
	sampleData := storage.GetSampleData()

	cases := []struct {
		stmt string
		want query.Plan
	}{
		{
			"select * from films",
			query.NewLoad("films", sampleData.Films.Schema),
		},
		// TODO
	}

	for _, c := range cases {
		stmt := parse(t, c.stmt)
		got, err := Plan(stmt, sampleData.Database)
		if err != nil {
			t.Fatalf("Plan returned error: %v", err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Query plan for\n%s\nis:\n%q\nwant:\n%v\n", c.stmt, got, c.want)
		}
	}
}
