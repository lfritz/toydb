package types

import "fmt"

type TableSchema struct {
	Columns []ColumnSchema
}

func (s TableSchema) Column(name string) (i int, t Type) {
	for i, c := range s.Columns {
		if c.Name == name {
			return i, c.Type
		}
	}
	panic(fmt.Sprintf("column not found: %s", name))
}

func (s TableSchema) Check(row []Value) error {
	if len(row) != len(s.Columns) {
		return fmt.Errorf("wrong number of values: expected %d, got %d", len(s.Columns), len(row))
	}
	for i := range s.Columns {
		if err := s.Columns[i].Check(row[i]); err != nil {
			return err
		}
	}
	return nil
}

type ColumnSchema struct {
	Name string
	Type Type
}

func (s ColumnSchema) Check(value Value) error {
	if s.Type != value.Type() {
		return fmt.Errorf("wrong type: expected %v, got %v", s.Type, value.Type())
	}
	return nil
}
