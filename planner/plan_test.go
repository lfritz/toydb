package planner

import (
	"reflect"
	"testing"

	"github.com/lfritz/toydb/query"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/storage"
	"github.com/lfritz/toydb/types"
)

func parse(t *testing.T, input string) *sql.SelectStatement {
	statement, err := sql.Parse(input)
	if err != nil {
		t.Fatalf("sql.Parse returned error for %q: %v", input, err)
	}
	return statement
}

func TestPlanValid(t *testing.T) {
	sampleData := storage.GetSampleData()

	join, err := query.NewJoin(
		query.JoinTypeInner,
		query.NewLoad("films", sampleData.Films.Schema),
		query.NewLoad("people", sampleData.People.Schema),
		&query.BinaryOperation{
			&query.ColumnReference{3, types.TypeDecimal},
			query.BinaryOperatorEq,
			&query.ColumnReference{4, types.TypeDecimal},
		},
	)
	if err != nil {
		t.Fatalf("query.NewJoin returned error: %v", err)
	}

	cases := []struct {
		stmt string
		want query.Plan
	}{
		{
			"select * from films",
			query.NewLoad("films", sampleData.Films.Schema),
		},
		{
			"select id, name, release_date, director from films",
			&query.Project{
				From: query.NewLoad("films", sampleData.Films.Schema),
				Columns: []query.OutputColumn{
					query.OutputColumn{"films.id", &query.ColumnReference{0, types.TypeDecimal}},
					query.OutputColumn{"films.name", &query.ColumnReference{1, types.TypeText}},
					query.OutputColumn{"films.release_date", &query.ColumnReference{2, types.TypeDate}},
					query.OutputColumn{"films.director", &query.ColumnReference{3, types.TypeDecimal}},
				},
			},
		},
		{
			"select id from films where name = 'The General'",
			&query.Project{
				From: &query.Select{
					From: query.NewLoad("films", sampleData.Films.Schema),
					Condition: &query.BinaryOperation{
						&query.ColumnReference{1, types.TypeText},
						query.BinaryOperatorEq,
						query.NewConstant(types.Txt("The General")),
					},
				},
				Columns: []query.OutputColumn{
					query.OutputColumn{"films.id", &query.ColumnReference{0, types.TypeDecimal}},
				},
			},
		},
		{
			"select films.id, people.id from films join people on films.director = people.id",
			&query.Project{
				From: join,
				Columns: []query.OutputColumn{
					query.OutputColumn{"films.id", &query.ColumnReference{0, types.TypeDecimal}},
					query.OutputColumn{"people.id", &query.ColumnReference{4, types.TypeDecimal}},
				},
			},
		},
	}

	for _, c := range cases {
		stmt := parse(t, c.stmt)
		got, err := Plan(stmt, sampleData.Database)
		if err != nil {
			t.Fatalf("Plan returned error: %v", err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Query plan for\n%s\nis:\n%swant:\n%v",
				c.stmt, query.Print(got), query.Print(c.want))
		}
	}
}

func TestPlanInvalid(t *testing.T) {
	sampleData := storage.GetSampleData()
	cases := []string{
		"select * from foo",
		"select foo from films",
		"select id from films where foo = 123",
	}
	for _, c := range cases {
		stmt := parse(t, c)
		_, err := Plan(stmt, sampleData.Database)
		if err == nil {
			t.Fatalf("Plan did not return error for: %s", c)
		}
	}
}
