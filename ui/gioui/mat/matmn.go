// This file is generated from mgl32/matmn.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"math"
)

// MatMxN is an arbitrary mxn matrix backed by a slice of floats.
//
// This is emphatically not recommended for hardcore n-dimensional
// linear algebra. For that purpose I recommend github.com/gonum/matrix or
// well-tested C libraries such as BLAS or LAPACK.
//
// This is meant to complement future algorithms that may require matrices larger than
// 4x4, but still relatively small (e.g. Jacobeans for inverse kinematics).
//
// It makes use of the same memory sync.Pool set that VecN does, with the same sizing rules.
//
// MatMN will always check if the receiver is nil on any method. Meaning MathMN(nil).Add(dst,m2)
// should always work. Except for the Reshape function, the semantics of this is to "propogate" nils
// forward, so if an invalid operation occurs in a long chain of matrix operations, the overall result will be nil.
type MatMxN[T Float] struct {
	m, n int
	dat  []T
}

// NewMatrix creates a matrix backed by a new slice of size m*n
func NewMatrix[T Float](m, n int) (mat *MatMxN[T]) {
	return &MatMxN[T]{m: m, n: n, dat: make([]T, m*n)}
}

// NewMatrixFromData returns a matrix with data specified by the data in src
//
// For instance, to create a 3x3 MatMN from a Mat3
//
//    m1 := mgl32.Rotate3DX(3.14159)
//    mat := mgl32.NewBackedMatrix(m1[:],3,3)
//
// will create an MN matrix matching the data in the original rotation matrix.
// This matrix is NOT backed by the initial slice; it's a copy of the data
//
// If m*n > cap(src), this function will panic.
func NewMatrixFromData[T Float](src []T, m, n int) *MatMxN[T] {
	internal := make([]T, m*n)
	copy(internal, src[:m*n])
	return &MatMxN[T]{m: m, n: n, dat: internal}
}

// CopyMatMN copies src into dst. This Reshapes dst
// to the same size as src.
//
// If dst or src is nil, this is a no-op
func CopyMatMN[T Float](dst, src *MatMxN[T]) {
	if dst == nil || src == nil {
		return
	}
	dst.Reshape(src.m, src.n)
	copy(dst.dat, src.dat)
}

// IdentN stores the NxN identity matrix in dst, reallocating as necessary.
func IdentN[T Float](dst *MatMxN[T], n int) *MatMxN[T] {
	dst = dst.Reshape(n, n)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				dst.Set(i, j, 1)
			} else {
				dst.Set(i, j, 0)
			}
		}
	}

	return dst
}

// DiagN creates an NxN diagonal matrix seeded by the diagonal vector diag.
// Meaning: for all entries, where i==j, dst.At(i,j) = diag[i]. Otherwise
// dst.At(i,j) = 0
//
// This reshapes dst to the correct size, returning/grabbing from the memory
// pool as necessary.
func DiagN[T Float](dst *MatMxN[T], diag *VecN[T]) *MatMxN[T] {
	dst = dst.Reshape(len(diag.vec), len(diag.vec))
	n := len(diag.vec)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				dst.Set(i, j, diag.vec[i])
			} else {
				dst.Set(i, j, 0)
			}
		}
	}

	return dst
}

// Zero reshapes the matrix to m by n and zeroes out all
// elements.
func (mat *MatMxN[T]) Zero(m, n int) {
	if mat == nil {
		return
	}

	mat.Reshape(m, n)
	for i := range mat.dat {
		mat.dat[i] = 0
	}
}

// destroy returns the underlying matrix slice to the memory pool
func (mat *MatMxN[T]) destroy() {
	if mat == nil {
		return
	}
	mat.m, mat.n = 0, 0
	mat.dat = nil
}

// Reshape reshapes the matrix to the desired dimensions.
// If the overall size of the new matrix (m*n) is bigger
// than the current size, the underlying slice will
// be grown, sending the current slice to the memory pool
// and grabbing a bigger one if necessary
//
// If the caller is a nil pointer, the return value will be a new
// matrix, as if NewMatrix(m,n) had been called. Otherwise it's
// simply the caller.
func (mat *MatMxN[T]) Reshape(m, n int) *MatMxN[T] {
	if mat == nil {
		return NewMatrix[T](m, n)
	}
	if m*n <= cap(mat.dat) {
		if mat.dat != nil {
			mat.dat = mat.dat[:m*n]
		} else {
			mat.dat = []T{}
		}
		mat.m, mat.n = m, n
		return mat
	}
	(*mat) = (*NewMatrix[T](m, n))
	return mat
}

