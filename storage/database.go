package storage

import (
	"fmt"

	"github.com/lfritz/toydb/types"
)

type Database struct {
	tables map[string]*types.Relation
}

func (d *Database) Table(name string) (*types.Relation, error) {
	t, ok := d.tables[name]
	if !ok {
		return nil, fmt.Errorf("table not found: %s", name)
	}
	return t, nil
}
