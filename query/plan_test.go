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
		NewConstant(types.Dat(1925, 1, 1)),
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
			{types.Dec("2"), types.Txt("The Kid"), types.Dat(1921, 1, 21), types.Dec("2")},
			{types.Dec("3"), types.Txt("Sherlock Jr."), types.Dat(1924, 4, 21), types.Dec("1")},
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
				types.ColumnSchema{"films.name", types.TypeText, false},
				types.ColumnSchema{"c1", types.TypeDecimal, false},
				types.ColumnSchema{"c2", types.TypeBoolean, false},
			},
		},
		Rows: [][]types.Value{
			{types.Txt("The General"), types.Dec("123"), types.Boo(false)},
			{types.Txt("The Kid"), types.Dec("123"), types.Boo(true)},
			{types.Txt("Sherlock Jr."), types.Dec("123"), types.Boo(true)},
		},
	}

	sampleData := storage.GetSampleData()
	l := NewLoad("films", sampleData.Films.Schema)
	comparison, err := NewBinaryOperation(
		NewColumnReference(2, types.TypeDate), // release_date
		BinaryOperatorLt,
		NewConstant(types.Dat(1925, 1, 1)),
	)
	if err != nil {
		t.Fatalf("NewBinaryOperation returned error: %v", err)
	}
	columns := []OutputColumn{
		OutputColumn{"films.name", NewColumnReference(1, types.TypeText)},
		OutputColumn{"c1", NewConstant(types.Dec("123"))},
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

func TestInnerJoin(t *testing.T) {
	sampleData := storage.GetSampleData()

	wantSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			types.ColumnSchema{"films.id", types.TypeDecimal, false},
			types.ColumnSchema{"films.name", types.TypeText, false},
			types.ColumnSchema{"films.release_date", types.TypeDate, false},
			types.ColumnSchema{"films.director", types.TypeDecimal, false},
			types.ColumnSchema{"people.id", types.TypeDecimal, false},
			types.ColumnSchema{"people.name", types.TypeText, false},
		},
	}
	wantRows := [][]types.Value{
		{types.Dec("1"), types.Txt("The General"), types.Dat(1926, 12, 31), types.Dec("1"), types.Dec("1"), types.Txt("Buster Keaton")},
		{types.Dec("2"), types.Txt("The Kid"), types.Dat(1921, 1, 21), types.Dec("2"), types.Dec("2"), types.Txt("Charlie Chaplin")},
		{types.Dec("3"), types.Txt("Sherlock Jr."), types.Dat(1924, 4, 21), types.Dec("1"), types.Dec("1"), types.Txt("Buster Keaton")},
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

func TestLeftOuterJoin(t *testing.T) {
	sampleData := storage.GetSampleData()

	wantSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			types.ColumnSchema{"people.id", types.TypeDecimal, false},
			types.ColumnSchema{"people.name", types.TypeText, false},
			types.ColumnSchema{"films.id", types.TypeDecimal, true},
			types.ColumnSchema{"films.name", types.TypeText, true},
			types.ColumnSchema{"films.release_date", types.TypeDate, true},
			types.ColumnSchema{"films.director", types.TypeDecimal, true},
		},
	}
	wantRows := [][]types.Value{
		{types.Dec("1"), types.Txt("Buster Keaton"), types.Dec("1"), types.Txt("The General"), types.Dat(1926, 12, 31), types.Dec("1")},
		{types.Dec("1"), types.Txt("Buster Keaton"), types.Dec("3"), types.Txt("Sherlock Jr."), types.Dat(1924, 4, 21), types.Dec("1")},
		{types.Dec("2"), types.Txt("Charlie Chaplin"), types.Dec("2"), types.Txt("The Kid"), types.Dat(1921, 1, 21), types.Dec("2")},
		{types.Dec("3"), types.Txt("Harold Lloyd"), types.NewNull(types.TypeDecimal), types.NewNull(types.TypeText), types.NewNull(types.TypeDate), types.NewNull(types.TypeDecimal)},
	}
	want := &types.Relation{
		Schema: wantSchema,
		Rows:   wantRows,
	}

	left := NewLoad("people", sampleData.People.Schema)
	right := NewLoad("films", sampleData.Films.Schema)
	condition, err := NewBinaryOperation(
		NewColumnReference(0, types.TypeDecimal),
		BinaryOperatorEq,
		NewColumnReference(5, types.TypeDecimal),
	)
	if err != nil {
		t.Fatalf("NewBinaryOperation returned error: %v", err)
	}
	join, err := NewJoin(JoinTypeLeftOuter, left, right, condition)
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

// TODO test right outer join
