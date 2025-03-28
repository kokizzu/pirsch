// Code generated by ./cmd/ch-gen-col, DO NOT EDIT.

package proto

// ColIPv6 represents IPv6 column.
type ColIPv6 []IPv6

// Compile-time assertions for ColIPv6.
var (
	_ ColInput  = ColIPv6{}
	_ ColResult = (*ColIPv6)(nil)
	_ Column    = (*ColIPv6)(nil)
)

// Rows returns count of rows in column.
func (c ColIPv6) Rows() int {
	return len(c)
}

// Reset resets data in row, preserving capacity for efficiency.
func (c *ColIPv6) Reset() {
	*c = (*c)[:0]
}

// Type returns ColumnType of IPv6.
func (ColIPv6) Type() ColumnType {
	return ColumnTypeIPv6
}

// Row returns i-th row of column.
func (c ColIPv6) Row(i int) IPv6 {
	return c[i]
}

// Append IPv6 to column.
func (c *ColIPv6) Append(v IPv6) {
	*c = append(*c, v)
}

// Append IPv6 slice to column.
func (c *ColIPv6) AppendArr(vs []IPv6) {
	*c = append(*c, vs...)
}

// LowCardinality returns LowCardinality for IPv6.
func (c *ColIPv6) LowCardinality() *ColLowCardinality[IPv6] {
	return &ColLowCardinality[IPv6]{
		index: c,
	}
}

// Array is helper that creates Array of IPv6.
func (c *ColIPv6) Array() *ColArr[IPv6] {
	return &ColArr[IPv6]{
		Data: c,
	}
}

// Nullable is helper that creates Nullable(IPv6).
func (c *ColIPv6) Nullable() *ColNullable[IPv6] {
	return &ColNullable[IPv6]{
		Values: c,
	}
}

// NewArrIPv6 returns new Array(IPv6).
func NewArrIPv6() *ColArr[IPv6] {
	return &ColArr[IPv6]{
		Data: new(ColIPv6),
	}
}
