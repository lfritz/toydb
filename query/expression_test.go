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
		Values: []types.Value{types.Boo(false), types.Txt("hello")},
	}
}

func TestConstantType(t *testing.T) {
	c := NewConstant(types.Dec("123"))
	got := c.Type()
	want := types.TypeDecimal
	if got != want {
		t.Errorf("c.Type() == %v, want %v", got, want)
	}
}

func TestConstantEvaluate(t *testing.T) {
	value := types.Dec("123")
	got := NewConstant(value).Evaluate(sampleRow())
	if got.Compare(value) != types.ComparedEq {
		t.Errorf("Evaluate returned %v, want %v", got, value)
	}
}

func TestColumnReferenceType(t *testing.T) {
	want := types.TypeText
	c := NewColumnReference(1, want)
	got := c.Type()
	if got != want {
		t.Errorf("c.Type() == %v, want %v", got, want)
	}
}

func TestColumnReferenceCheck(t *testing.T) {
	cases := []struct {
		index     int
		typ       types.Type
		wantError bool
	}{
		{1, types.TypeText, false},
		{1, types.TypeBoolean, true},
		{2, types.TypeBoolean, true},
	}
	schema := sampleSchema()
	for _, c := range cases {
		err := NewColumnReference(c.index, c.typ).Check(schema)
		if err != nil && !c.wantError {
			t.Errorf("NewColumnReference(%v, %v).Check returned error: %v", c.index, c.typ, err)
		}
		if err == nil && c.wantError {
			t.Errorf("NewColumnReference(%v, %v).Check did not return error", c.index, c.typ)
		}
	}
}

func TestColumnReferenceEvaluate(t *testing.T) {
	c := NewColumnReference(1, types.TypeText)
	got := c.Evaluate(sampleRow())
	want := types.Txt("hello")
	if got.Compare(want) != types.ComparedEq {
		t.Errorf("Evaluate returned %v, want %v", got, want)
	}
}

func binaryOperation(t *testing.T, left, right string, op BinaryOperator) *BinaryOperation {
	l := NewConstant(types.Dec(left))
	r := NewConstant(types.Dec(right))
	expression, err := NewBinaryOperation(l, op, r)
	if err != nil {
		t.Fatalf("NewBinaryOperation(%v, %v, %v) returned error: %v", l, op, r, err)
	}
	return expression
}

func TestBinaryOperationType(t *testing.T) {
	expression := binaryOperation(t, "123", "456", BinaryOperatorNe)
	want := types.TypeBoolean
	got := expression.Type()
	if got != want {
		t.Errorf("expression.Type() == %v, want %v", got, want)
	}
}

func TestBinaryOperationEvaluate(t *testing.T) {
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
		expression := binaryOperation(t, c.left, c.right, c.op)
		want := types.Boo(c.want)
		got := expression.Evaluate(row)
		if got.Compare(want) != types.ComparedEq {
			t.Errorf("Evaluate returned %v, want %v", got, c.want)
		}
	}

	left := NewConstant(types.Dec("133"))
	right := NewConstant(types.Boo(false))
	op := BinaryOperatorEq
	_, err := NewBinaryOperation(left, op, right)
	if err == nil {
		t.Errorf("NewBinaryOperation(%v, %v, %v) did not return error", left, op, right)
	}
}

func TestExpressionString(t *testing.T) {
	constant := NewConstant(types.Dec("123"))
	columnReference := NewColumnReference(1, types.TypeDecimal)
	binaryOperation, err := NewBinaryOperation(constant, BinaryOperatorEq, columnReference)
	if err != nil {
		t.Fatalf("NewBinaryOperation returned an error: %v", err)
	}
	cases := []struct {
		e    Expression
		want string
	}{
		{constant, "Constant(123)"},
		{columnReference, "ColumnReference(1, decimal)"},
		{binaryOperation, "BinaryOperation(Constant(123) eq ColumnReference(1, decimal))"},
	}
	for _, c := range cases {
		got := c.e.String()
		if got != c.want {
			t.Errorf("String() returned %q, want %q", got, c.want)
		}
	}
}
