// Code generated by ./cmd/ch-gen-col, DO NOT EDIT.

package proto

// ColInt32 represents Int32 column.
type ColInt32 []int32

// Compile-time assertions for ColInt32.
var (
	_ ColInput  = ColInt32{}
	_ ColResult = (*ColInt32)(nil)
	_ Column    = (*ColInt32)(nil)
)

// Rows returns count of rows in column.
func (c ColInt32) Rows() int {
	return len(c)
}

// Reset resets data in row, preserving capacity for efficiency.
func (c *ColInt32) Reset() {
	*c = (*c)[:0]
}

// Type returns ColumnType of Int32.
func (ColInt32) Type() ColumnType {
	return ColumnTypeInt32
}

// Row returns i-th row of column.
func (c ColInt32) Row(i int) int32 {
	return c[i]
}

// Append int32 to column.
func (c *ColInt32) Append(v int32) {
	*c = append(*c, v)
}

// Append int32 slice to column.
func (c *ColInt32) AppendArr(vs []int32) {
	*c = append(*c, vs...)
}

// LowCardinality returns LowCardinality for Int32.
func (c *ColInt32) LowCardinality() *ColLowCardinality[int32] {
	return &ColLowCardinality[int32]{
		index: c,
	}
}

// Array is helper that creates Array of int32.
func (c *ColInt32) Array() *ColArr[int32] {
	return &ColArr[int32]{
		Data: c,
	}
}

// Nullable is helper that creates Nullable(int32).
func (c *ColInt32) Nullable() *ColNullable[int32] {
	return &ColNullable[int32]{
		Values: c,
	}
}

// NewArrInt32 returns new Array(Int32).
func NewArrInt32() *ColArr[int32] {
	return &ColArr[int32]{
		Data: new(ColInt32),
	}
}