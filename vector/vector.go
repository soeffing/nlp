package vector

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

var (
	// ErrVector dimension error
	ErrVector = errors.New("dimension not equal")
	// ErrZero for zero vectors
	ErrZero = errors.New("zero vector")
	// ErrNotSupported for not supported types
	ErrNotSupported = errors.New("not support")
)

// Vector type
type Vector []float64

// NewVector creates a new vector
// Create a new vector who's dimension is dim, and components are initialized by c,
// if len(c) > dim, the redundant components of c are discarded,
// if len(c) < dim, the remain component values of new vector are 0
func NewVector(dim int, c ...float64) Vector {
	if dim <= 0 {
		return nil
	}
	v := make([]float64, dim)
	copy(v, c)

	return v
}

// Copy copies a vector
func Copy(v Vector) Vector {
	if v == nil {
		return v
	}
	return v.Fork()
}

// Dim defines vector dimension
func Dim(v Vector) int {
	return len(v)
}

// DimEqual check if vectors are equal
func DimEqual(v1, v2 Vector) bool {
	return (v1 != nil && v2 != nil) && Dim(v1) == Dim(v2)
}

// Add two vectors and return a new vector as the result
func Add(v1, v2 Vector) (Vector, error) {
	if !DimEqual(v1, v2) {
		return nil, ErrVector
	}

	v := NewVector(Dim(v1))
	for i := range v1 {
		v[i] = v1[i] + v2[i]
	}
	return v, nil
}

//Sub subtracts v2 from v1 and return a new vector as the result
func Sub(v1, v2 Vector) (Vector, error) {
	if !DimEqual(v1, v2) {
		return nil, ErrVector
	}

	v := NewVector(Dim(v1))

	for i := range v1 {
		v[i] = v1[i] - v2[i]
	}

	return v, nil
}

// Multi multiplies v by scalar a and return a new vector
func Multi(v Vector, a float64) Vector {
	w := NewVector(Dim(v))
	for i := range v {
		w[i] = v[i] * a
	}
	return w
}

// Div divides v by scalar a, return a new vector as the result
func Div(v Vector, a float64) Vector {
	if a == 0 {
		return v
	}
	w := NewVector(Dim(v))
	f := 1.0 / a
	for i := range v {
		w[i] = v[i] * f
	}
	return w
}

// Dot product (or inner product) of the two vectors v1 and v2
func Dot(v1, v2 Vector) (float64, error) {
	var dot float64

	if !DimEqual(v1, v2) {
		return 0, ErrVector
	}

	for i := range v1 {
		dot += v1[i] * v2[i]
	}

	return dot, nil
}

// Cross product (or vector product) of the two 3D vectors v1 and v2,
// v1×v2 = ||v1|| ||v2||sin(θ)N,
// θ is angle between v1 and v2, N is the unit vector that is perpendicular to both v1 and v2
func Cross(v1, v2 Vector) (Vector, error) {
	if !DimEqual(v1, v2) {
		return nil, ErrVector
	}
	if Dim(v1) != 3 {
		return nil, ErrNotSupported
	}

	x, y, z := 0, 1, 2
	cross := Vector{
		v1[y]*v2[z] - v1[z]*v2[y],
		v1[z]*v2[x] - v1[x]*v2[z],
		v1[x]*v2[y] - v1[y]*v2[x],
	}
	return cross, nil
}

// Reflect ...
// return the reflect vector of vector v, n is a normal to a surface,
// R = v - (2n·v)n
func Reflect(v, n Vector) (Vector, error) {
	if !DimEqual(v, n) {
		return nil, ErrVector
	}
	n = Normalize(n)
	v = Normalize(v)

	d, _ := Dot(n, v)
	return Sub(v, Multi(n, d*2))
}

// Refract ...
func Refract(v, n Vector, eta float64) (Vector, error) {
	if !DimEqual(v, n) {
		return nil, ErrVector
	}
	n = Normalize(n)
	v = Normalize(v)

	d, _ := Dot(n, v)
	k := 1 - eta*eta*(1-d*d)
	if k < 0.0 {
		return NewVector(v.Dim()), nil
	}

	return Sub(Multi(v, eta), Multi(n, eta*d+math.Sqrt(k)))
}

// Distance between two vectors
func Distance(v1, v2 Vector) (float64, error) {
	v, err := Sub(v1, v2)
	if err != nil {
		return 0, err
	}
	return v.Length(), nil
}

