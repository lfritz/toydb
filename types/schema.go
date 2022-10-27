package types

import (
	"errors"
	"fmt"
	"strings"
)

type TableSchema struct {
	Columns []ColumnSchema
}

func (s TableSchema) Column(name string) (i int, t Type, ok bool) {
	for i, c := range s.Columns {
		if c.Name == name {
			return i, c.Type, true
		}
	}
	return
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

func (s TableSchema) Prefix(name string) TableSchema {
	columns := make([]ColumnSchema, len(s.Columns))
	for i, col := range s.Columns {
		columns[i] = ColumnSchema{
			Name: fmt.Sprintf("%s.%s", name, col.Name),
			Type: col.Type,
		}
	}
	return TableSchema{columns}
}

func (s TableSchema) String() string {
	list := make([]string, len(s.Columns))
	for i, c := range s.Columns {
		list[i] = c.String()
	}
	return fmt.Sprintf("TableSchema(%s)", strings.Join(list, ", "))
}

type ColumnSchema struct {
	Name string
	Type Type
	Null bool
}

func (s ColumnSchema) Check(value Value) error {
	if value.Type() != s.Type {
		return fmt.Errorf("wrong type: expected %v, got %v", s.Type, value.Type())
	}
	if value.Null() && !s.Null {
		return errors.New("value cannot be null")
	}
	return nil
}

func (s ColumnSchema) String() string {
	if s.Null {
		return fmt.Sprintf("%s %s null", s.Name, s.Type)
	} else {
		return fmt.Sprintf("%s %s not null", s.Name, s.Type)
	}
}
