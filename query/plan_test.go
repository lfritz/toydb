package query

import (
	"reflect"
	"testing"

	"github.com/lfritz/toydb/storage"
	"github.com/lfritz/toydb/types"
)

func TestLoad(t *testing.T) {
	sampleData := storage.GetSampleData()
	l := NewLoad("films", sampleData.Films.Schema)
	got := l.Run(sampleData.Database)
	want := sampleData.Films
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Run returned %v, want %v", got, want)
	}
}

func TestSelect(t *testing.T) {
	sampleData := storage.GetSampleData()
	l := NewLoad("films", sampleData.Films.Schema)
	condition, err := NewBinaryOperation(
		NewColumnReference(2, types.TypeDate), // release_date
		BinaryOperatorLt,
		NewConstant(types.NewDate(1925, 1, 1)),
	)
	if err != nil {
		t.Fatalf("NewBinaryOperation returned error: %v", err)
	}
	s, err := NewSelect(l, condition)
	if err != nil {
		t.Fatalf("NewSelect returned error: %v", err)
	}

	got := s.Run(sampleData.Database)
	want := &types.Relation{
		Schema: sampleData.Films.Schema,
		Rows: [][]types.Value{
			{types.NewDecimal("2"), types.NewText("The Kid"), types.NewDate(1921, 1, 21), types.NewDecimal("2")},
			{types.NewDecimal("3"), types.NewText("Sherlock Jr."), types.NewDate(1924, 4, 21), types.NewDecimal("1")},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Run returned %v, want %v", got, want)
	}
}

func TestProject(t *testing.T) {
	want := &types.Relation{
		Schema: types.TableSchema{
			Columns: []types.ColumnSchema{
				types.ColumnSchema{"films.name", types.TypeText},
				types.ColumnSchema{"c1", types.TypeDecimal},
				types.ColumnSchema{"c2", types.TypeBoolean},
			},
		},
		Rows: [][]types.Value{
			{types.NewText("The General"), types.NewDecimal("123"), types.NewBoolean(false)},
			{types.NewText("The Kid"), types.NewDecimal("123"), types.NewBoolean(true)},
			{types.NewText("Sherlock Jr."), types.NewDecimal("123"), types.NewBoolean(true)},
		},
	}

	sampleData := storage.GetSampleData()
	l := NewLoad("films", sampleData.Films.Schema)
	comparison, err := NewBinaryOperation(
		NewColumnReference(2, types.TypeDate), // release_date
		BinaryOperatorLt,
		NewConstant(types.NewDate(1925, 1, 1)),
	)
	if err != nil {
		t.Fatalf("NewBinaryOperation returned error: %v", err)
	}
	columns := []OutputColumn{
		OutputColumn{"films.name", NewColumnReference(1, types.TypeText)},
		OutputColumn{"c1", NewConstant(types.NewDecimal("123"))},
		OutputColumn{"c2", comparison},
	}
	p, err := NewProject(l, columns)
	if err != nil {
		t.Fatalf("NewProject returned error: %v", err)
	}
	got := p.Run(sampleData.Database)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Run returned %v, want %v", got, want)
	}
}

func TestJoin(t *testing.T) {
	sampleData := storage.GetSampleData()

	wantSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			types.ColumnSchema{"films.id", types.TypeDecimal},
			types.ColumnSchema{"films.name", types.TypeText},
			types.ColumnSchema{"films.release_date", types.TypeDate},
			types.ColumnSchema{"films.director", types.TypeDecimal},
			types.ColumnSchema{"people.id", types.TypeDecimal},
			types.ColumnSchema{"people.name", types.TypeText},
		},
	}
	wantRows := [][]types.Value{
		{types.NewDecimal("1"), types.NewText("The General"), types.NewDate(1926, 12, 31), types.NewDecimal("1"), types.NewDecimal("1"), types.NewText("Buster Keaton")},
		{types.NewDecimal("2"), types.NewText("The Kid"), types.NewDate(1921, 1, 21), types.NewDecimal("2"), types.NewDecimal("2"), types.NewText("Charlie Chaplin")},
		{types.NewDecimal("3"), types.NewText("Sherlock Jr."), types.NewDate(1924, 4, 21), types.NewDecimal("1"), types.NewDecimal("1"), types.NewText("Buster Keaton")},
	}
	want := &types.Relation{
		Schema: wantSchema,
		Rows:   wantRows,
	}

	left := NewLoad("films", sampleData.Films.Schema)
	right := NewLoad("people", sampleData.People.Schema)
	condition, err := NewBinaryOperation(
		NewColumnReference(3, types.TypeDecimal),
		BinaryOperatorEq,
		NewColumnReference(4, types.TypeDecimal),
	)
	if err != nil {
		t.Fatalf("NewBinaryOperation returned error: %v", err)
	}
	join, err := NewJoin(JoinTypeInner, left, right, condition)
	if err != nil {
		t.Fatalf("NewJoin returned error: %v", err)
	}

	gotSchema := join.Schema()
	if !reflect.DeepEqual(gotSchema, wantSchema) {
		t.Errorf("Schema returned %v, want %v", gotSchema, wantSchema)
	}

	got := join.Run(sampleData.Database)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Run returned %v, want %v", got, want)
	}
}
