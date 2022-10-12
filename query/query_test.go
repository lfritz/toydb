package query

import (
	"reflect"
	"testing"

	"github.com/lfritz/toydb/storage"
	"github.com/lfritz/toydb/types"
)

func sampleRows() [][]types.Value {
	return [][]types.Value{
		[]types.Value{types.NewBoolean(false), types.NewText("hello")},
		[]types.Value{types.NewBoolean(true), types.NewText("ciao")},
	}
}

func sampleDatabase(t *testing.T) *storage.Database {
	db := storage.NewDatabase()
	table, err := db.CreateTable("mytable", sampleSchema())
	if err != nil {
		t.Fatalf("CreateTable returned error: %v", err)
	}

	for _, row := range sampleRows() {
		err := table.Insert(row)
		if err != nil {
			t.Fatalf("Insert returned error: %v", err)
		}
	}

	return db
}

func TestLoad(t *testing.T) {
	schema := sampleSchema()
	db := sampleDatabase(t)
	l := NewLoad("mytable", schema)
	got := l.Run(db)
	want := &types.Relation{
		Schema: schema,
		Rows:   sampleRows(),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Run returned %v, want %v", got, want)
	}
}

func TestSelect(t *testing.T) {
	schema := sampleSchema()
	db := sampleDatabase(t)
	l := NewLoad("mytable", schema)
	condition, err := NewBinaryOperation(
		NewColumnReference(0, types.TypeBoolean),
		NewConstant(types.NewBoolean(true)),
		BinaryOperatorEq,
	)
	if err != nil {
		t.Fatalf("NewBinaryOperation returned error: %v", err)
	}
	s, err := NewSelect(l, condition)
	if err != nil {
		t.Fatalf("NewSelect returned error: %v", err)
	}

	got := s.Run(db)
	want := &types.Relation{
		Schema: schema,
		Rows:   [][]types.Value{[]types.Value{types.NewBoolean(true), types.NewText("ciao")}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Run returned %v, want %v", got, want)
	}
}