// IsZero check
func IsZero(v Vector) bool {
	return v == nil || v.IsZero()
}

// Normalize a vector
func Normalize(v Vector) Vector {
	if IsZero(v) {
		return v
	}
	return Div(v, v.Length())
}

// Angle between two vectors
func Angle(v1, v2 Vector) (float64, error) {
	d, err := Dot(Normalize(v1), Normalize(v2))
	if err != nil {
		return 0, err
	}
	return math.Acos(d) * 180.0 / math.Pi, nil
}

// Init initializes vectors
func (v Vector) Init(c ...float64) {
	copy(v, c)
}

// InitV single vectors?
func (v Vector) InitV(vec Vector) {
	copy(v, vec)
}

// Set value in index, return old value
func (v Vector) Set(index int, value float64) float64 {
	if index <= 0 || index > v.Dim() {
		return 0
	}
	old := v[index-1]
	v[index-1] = value
	return old
}

// Get value by index
func (v Vector) Get(index int) float64 {
	if index > 0 && index <= v.Dim() {
		return v[index-1]
	}
	return 0
}

// Dim of vectors
func (v Vector) Dim() int {
	return len(v)
}

// Equal wether two vectors are equal, means all components are equal
func (v Vector) Equal(vec Vector) bool {
	if vec == nil || v.Dim() != vec.Dim() {
		return false
	}

	equal := true
	for i := range v {
		if v[i] != vec[i] {
			equal = false
			break
		}
	}

	return equal
}

// Zero sets vector to zero vector(all components are zeros)
func (v Vector) Zero() Vector {
	for i := range v {
		v[i] = 0
	}
	return v
}

// IsZero check
func (v Vector) IsZero() bool {
	b := true

	for i := range v {
		if v[i] != 0 {
			b = false
			break
		}
	}

	return b
}

// Negate vectors check
func (v Vector) Negate() Vector {
	for i := range v {
		v[i] = -v[i]
	}
	return v
}

// Add c...
func (v Vector) Add(vec Vector) Vector {
	if vec == nil || v.Dim() != vec.Dim() {
		return v
	}

	for i := range v {
		v[i] += vec[i]
	}

	return v
}

// Sub c...
func (v Vector) Sub(vec Vector) Vector {
	if vec == nil || v.Dim() != vec.Dim() {
		return v
	}

	for i := range v {
		v[i] -= vec[i]
	}

	return v
}

// Multi c...
func (v Vector) Multi(a float64) Vector {
	for i := range v {
		v[i] *= a
	}
	return v
}

// Div c...
func (v Vector) Div(a float64) Vector {
	if a == 0 {
		return v
	}
	f := 1.0 / a
	for i := range v {
		v[i] *= f
	}

	return v
}

// Normalize c...
// Unit vector
func (v Vector) Normalize() Vector {
	if !v.IsZero() {
		v.Div(v.Length())
	}

	return v
}

// Length c...
// The length (or magnitude) of the vector
// ||v||
func (v Vector) Length() float64 {
	d, _ := Dot(v, v)
	return math.Sqrt(d)
}

// Fork c...
func (v Vector) Fork() Vector {
	w := make([]float64, v.Dim())
	copy(w, v)
	return w
}

// Parallel c...
// v parallel vector in vec
func (v Vector) Parallel(vec Vector) (Vector, error) {
	if !DimEqual(v, vec) {
		return nil, ErrVector
	}
	norm := Normalize(vec)
	dot, _ := Dot(norm, v)
	return Multi(norm, dot), nil
}

// Perp c...ar vector in vec
// v perpendicul
func (v Vector) Perp(vec Vector) (Vector, error) {
	if !DimEqual(v, vec) {
		return nil, ErrVector
	}
	w, _ := v.Parallel(vec)

	return Sub(v, w)
}

// ToSlice c...
func (v Vector) ToSlice() []float64 {
	s := make([]float64, v.Dim())
	copy(s, v)
	return s
}

// ToSlice32 c...
func (v Vector) ToSlice32() []float32 {
	s := make([]float32, v.Dim())
	for i, value := range v {
		s[i] = float32(value)
	}
	return s
}

// String c...
func (v Vector) String() string {
	buf := new(bytes.Buffer)

	for i := 1; i <= v.Dim(); i++ {
		fmt.Fprintf(buf, " %.2f", v.Get(i))
	}

	return buf.String()
}
