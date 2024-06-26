// This file is generated from mgl32/vector.go; DO NOT EDIT

// Copyright 2014 The go-gl/mathgl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is generated by codegen.go; DO NOT EDIT
// Edit vector.tmpl and run "go generate" to make changes.

package mat

import (
	"math"
)

type Vec2[T Float] [2]T
type Vec3[T Float] [3]T
type Vec4[T Float] [4]T

// Vec3 constructs a 3-dimensional vector by appending the given coordinates.
func (v Vec2[T]) Vec3(z T) Vec3[T] {
	return Vec3[T]{v[0], v[1], z}
}

// Vec4 constructs a 4-dimensional vector by appending the given coordinates.
func (v Vec2[T]) Vec4(z, w T) Vec4[T] {
	return Vec4[T]{v[0], v[1], z, w}
}

// Vec4 constructs a 4-dimensional vector by appending the given coordinates.
func (v Vec3[T]) Vec4(w T) Vec4[T] {
	return Vec4[T]{v[0], v[1], v[2], w}
}

// Vec2 constructs a 2-dimensional vector by discarding coordinates.
func (v Vec3[T]) Vec2() Vec2[T] {
	return Vec2[T]{v[0], v[1]}
}

// Vec2 constructs a 2-dimensional vector by discarding coordinates.
func (v Vec4[T]) Vec2() Vec2[T] {
	return Vec2[T]{v[0], v[1]}
}

// Vec3 constructs a 3-dimensional vector by discarding coordinates.
func (v Vec4[T]) Vec3() Vec3[T] {
	return Vec3[T]{v[0], v[1], v[2]}
}

// Elem extracts the elements of the vector for direct value assignment.
func (v Vec2[T]) Elem() (x, y T) {
	return v[0], v[1]
}

// Elem extracts the elements of the vector for direct value assignment.
func (v Vec3[T]) Elem() (x, y, z T) {
	return v[0], v[1], v[2]
}

// Elem extracts the elements of the vector for direct value assignment.
func (v Vec4[T]) Elem() (x, y, z, w T) {
	return v[0], v[1], v[2], v[3]
}

// Cross is the vector cross product. This operation is only defined on 3D
// vectors. It is equivalent to Vec3[T]{v1[1]*v2[2]-v1[2]*v2[1],
// v1[2]*v2[0]-v1[0]*v2[2], v1[0]*v2[1] - v1[1]*v2[0]}. Another interpretation
// is that it's the vector whose magnitude is |v1||v2|sin(theta) where theta is
// the angle between v1 and v2.
//
// The cross product is most often used for finding surface normals. The cross
// product of vectors will generate a vector that is perpendicular to the plane
// they form.
//
// Technically, a generalized cross product exists as an "(N-1)ary" operation
// (that is, the 4D cross product requires 3 4D vectors). But the binary 3D (and
// 7D) cross product is the most important. It can be considered the area of a
// parallelogram with sides v1 and v2.
//
// Like the dot product, the cross product is roughly a measure of
// directionality. Two normalized perpendicular vectors will return a vector
// with a magnitude of 1.0 or -1.0 and two parallel vectors will return a vector
// with magnitude 0.0. The cross product is "anticommutative" meaning
// v1.Cross(v2) = -v2.Cross(v1), this property can be useful to know when
// finding normals, as taking the wrong cross product can lead to the opposite
// normal of the one you want.
func (v1 Vec3[T]) Cross(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{v1[1]*v2[2] - v1[2]*v2[1], v1[2]*v2[0] - v1[0]*v2[2], v1[0]*v2[1] - v1[1]*v2[0]}
}

// Quat reinterprets this vector as a quaternion, with the individual elements
// staying the same.
func (v Vec4[T]) Quat() Quat[T] {
	return Quat[T]{v[3], Vec3[T]{v[0], v[1], v[2]}}
}

// Add performs element-wise addition between two vectors. It is equivalent to iterating
// over every element of v1 and adding the corresponding element of v2 to it.
func (v1 Vec2[T]) Add(v2 Vec2[T]) Vec2[T] {
	return Vec2[T]{v1[0] + v2[0], v1[1] + v2[1]}
}