// InferMatrix infers an MxN matrix from a constant matrix from this package.
// For instance, a Mat2x3 inferred with this function will work just like
// NewMatrixFromData(m[:],2,3) where m is the Mat2x3. This uses a type switch.
//
// I personally recommend using NewMatrixFromData, because it avoids a
// potentially costly type switch. However, this is also more robust and less
// error prone if you change the size of your matrix somewhere.
//
// If the value passed in is not recognized, it returns an InferMatrixError.
func (mat *MatMxN[T]) InferMatrix(m interface{}) (*MatMxN[T], error) {
	switch raw := m.(type) {
	case Mat2[T]:
		return NewMatrixFromData(raw[:], 2, 2), nil
	case Mat2x3[T]:
		return NewMatrixFromData(raw[:], 2, 3), nil
	case Mat2x4[T]:
		return NewMatrixFromData(raw[:], 2, 4), nil
	case Mat3[T]:
		return NewMatrixFromData(raw[:], 3, 3), nil
	case Mat3x2[T]:
		return NewMatrixFromData(raw[:], 3, 2), nil
	case Mat3x4[T]:
		return NewMatrixFromData(raw[:], 3, 4), nil
	case Mat4[T]:
		return NewMatrixFromData(raw[:], 4, 4), nil
	case Mat4x2[T]:
		return NewMatrixFromData(raw[:], 4, 2), nil
	case Mat4x3[T]:
		return NewMatrixFromData(raw[:], 4, 3), nil
	default:
		return nil, InferMatrixError{}
	}
}

// Trace returns the trace of a square matrix (sum of all diagonal elements). If
// the matrix is nil, or not square, the result will be NaN.
func (mat *MatMxN[T]) Trace() T {
	if mat == nil || mat.m != mat.n {
		return T(math.NaN())
	}

	var out T
	for i := 0; i < mat.m; i++ {
		out += mat.At(i, i)
	}

	return out
}

// Transpose takes the transpose of mat and puts it in dst.
//
// If dst is not of the correct dimensions, it will be Reshaped,
// if dst and mat are the same, a temporary matrix of the correct size will
// be allocated; these resources will be released via the memory pool.
//
// This should be improved in the future.
func (mat *MatMxN[T]) Transpose(dst *MatMxN[T]) (t *MatMxN[T]) {
	if mat == nil {
		return nil
	}

	if dst == mat {
		dst = NewMatrix[T](mat.n, mat.m)

		// Copy data to correct matrix,
		// delete temporary buffer,
		// and set the return value to the
		// correct one
		defer func() {
			copy(mat.dat, dst.dat)

			mat.m, mat.n = mat.n, mat.m

			dst.destroy()
			t = mat
		}()

		return mat
	} else {
		dst = dst.Reshape(mat.n, mat.m)
	}

	for r := 0; r < mat.m; r++ {
		for c := 0; c < mat.n; c++ {
			dst.dat[r*dst.m+c] = mat.dat[c*mat.m+r]
		}
	}

	return dst
}

// Raw returns the raw slice backing this matrix
func (mat *MatMxN[T]) Raw() []T {
	if mat == nil {
		return nil
	}
	return mat.dat
}

// NumRows returns the number of rows in this matrix
func (mat *MatMxN[T]) NumRows() int {
	return mat.m
}

// NumCols returns the number of columns in this matrix
func (mat *MatMxN[T]) NumCols() int {
	return mat.n
}

// NumRowCols returns the number of rows and columns in this matrix
// as a single operation
func (mat *MatMxN[T]) NumRowCols() (rows, cols int) {
	return mat.m, mat.n
}

// At returns the element at the given row and column.
// This is garbage in/garbage out and does no bounds
// checking. If the computation happens to lead to an invalid
// element, it will be returned; or it may panic.
func (mat *MatMxN[T]) At(row, col int) T {
	return mat.dat[col*mat.m+row]
}

// Set sets the element at the given row and column.
// This is garbage in/garbage out and does no bounds
// checking. If the computation happens to lead to an invalid
// element, it will be set; or it may panic.
func (mat *MatMxN[T]) Set(row, col int, val T) {
	mat.dat[col*mat.m+row] = val
}

// Add is the arithemtic + operator defined on a MatMxN.
func (mat *MatMxN[T]) Add(dst *MatMxN[T], addend *MatMxN[T]) *MatMxN[T] {
	if mat == nil || addend == nil || mat.m != addend.m || mat.n != addend.n {
		return nil
	}

	dst = dst.Reshape(mat.m, mat.n)

	// No need to care about rows and columns
	// since it's element-wise anyway
	for i, el := range mat.dat {
		dst.dat[i] = el + addend.dat[i]
	}

	return dst
}

// Sub is the arithemtic - operator defined on a MatMxN.
func (mat *MatMxN[T]) Sub(dst *MatMxN[T], subtrahend *MatMxN[T]) *MatMxN[T] {
	if mat == nil || subtrahend == nil || mat.m != subtrahend.m || mat.n != subtrahend.n {
		return nil
	}

	dst = dst.Reshape(mat.m, mat.n)

	// No need to care about rows and columns
	// since it's element-wise anyway
	for i, el := range mat.dat {
		dst.dat[i] = el - subtrahend.dat[i]
	}

	return dst
}

