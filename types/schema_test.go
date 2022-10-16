package types

import "testing"

func TestColumnSchemaString(t *testing.T) {
	schema := ColumnSchema{"name", TypeText}
	want := "name text"
	got := schema.String()
	if got != want {
		t.Errorf("schema.String() == %q, want %q", got, want)
	}
}

func TestTableSchema(t *testing.T) {
	schema := TableSchema{
		Columns: []ColumnSchema{
			ColumnSchema{"id", TypeDecimal},
			ColumnSchema{"name", TypeText},
		},
	}
	want := "TableSchema(id decimal, name text)"
	got := schema.String()
	if got != want {
		t.Errorf("schema.String() == %q, want %q", got, want)
	}
}
