package toydb

import (
	"reflect"
	"testing"

	"github.com/lfritz/toydb/planner"
	"github.com/lfritz/toydb/sql"
	"github.com/lfritz/toydb/storage"
	"github.com/lfritz/toydb/types"
)

func TestAll(t *testing.T) {
	sampleData := storage.GetSampleData()
	query := `select films.name, people.name from films join people on films.director = people.id`
	want := &types.Relation{
		Schema: types.TableSchema{
			Columns: []types.ColumnSchema{
				{"films.name", types.TypeText, false},
				{"people.name", types.TypeText, false},
			},
		},
		Rows: [][]types.Value{
			[]types.Value{types.Txt("The General"), types.Txt("Buster Keaton")},
			[]types.Value{types.Txt("The Kid"), types.Txt("Charlie Chaplin")},
			[]types.Value{types.Txt("Sherlock Jr."), types.Txt("Buster Keaton")},
		},
	}

	stmt, err := sql.Parse(query)
	if err != nil {
		t.Fatalf("sql.Parse returned error: %v", err)
	}

	plan, err := planner.Plan(stmt, sampleData.Database)
	if err != nil {
		t.Fatalf("planner.Plan returned error: %v", err)
	}

	got := plan.Run(sampleData.Database)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Query result for\n%s\ngot:\n%s\nwant:\n%s\n", query, got, want)
	}
}