// MulMxN performs matrix multiplication on MxN matrix mat and NxO matrix mul,
// storing the result in dst. This returns dst, or nil if the operation is not
// able to be performed.
//
// If mat == dst, or mul == dst a temporary matrix will be used.
//
// This uses the naive algorithm (though on smaller matrices, this can actually
// be faster; about len(mat)+len(mul) < ~100)
func (mat *MatMxN[T]) MulMxN(dst *MatMxN[T], mul *MatMxN[T]) *MatMxN[T] {
	if mat == nil || mul == nil || mat.n != mul.m {
		return nil
	}

	if dst == mul {
		mul = NewMatrix[T](mul.m, mul.n)
		copy(mul.dat, dst.dat)

		// If mat==dst==mul, we need to change
		// mat too or we have a bug
		if mat == dst {
			mat = mul
		}

		defer mul.destroy()
	} else if dst == mat {
		mat = NewMatrix[T](mat.m, mat.n)
		copy(mat.dat, dst.dat)

		defer mat.destroy()
	}

	dst = dst.Reshape(mat.m, mul.n)
	for r1 := 0; r1 < mat.m; r1++ {
		for c2 := 0; c2 < mul.n; c2++ {

			dst.dat[c2*mat.m+r1] = 0
			for i := 0; i < mat.n; i++ {
				dst.dat[c2*mat.m+r1] += mat.dat[i*mat.m+r1] * mul.dat[c2*mul.m+i]
			}

		}
	}

	return dst
}

// Mul performs a scalar multiplication between mat and some constant c,
// storing the result in dst. Mat and dst can be equal. If dst is not the
// correct size, a Reshape will occur.
func (mat *MatMxN[T]) Mul(dst *MatMxN[T], c T) *MatMxN[T] {
	if mat == nil {
		return nil
	}

	dst = dst.Reshape(mat.m, mat.n)

	for i, el := range mat.dat {
		dst.dat[i] = el * c
	}

	return dst
}

// MulNx1 multiplies the matrix by a vector of size n. If mat or v is nil, this
// returns nil. If the number of columns in mat does not match the Size of v,
// this also returns nil.
//
// Dst will be resized if it's not big enough. If dst == v; a temporary vector
// will be allocated and returned via the realloc callback when complete.
func (mat *MatMxN[T]) MulNx1(dst, v *VecN[T]) *VecN[T] {
	if mat == nil || v == nil || mat.n != len(v.vec) {
		return nil
	}
	if dst == v {
		v = NewVecN[T](len(v.vec))
		copy(v.vec, dst.vec)
		defer v.destroy()
	}
	dst = dst.Resize(mat.m)
	for r := 0; r < mat.m; r++ {
		dst.vec[r] = 0
		for c := 0; c < mat.n; c++ {
			dst.vec[r] += mat.At(r, c) * v.vec[c]
		}
	}
	return dst
}

// ApproxEqual returns whether the two vectors are approximately equal (See
// FloatEqual).
func (mat *MatMxN[T]) ApproxEqual(m2 *MatMxN[T]) bool {
	if mat == m2 {
		return true
	}
	if mat.m != m2.m || mat.n != m2.n {
		return false
	}

	for i, el := range mat.dat {
		if !FloatEqual(el, m2.dat[i]) {
			return false
		}
	}

	return true
}

// ApproxEqualThreshold returns whether the two vectors are approximately equal
// to within the given threshold given by "epsilon" (See ApproxEqualThreshold).
func (mat *MatMxN[T]) ApproxEqualThreshold(m2 *MatMxN[T], epsilon T) bool {
	if mat == m2 {
		return true
	}
	if mat.m != m2.m || mat.n != m2.n {
		return false
	}

	for i, el := range mat.dat {
		if !FloatEqualThreshold(el, m2.dat[i], epsilon) {
			return false
		}
	}

	return true
}

// ApproxEqualFunc returns whether the two vectors are approximately equal,
// given a function which compares two scalar elements.
func (mat *MatMxN[T]) ApproxEqualFunc(m2 *MatMxN[T], comp func(T, T) bool) bool {
	if mat == m2 {
		return true
	}
	if mat.m != m2.m || mat.n != m2.n {
		return false
	}

	for i, el := range mat.dat {
		if !comp(el, m2.dat[i]) {
			return false
		}
	}

	return true
}

// InferMatrixError may be returned by InferMatrix.
//
// Make sure you're using a constant matrix such as Mat3 from within the same
// package (meaning: mgl32.MatMxN can't handle a mgl64.Mat2x3).
type InferMatrixError struct{}

func (me InferMatrixError) Error() string {
	return "could not infer matrix. Make sure you're using a constant matrix such as Mat3 from within the same package (meaning: mgl32.MatMxN can't handle a mgl64.Mat2x3)."
}

// RectangularMatrixError is returned when a rectangular matrix was expected but
// not given.
type RectangularMatrixError struct{}

func (mse RectangularMatrixError) Error() string {
	return "the matrix was the wrong shape, needed a square matrix."
}

// NilMatrixError is returned when an operand to a function was unexpectedly 'nil'.
type NilMatrixError struct{}

func (me NilMatrixError) Error() string {
	return "the matrix is nil"
}
