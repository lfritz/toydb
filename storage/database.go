package storage

import (
	"fmt"

	"github.com/lfritz/toydb/types"
)

type Database struct {
	tables map[string]*types.Relation
}

func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*types.Relation),
	}
}

func (d *Database) Table(name string) (*types.Relation, error) {
	t, ok := d.tables[name]
	if !ok {
		return nil, fmt.Errorf("table not found: %s", name)
	}
	return t, nil
}

func (d *Database) CreateTable(name string, schema types.TableSchema) (*types.Relation, error) {
	_, exists := d.tables[name]
	if exists {
		return nil, fmt.Errorf("table already exists: %s", name)
	}
	table := &types.Relation{
		Schema: schema,
	}
	d.tables[name] = table
	return table, nil
}
