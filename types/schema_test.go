package types

import "testing"

func TestColumnSchemaString(t *testing.T) {
	schema := ColumnSchema{"name", TypeText, false}
	want := "name text not null"
	got := schema.String()
	if got != want {
		t.Errorf("schema.String() == %q, want %q", got, want)
	}
}

func TestTableSchema(t *testing.T) {
	schema := TableSchema{
		Columns: []ColumnSchema{
			ColumnSchema{"id", TypeDecimal, false},
			ColumnSchema{"name", TypeText, true},
		},
	}
	want := "TableSchema(id decimal not null, name text null)"
	got := schema.String()
	if got != want {
		t.Errorf("schema.String() == %q, want %q", got, want)
	}
}