// Sub performs element-wise subtraction between two vectors. It is equivalent to iterating
// over every element of v1 and subtracting the corresponding element of v2 from it.
func (v1 Vec2[T]) Sub(v2 Vec2[T]) Vec2[T] {
	return Vec2[T]{v1[0] - v2[0], v1[1] - v2[1]}
}

// Mul performs a scalar multiplication between the vector and some constant value
// c. This is equivalent to iterating over every vector element and multiplying by c.
func (v1 Vec2[T]) Mul(c T) Vec2[T] {
	return Vec2[T]{v1[0] * c, v1[1] * c}
}

// Dot returns the dot product of this vector with another. There are multiple ways
// to describe this value. One is the multiplication of their lengths and cos(theta) where
// theta is the angle between the vectors: v1.v2 = |v1||v2|cos(theta).
//
// The other (and what is actually done) is the sum of the element-wise multiplication of all
// elements. So for instance, two Vec3s would yield v1.x * v2.x + v1.y * v2.y + v1.z * v2.z.
//
// This means that the dot product of a vector and itself is the square of its Len (within
// the bounds of floating points error).
//
// The dot product is roughly a measure of how closely two vectors are to pointing in the same
// direction. If both vectors are normalized, the value will be -1 for opposite pointing,
// one for same pointing, and 0 for perpendicular vectors.
func (v1 Vec2[T]) Dot(v2 Vec2[T]) T {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

// Len returns the vector's length. Note that this is NOT the dimension of
// the vector (len(v)), but the mathematical length. This is equivalent to the square
// root of the sum of the squares of all elements. E.G. for a Vec2 it's
// math.Hypot(v[0], v[1]).
func (v1 Vec2[T]) Len() T {
	return T(math.Hypot(float64(v1[0]), float64(v1[1])))
}

// LenSqr returns the vector's square length. This is equivalent to the sum of the squares of all elements.
func (v1 Vec2[T]) LenSqr() T {
	return v1[0]*v1[0] + v1[1]*v1[1]
}

// Normalize normalizes the vector. Normalization is (1/|v|)*v,
// making this equivalent to v.Scale(1/v.Len()). If the len is 0.0,
// this function will return an infinite value for all elements due
// to how floating point division works in Go (n/0.0 = math.Inf(Sign(n))).
//
// Normalization makes a vector's Len become 1.0 (within the margin of floating point error),
// while maintaining its directionality.
//
// (Can be seen here: http://play.golang.org/p/Aaj7SnbqIp )
func (v1 Vec2[T]) Normalize() Vec2 [T]{
	l :=(T)(1.0 / v1.Len())
	return Vec2[T]{v1[0] * l, v1[1] * l}
}

// ApproxEqual takes in a vector and does an element-wise approximate float
// comparison as if FloatEqual had been used
func (v1 Vec2[T]) ApproxEqual(v2 Vec2[T]) bool {
	for i := range v1 {
		if !FloatEqual(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

// ApproxEqualThreshold takes in a threshold for comparing two floats, and uses
// it to do an element-wise comparison of the vector to another.
func (v1 Vec2[T]) ApproxEqualThreshold(v2 Vec2[T], threshold T) bool {
	for i := range v1 {
		if !FloatEqualThreshold(v1[i], v2[i], threshold) {
			return false
		}
	}
	return true
}

// ApproxFuncEqual takes in a func that compares two floats, and uses it to do an element-wise
// comparison of the vector to another. This is intended to be used with FloatEqualFunc
func (v1 Vec2[T]) ApproxFuncEqual(v2 Vec2[T], eq func(T, T) bool) bool {
	for i := range v1 {
		if !eq(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

// X is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec2[T]) X() T {
	return v[0]
}

// Y is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec2[T]) Y() T {
	return v[1]
}

// OuterProd2 does the vector outer product
// of two vectors. The outer product produces an
// 2x2 matrix. E.G. a Vec2 * Vec2 = Mat2.
//
// The outer product can be thought of as the "opposite"
// of the Dot product. The Dot product treats both vectors like matrices
// oriented such that the left one has N columns and the right has N rows.
// So Vec3.Vec3 = Mat1x3*Mat3x1 = Mat1 = Scalar.
//
// The outer product orients it so they're facing "outward": Vec2*Vec3
// = Mat2x1*Mat1x3 = Mat2x3.
func (v1 Vec2[T]) OuterProd2(v2 Vec2[T]) Mat2[T] {
	return Mat2[T]{v1[0] * v2[0], v1[1] * v2[0], v1[0] * v2[1], v1[1] * v2[1]}
}

// OuterProd3 does the vector outer product
// of two vectors. The outer product produces an
// 2x3 matrix. E.G. a Vec2 * Vec3 = Mat2x3.
//
// The outer product can be thought of as the "opposite"
// of the Dot product. The Dot product treats both vectors like matrices
// oriented such that the left one has N columns and the right has N rows.
// So Vec3.Vec3 = Mat1x3*Mat3x1 = Mat1 = Scalar.
//
// The outer product orients it so they're facing "outward": Vec2*Vec3
// = Mat2x1*Mat1x3 = Mat2x3.
func (v1 Vec2[T]) OuterProd3(v2 Vec3[T]) Mat2x3[T] {
	return Mat2x3[T]{v1[0] * v2[0], v1[1] * v2[0], v1[0] * v2[1], v1[1] * v2[1], v1[0] * v2[2], v1[1] * v2[2]}
}

// OuterProd4 does the vector outer product
// of two vectors. The outer product produces an
// 2x4 matrix. E.G. a Vec2 * Vec4 = Mat2x4.
//
// The outer product can be thought of as the "opposite"
// of the Dot product. The Dot product treats both vectors like matrices
// oriented such that the left one has N columns and the right has N rows.
// So Vec3.Vec3 = Mat1x3*Mat3x1 = Mat1 = Scalar.
//
// The outer product orients it so they're facing "outward": Vec2*Vec3
// = Mat2x1*Mat1x3 = Mat2x3.
func (v1 Vec2[T]) OuterProd4(v2 Vec4[T]) Mat2x4[T] {
	return Mat2x4[T]{v1[0] * v2[0], v1[1] * v2[0], v1[0] * v2[1], v1[1] * v2[1], v1[0] * v2[2], v1[1] * v2[2], v1[0] * v2[3], v1[1] * v2[3]}
}

// Add performs element-wise addition between two vectors. It is equivalent to iterating
// over every element of v1 and adding the corresponding element of v2 to it.
func (v1 Vec3[T]) Add(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

// Sub performs element-wise subtraction between two vectors. It is equivalent to iterating
// over every element of v1 and subtracting the corresponding element of v2 from it.
func (v1 Vec3[T]) Sub(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

// Mul performs a scalar multiplication between the vector and some constant value
// c. This is equivalent to iterating over every vector element and multiplying by c.
func (v1 Vec3[T]) Mul(c T) Vec3[T] {
	return Vec3[T]{v1[0] * c, v1[1] * c, v1[2] * c}
}

// Dot returns the dot product of this vector with another. There are multiple ways
// to describe this value. One is the multiplication of their lengths and cos(theta) where
// theta is the angle between the vectors: v1.v2 = |v1||v2|cos(theta).
//
// The other (and what is actually done) is the sum of the element-wise multiplication of all
// elements. So for instance, two Vec3s would yield v1.x * v2.x + v1.y * v2.y + v1.z * v2.z.
//
// This means that the dot product of a vector and itself is the square of its Len (within
// the bounds of floating points error).
//
// The dot product is roughly a measure of how closely two vectors are to pointing in the same
// direction. If both vectors are normalized, the value will be -1 for opposite pointing,
// one for same pointing, and 0 for perpendicular vectors.
func (v1 Vec3[T]) Dot(v2 Vec3[T]) T {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

// Len returns the vector's length. Note that this is NOT the dimension of
// the vector (len(v)), but the mathematical length. This is equivalent to the square
// root of the sum of the squares of all elements. E.G. for a Vec2 it's
// math.Hypot(v[0], v[1]).
func (v1 Vec3[T]) Len() T {
	return T(math.Sqrt(float64(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])))
}

// LenSqr returns the vector's square length. This is equivalent to the sum of the squares of all elements.
func (v1 Vec3[T]) LenSqr() T {
	return v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2]
}

// Normalize normalizes the vector. Normalization is (1/|v|)*v,
// making this equivalent to v.Scale(1/v.Len()). If the len is 0.0,
// this function will return an infinite value for all elements due
// to how floating point division works in Go (n/0.0 = math.Inf(Sign(n))).
//
// Normalization makes a vector's Len become 1.0 (within the margin of floating point error),
// while maintaining its directionality.
//
// (Can be seen here: http://play.golang.org/p/Aaj7SnbqIp )
func (v1 Vec3[T]) Normalize() Vec3[T] {
	l := (T)(1.0 / v1.Len())
	return Vec3[T]{v1[0] * l, v1[1] * l, v1[2] * l}
}

// ApproxEqual takes in a vector and does an element-wise approximate float
// comparison as if FloatEqual had been used
func (v1 Vec3[T]) ApproxEqual(v2 Vec3[T]) bool {
	for i := range v1 {
		if !FloatEqual(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

// ApproxEqualThreshold takes in a threshold for comparing two floats, and uses
// it to do an element-wise comparison of the vector to another.
func (v1 Vec3[T]) ApproxEqualThreshold(v2 Vec3[T], threshold T) bool {
	for i := range v1 {
		if !FloatEqualThreshold(v1[i], v2[i], threshold) {
			return false
		}
	}
	return true
}

// ApproxFuncEqual takes in a func that compares two floats, and uses it to do an element-wise
// comparison of the vector to another. This is intended to be used with FloatEqualFunc
func (v1 Vec3[T]) ApproxFuncEqual(v2 Vec3[T], eq func(T, T) bool) bool {
	for i := range v1 {
		if !eq(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

// X is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec3[T]) X() T {
	return v[0]
}

// Y is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec3[T]) Y() T {
	return v[1]
}

// Z is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec3[T]) Z() T {
	return v[2]
}

// OuterProd2 does the vector outer product
// of two vectors. The outer product produces an
// 3x2 matrix. E.G. a Vec3 * Vec2 = Mat3x2.
//
// The outer product can be thought of as the "opposite"
// of the Dot product. The Dot product treats both vectors like matrices
// oriented such that the left one has N columns and the right has N rows.
// So Vec3.Vec3 = Mat1x3*Mat3x1 = Mat1 = Scalar.
//
// The outer product orients it so they're facing "outward": Vec2*Vec3
// = Mat2x1*Mat1x3 = Mat2x3.
func (v1 Vec3[T]) OuterProd2(v2 Vec2[T]) Mat3x2[T] {
	return Mat3x2[T]{v1[0] * v2[0], v1[1] * v2[0], v1[2] * v2[0], v1[0] * v2[1], v1[1] * v2[1], v1[2] * v2[1]}
}

// OuterProd3 does the vector outer product
// of two vectors. The outer product produces an
// 3x3 matrix. E.G. a Vec3 * Vec3 = Mat3.
//
// The outer product can be thought of as the "opposite"
// of the Dot product. The Dot product treats both vectors like matrices
// oriented such that the left one has N columns and the right has N rows.
// So Vec3.Vec3 = Mat1x3*Mat3x1 = Mat1 = Scalar.
//
// The outer product orients it so they're facing "outward": Vec2*Vec3
// = Mat2x1*Mat1x3 = Mat2x3.
func (v1 Vec3[T]) OuterProd3(v2 Vec3[T]) Mat3[T] {
	return Mat3[T]{v1[0] * v2[0], v1[1] * v2[0], v1[2] * v2[0], v1[0] * v2[1], v1[1] * v2[1], v1[2] * v2[1], v1[0] * v2[2], v1[1] * v2[2], v1[2] * v2[2]}
}

// OuterProd4 does the vector outer product
// of two vectors. The outer product produces an
// 3x4 matrix. E.G. a Vec3 * Vec4 = Mat3x4.
//
// The outer product can be thought of as the "opposite"
// of the Dot product. The Dot product treats both vectors like matrices
// oriented such that the left one has N columns and the right has N rows.
// So Vec3.Vec3 = Mat1x3*Mat3x1 = Mat1 = Scalar.
//
// The outer product orients it so they're facing "outward": Vec2*Vec3
// = Mat2x1*Mat1x3 = Mat2x3.
func (v1 Vec3[T]) OuterProd4(v2 Vec4[T]) Mat3x4[T] {
	return Mat3x4[T]{v1[0] * v2[0], v1[1] * v2[0], v1[2] * v2[0], v1[0] * v2[1], v1[1] * v2[1], v1[2] * v2[1], v1[0] * v2[2], v1[1] * v2[2], v1[2] * v2[2], v1[0] * v2[3], v1[1] * v2[3], v1[2] * v2[3]}
}

// Add performs element-wise addition between two vectors. It is equivalent to iterating
// over every element of v1 and adding the corresponding element of v2 to it.
func (v1 Vec4[T]) Add(v2 Vec4[T]) Vec4[T] {
	return Vec4[T]{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2], v1[3] + v2[3]}
}

// Sub performs element-wise subtraction between two vectors. It is equivalent to iterating
// over every element of v1 and subtracting the corresponding element of v2 from it.
func (v1 Vec4[T]) Sub(v2 Vec4[T]) Vec4[T] {
	return Vec4[T]{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2], v1[3] - v2[3]}
}

// Mul performs a scalar multiplication between the vector and some constant value
// c. This is equivalent to iterating over every vector element and multiplying by c.
func (v1 Vec4[T]) Mul(c T) Vec4[T] {
	return Vec4[T]{v1[0] * c, v1[1] * c, v1[2] * c, v1[3] * c}
}

// Dot returns the dot product of this vector with another. There are multiple ways
// to describe this value. One is the multiplication of their lengths and cos(theta) where
// theta is the angle between the vectors: v1.v2 = |v1||v2|cos(theta).
//
// The other (and what is actually done) is the sum of the element-wise multiplication of all
// elements. So for instance, two Vec3s would yield v1.x * v2.x + v1.y * v2.y + v1.z * v2.z.
//
// This means that the dot product of a vector and itself is the square of its Len (within
// the bounds of floating points error).
//
// The dot product is roughly a measure of how closely two vectors are to pointing in the same
// direction. If both vectors are normalized, the value will be -1 for opposite pointing,
// one for same pointing, and 0 for perpendicular vectors.
func (v1 Vec4[T]) Dot(v2 Vec4[T]) T {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2] + v1[3]*v2[3]
}

// Len returns the vector's length. Note that this is NOT the dimension of
// the vector (len(v)), but the mathematical length. This is equivalent to the square
// root of the sum of the squares of all elements. E.G. for a Vec2 it's
// math.Hypot(v[0], v[1]).
func (v1 Vec4[T]) Len() T {
	return T(math.Sqrt(float64(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2] + v1[3]*v1[3])))
}

// LenSqr returns the vector's square length. This is equivalent to the sum of the squares of all elements.
func (v1 Vec4[T]) LenSqr() T {
	return v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2] + v1[3]*v1[3]
}

// Normalize normalizes the vector. Normalization is (1/|v|)*v,
// making this equivalent to v.Scale(1/v.Len()). If the len is 0.0,
// this function will return an infinite value for all elements due
// to how floating point division works in Go (n/0.0 = math.Inf(Sign(n))).
//
// Normalization makes a vector's Len become 1.0 (within the margin of floating point error),
// while maintaining its directionality.
//
// (Can be seen here: http://play.golang.org/p/Aaj7SnbqIp )
func (v1 Vec4[T]) Normalize() Vec4[T] {
	l := 1.0 / v1.Len()
	return Vec4[T]{v1[0] * l, v1[1] * l, v1[2] * l, v1[3] * l}
}

// ApproxEqual takes in a vector and does an element-wise approximate float
// comparison as if FloatEqual had been used
func (v1 Vec4[T]) ApproxEqual(v2 Vec4[T]) bool {
	for i := range v1 {
		if !FloatEqual(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

// ApproxEqualThreshold takes in a threshold for comparing two floats, and uses
// it to do an element-wise comparison of the vector to another.
func (v1 Vec4[T]) ApproxEqualThreshold(v2 Vec4[T], threshold T) bool {
	for i := range v1 {
		if !FloatEqualThreshold(v1[i], v2[i], threshold) {
			return false
		}
	}
	return true
}

// ApproxFuncEqual takes in a func that compares two floats, and uses it to do an element-wise
// comparison of the vector to another. This is intended to be used with FloatEqualFunc
func (v1 Vec4[T]) ApproxFuncEqual(v2 Vec4[T], eq func(T, T) bool) bool {
	for i := range v1 {
		if !eq(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

// X is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec4[T]) X() T {
	return v[0]
}

// Y is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec4[T]) Y() T {
	return v[1]
}

// Z is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec4[T]) Z() T {
	return v[2]
}

// W is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec4[T]) W() T {
	return v[3]
}

// OuterProd2 does the vector outer product
// of two vectors. The outer product produces an
// 4x2 matrix. E.G. a Vec4 * Vec2 = Mat4x2.
//
// The outer product can be thought of as the "opposite"
// of the Dot product. The Dot product treats both vectors like matrices
// oriented such that the left one has N columns and the right has N rows.
// So Vec3.Vec3 = Mat1x3*Mat3x1 = Mat1 = Scalar.
//
// The outer product orients it so they're facing "outward": Vec2*Vec3
// = Mat2x1*Mat1x3 = Mat2x3.
func (v1 Vec4[T]) OuterProd2(v2 Vec2[T]) Mat4x2[T] {
	return Mat4x2[T]{v1[0] * v2[0], v1[1] * v2[0], v1[2] * v2[0], v1[3] * v2[0], v1[0] * v2[1], v1[1] * v2[1], v1[2] * v2[1], v1[3] * v2[1]}
}

// OuterProd3 does the vector outer product
// of two vectors. The outer product produces an
// 4x3 matrix. E.G. a Vec4 * Vec3 = Mat4x3.
//
// The outer product can be thought of as the "opposite"
// of the Dot product. The Dot product treats both vectors like matrices
// oriented such that the left one has N columns and the right has N rows.
// So Vec3.Vec3 = Mat1x3*Mat3x1 = Mat1 = Scalar.
//
// The outer product orients it so they're facing "outward": Vec2*Vec3
// = Mat2x1*Mat1x3 = Mat2x3.
func (v1 Vec4[T]) OuterProd3(v2 Vec3[T]) Mat4x3[T] {
	return Mat4x3[T]{v1[0] * v2[0], v1[1] * v2[0], v1[2] * v2[0], v1[3] * v2[0], v1[0] * v2[1], v1[1] * v2[1], v1[2] * v2[1], v1[3] * v2[1], v1[0] * v2[2], v1[1] * v2[2], v1[2] * v2[2], v1[3] * v2[2]}
}

// OuterProd4 does the vector outer product
// of two vectors. The outer product produces an
// 4x4 matrix. E.G. a Vec4 * Vec4 = Mat4.
//
// The outer product can be thought of as the "opposite"
// of the Dot product. The Dot product treats both vectors like matrices
// oriented such that the left one has N columns and the right has N rows.
// So Vec3.Vec3 = Mat1x3*Mat3x1 = Mat1 = Scalar.
//
// The outer product orients it so they're facing "outward": Vec2*Vec3
// = Mat2x1*Mat1x3 = Mat2x3.
func (v1 Vec4[T]) OuterProd4(v2 Vec4[T]) Mat4[T]{
	return Mat4[T]{v1[0] * v2[0], v1[1] * v2[0], v1[2] * v2[0], v1[3] * v2[0], v1[0] * v2[1], v1[1] * v2[1], v1[2] * v2[1], v1[3] * v2[1], v1[0] * v2[2], v1[1] * v2[2], v1[2] * v2[2], v1[3] * v2[2], v1[0] * v2[3], v1[1] * v2[3], v1[2] * v2[3], v1[3] * v2[3]}
}


func (v1 Vec2[T]) Reflect(v2 Vec2[T]) Vec2[T]{
	p := 2 * v1.Dot(v2)
	return Vec2[T]{v1[0] - p * v2[0],v1[1] - p * v2[1]}
}
func (v1 Vec3[T]) Reflect(v2 Vec3[T]) Vec3[T]{
	p := 2 * v1.Dot(v2)
	return Vec3[T]{v1[0] - p * v2[0],v1[1] - p * v2[1],v1[2] - p * v2[2] }
}
func (v1 Vec4[T]) Reflect(v2 Vec4[T]) Vec4[T]{
	p := 2 * v1.Dot(v2)
	return Vec4[T]{v1[0] - p * v2[0],v1[1] - p * v2[1],v1[2] - p * v2[2],v1[3] - p * v2[3]}
}
