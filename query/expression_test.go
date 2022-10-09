package query

import (
	"testing"

	"github.com/lfritz/toydb/types"
)

func sampleSchema() types.TableSchema {
	return types.TableSchema{
		Columns: []types.ColumnSchema{
			types.ColumnSchema{
				Name: "foo",
				Type: types.TypeBoolean,
			},
			types.ColumnSchema{
				Name: "bar",
				Type: types.TypeText,
			},
		},
	}
}

func sampleRow() *types.Row {
	return &types.Row{
		Schema: sampleSchema(),
		Values: []types.Value{types.NewBoolean(false), types.NewText("hello")},
	}
}

func TestContant(t *testing.T) {
	value := types.NewDecimal("123")
	c := NewConstant(value)

	gotType := c.Type()
	wantType := types.TypeDecimal
	if gotType != wantType {
		t.Errorf("c.Type() == %v, want %v", gotType, wantType)
	}

	got := c.Evaluate(sampleRow())
	want := value
	if got.Compare(want) != types.ComparedEq {
		t.Errorf("Evaluate returned %v, want %v", got, want)
	}
}

func TestColumnReference(t *testing.T) {
	c := NewColumnReference(1, types.TypeText)

	got := c.Evaluate(sampleRow())
	want := types.NewText("hello")
	if got.Compare(want) != types.ComparedEq {
		t.Errorf("Evaluate returned %v, want %v", got, want)
	}
}

func TestBinaryOperation(t *testing.T) {
	cases := []struct {
		left, right string
		op          BinaryOperator
		want        bool
	}{
		{"123", "123", BinaryOperatorEq, true},
		{"123", "456", BinaryOperatorEq, false},
		{"123", "456", BinaryOperatorGt, false},
		{"123", "456", BinaryOperatorLt, true},
		{"123", "456", BinaryOperatorLe, true},
		{"123", "456", BinaryOperatorGe, false},
	}

	row := sampleRow()
	for _, c := range cases {
		left := NewConstant(types.NewDecimal(c.left))
		right := NewConstant(types.NewDecimal(c.right))
		expression, err := NewBinaryOperation(left, right, c.op)
		if err != nil {
			t.Errorf("NewBinaryOperation(%v, %v, %v) returned error: %v", left, right, c.op, err)
			continue
		}
		want := types.NewBoolean(c.want)
		got := expression.Evaluate(row)
		if got.Compare(want) != types.ComparedEq {
			t.Errorf("Evaluate returned %v, want %v", got, c.want)
		}
	}

	left := NewConstant(types.NewDecimal("133"))
	right := NewConstant(types.NewBoolean(false))
	op := BinaryOperatorEq
	_, err := NewBinaryOperation(left, right, op)
	if err == nil {
		t.Errorf("NewBinaryOperation(%v, %v, %v) did not return error", left, right, op)
	}
}
