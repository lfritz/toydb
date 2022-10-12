package types

// A Relation represents a table or the result of a query.
type Relation struct {
	Schema TableSchema
	Rows   [][]Value
}

func (r *Relation) Row(i int) *Row {
	return &Row{
		Schema: r.Schema,
		Values: r.Rows[i],
	}
}

func (r *Relation) Insert(row []Value) error {
	if err := r.Schema.Check(row); err != nil {
		return err
	}
	r.Rows = append(r.Rows, row)
	return nil
}

type Row struct {
	Schema TableSchema
	Values []Value
}
