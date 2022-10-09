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

type ColumnSchema struct {
	Name string
	Type Type
}
