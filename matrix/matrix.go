// Package matrix ...
package matrix

import (
	"bytes"
	"fmt"
	"github.com/ginuerzh/math3d/vector"
)

// Matrix is a bunch of vectors
type Matrix []vector.Vector

// NewMatrix ...
func NewMatrix(row, col int) Matrix {
	m := make([]vector.Vector, col)
	for i := range m {
		m[i] = vector.NewVector(row)
	}
	return m
}

// NewIdentityMatrix ...
func NewIdentityMatrix(row int) Matrix {
	m := NewMatrix(row, row)
	for i := 1; i <= m.Cols(); i++ {
		m.Column(i).Set(i, 1)
	}
	return m
}

// MultiSM ...trix m multiply by scale, return a new matrix
// ma
func MultiSM(scale float64, m Matrix) Matrix {
	return m.Fork().MultiS(scale)
}

// MultiMM ...
func MultiMM(m Matrix, ms ...Matrix) Matrix {
	nm := m
	for _, mm := range ms {
		nm = nm.multiM(mm)
	}

	return nm
}

// Init ...
func (m Matrix) Init(cols ...vector.Vector) {
	for i := range cols {
		copy(m.Column(i+1), cols[i])
	}
}

// InitColumn ...
func (m Matrix) InitColumn(col int, v vector.Vector) {
	if col <= 0 || col > m.Cols() {
		return
	}

	m.Column(col).InitV(v)
}

// InitRow ...
func (m Matrix) InitRow(row int, v vector.Vector) {
	if row <= 0 || row > m.Rows() {
		return
	}

	for i := 1; i <= m.Cols(); i++ {
		if v.Dim() < i {
			break
		}
		m.Set(row, i, v.Get(i))
	}
}

// Get ...
func (m Matrix) Get(row, col int) float64 {
	if row <= 0 || row > m.Rows() || col <= 0 || col > m.Cols() {
		return 0
	}
	return m.Column(col).Get(row)
}

// Set ...[row, col], return old value
// set value in
func (m Matrix) Set(row, col int, value float64) float64 {
	if row <= 0 || row > m.Rows() || col <= 0 || col > m.Cols() {
		return 0
	}

	return m.Column(col).Set(row, value)
}

// Cols ...
func (m Matrix) Cols() int {
	return len(m)
}

// Rows ...
func (m Matrix) Rows() int {
	return len(m[0])
}

// Column ...
// return column col in matrix, col starts from 1,
// NOTE: the returned vector is a reference to the row vector in matrix m
func (m Matrix) Column(col int) vector.Vector {
	if col > 0 && col <= len(m) {
		return m[col-1]
	}
	return nil
}

// Row ...urned vector is a new vector, not a reference to the vector in matrix
// return row 'row' in matrix, row starts from 1,
// NOTE: the ret
func (m Matrix) Row(row int) vector.Vector {
	if row > 0 && row <= m.Rows() {
		v := vector.NewVector(m.Cols())
		for i := 1; i <= m.Cols(); i++ {
			v.Set(i, m.Column(i).Get(row))
		}
		return v
	}

	return nil
}

// Transpose ...
func (m Matrix) Transpose() Matrix {
	tran := NewMatrix(m.Cols(), m.Rows())

	for i := 1; i <= tran.Cols(); i++ {
		tran.InitColumn(i, m.Row(i))
	}
	return tran
}

// MultiS ...
func (m Matrix) MultiS(scale float64) Matrix {
	for i := 1; i <= m.Cols(); i++ {
		m.Column(i).Multi(scale)
	}
	return m
}

// multiM ...matrix
// return a new
func (m Matrix) multiM(m2 Matrix) Matrix {
	if m.Cols() != m2.Rows() {
		return nil
	}
	mm := NewMatrix(m.Rows(), m2.Cols())
	for i := 1; i <= mm.Cols(); i++ {
		col := mm.Column(i)
		for j := 1; j <= col.Dim(); j++ {
			dot, _ := vector.Dot(m.Row(j), m2.Column(i))
			col.Set(j, dot)
		}
	}

	return mm
}

// Fork ...
func (m Matrix) Fork() Matrix {
	mx := NewMatrix(m.Rows(), m.Cols())
	mx.Init(m...)
	return mx
}

// ToSlice ...
func (m Matrix) ToSlice() []float64 {
	s := make([]float64, 0, m.Rows()*m.Cols())

	for i := 1; i <= m.Cols(); i++ {
		s = append(s, m.Column(i)...)
	}

	return s
}

// ToSlice32 ...
func (m Matrix) ToSlice32() []float32 {
	s := make([]float32, 0, m.Rows()*m.Cols())

	for i := 1; i <= m.Cols(); i++ {
		s = append(s, m.Column(i).ToSlice32()...)
	}

	return s
}

// String ...
func (m Matrix) String() string {
	buf := new(bytes.Buffer)

	fmt.Fprintln(buf)
	for i := 1; i <= m.Rows(); i++ {
		fmt.Fprintln(buf, m.Row(i))
	}

	return buf.String()
}
