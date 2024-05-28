// This file is generated from mgl32/transform.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import "math"

// Rotate2D returns a rotation Matrix about a angle in 2-D space. Specifically about the origin.
// It is a 2x2 matrix, if you need a 3x3 for Homogeneous math (e.g. composition with a Translation matrix)
// see HomogRotate2D
func Rotate2D[T Float](angle T) Mat2[T] {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := T(math.Sin(float64(angle))), T(math.Cos(float64(angle)))
	return Mat2[T] {cos, sin, -sin, cos}
}

// Rotate3DX returns a 3x3 (non-homogeneous) Matrix that rotates by angle about the X-axis
//
// Where c is cos(angle) and s is sin(angle)
//    [1 0 0]
//    [0 c -s]
//    [0 s c ]
func Rotate3DX[T Float](angle T) Mat3[T] {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := T(math.Sin(float64(angle))), T(math.Cos(float64(angle)))

	return Mat3[T]{1, 0, 0, 0, cos, sin, 0, -sin, cos}
}

// Rotate3DY returns a 3x3 (non-homogeneous) Matrix that rotates by angle about the Y-axis
//
// Where c is cos(angle) and s is sin(angle)
//    [c 0 s]
//    [0 1 0]
//    [s 0 c ]
func Rotate3DY[T Float](angle T) Mat3[T] {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := T(math.Sin(float64(angle))), T(math.Cos(float64(angle)))

	return Mat3[T]{cos, 0, -sin, 0, 1, 0, sin, 0, cos}
}

// Rotate3DZ returns a 3x3 (non-homogeneous) Matrix that rotates by angle about the Z-axis
//
// Where c is cos(angle) and s is sin(angle)
//    [c -s 0]
//    [s c 0]
//    [0 0 1 ]
func Rotate3DZ[T Float](angle T) Mat3[T] {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := T(math.Sin(float64(angle))), T(math.Cos(float64(angle)))

	return Mat3[T]{cos, sin, 0, -sin, cos, 0, 0, 0, 1}
}

// Translate2D returns a homogeneous (3x3 for 2D-space) Translation matrix that moves a point by Tx units in the x-direction and Ty units in the y-direction
//
//    [[1, 0, Tx]]
//    [[0, 1, Ty]]
//    [[0, 0, 1 ]]
func Translate2D[T Float](Tx, Ty T) Mat3[T] {
	return Mat3[T]{1, 0, 0, 0, 1, 0, T(Tx), T(Ty), 1}
}

// Translate3D returns a homogeneous (4x4 for 3D-space) Translation matrix that moves a point by Tx units in the x-direction, Ty units in the y-direction,
// and Tz units in the z-direction
//
//    [[1, 0, 0, Tx]]
//    [[0, 1, 0, Ty]]
//    [[0, 0, 1, Tz]]
//    [[0, 0, 0, 1 ]]
func Translate3D[T Float](Tx, Ty, Tz T) Mat4[T] {
	return Mat4[T]{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, T(Tx), T(Ty), T(Tz), 1}
}

// HomogRotate2D is the same as Rotate2D, except homogeneous (3x3 with the extra row/col being all zeroes with a one in the bottom right)
func HomogRotate2D[T Float](angle T) Mat3[T] {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := T(math.Sin(float64(angle))), T(math.Cos(float64(angle)))
	return Mat3[T]{cos, sin, 0, -sin, cos, 0, 0, 0, 1}
}

// HomogRotate3DX is the same as Rotate3DX, except homogeneous (4x4 with the extra row/col being all zeroes with a one in the bottom right)
func HomogRotate3DX[T Float](angle T) Mat4[T] {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := T(math.Sin(float64(angle))), T(math.Cos(float64(angle)))

	return Mat4[T]{1, 0, 0, 0, 0, cos, sin, 0, 0, -sin, cos, 0, 0, 0, 0, 1}
}

// HomogRotate3DY is the same as Rotate3DY, except homogeneous (4x4 with the extra row/col being all zeroes with a one in the bottom right)
func HomogRotate3DY[T Float](angle T) Mat4[T] {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := T(math.Sin(float64(angle))), T(math.Cos(float64(angle)))
	return Mat4[T]{cos, 0, -sin, 0, 0, 1, 0, 0, sin, 0, cos, 0, 0, 0, 0, 1}
}

// HomogRotate3DZ is the same as Rotate3DZ, except homogeneous (4x4 with the extra row/col being all zeroes with a one in the bottom right)
func HomogRotate3DZ[T Float](angle T) Mat4[T] {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := T(math.Sin(float64(angle))), T(math.Cos(float64(angle)))
	return Mat4[T]{cos, sin, 0, 0, -sin, cos, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// Scale3D creates a homogeneous 3D scaling matrix
// [[ scaleX, 0     , 0     , 0 ]]
// [[ 0     , scaleY, 0     , 0 ]]
// [[ 0     , 0     , scaleZ, 0 ]]
// [[ 0     , 0     , 0     , 1 ]]
func Scale3D[T Float](scaleX, scaleY, scaleZ T) Mat4[T] {

	return Mat4[T]{T(scaleX), 0, 0, 0, 0, T(scaleY), 0, 0, 0, 0, T(scaleZ), 0, 0, 0, 0, 1}
}

// Scale2D creates a homogeneous 2D scaling matrix
// [[ scaleX, 0     , 0 ]]
// [[ 0     , scaleY, 0 ]]
// [[ 0     , 0     , 1 ]]
func Scale2D[T Float](scaleX, scaleY T) Mat3[T] {
	return Mat3[T]{T(scaleX), 0, 0, 0, T(scaleY), 0, 0, 0, 1}
}

// ShearX2D creates a homogeneous 2D shear matrix along the X-axis
func ShearX2D[T Float](shear T) Mat3[T] {
	return Mat3[T]{1, 0, 0, T(shear), 1, 0, 0, 0, 1}
}

// ShearY2D creates a homogeneous 2D shear matrix along the Y-axis
func ShearY2D[T Float](shear T) Mat3[T] {
	return Mat3[T]{1, T(shear), 0, 0, 1, 0, 0, 0, 1}
}

// ShearX3D creates a homogeneous 3D shear matrix along the X-axis
func ShearX3D[T Float](shearY, shearZ T) Mat4[T] {

	return Mat4[T]{1, T(shearY), T(shearZ), 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// ShearY3D creates a homogeneous 3D shear matrix along the Y-axis
func ShearY3D[T Float](shearX, shearZ T) Mat4[T] {
	return Mat4[T]{1, 0, 0, 0, T(shearX), 1, T(shearZ), 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// ShearZ3D creates a homogeneous 3D shear matrix along the Z-axis
func ShearZ3D[T Float](shearX, shearY T) Mat4[T] {
	return Mat4[T]{1, 0, 0, 0, 0, 1, 0, 0, T(shearX), T(shearY), 1, 0, 0, 0, 0, 1}
}

// HomogRotate3D creates a 3D rotation Matrix that rotates by (radian) angle about some arbitrary axis given by a normalized Vector.
// It produces a homogeneous matrix (4x4)
//
// Where c is cos(angle) and s is sin(angle), and x, y, and z are the first, second, and third elements of the axis vector (respectively):
//
//    [[ x^2(1-c)+c, xy(1-c)-zs, xz(1-c)+ys, 0 ]]
//    [[ xy(1-c)+zs, y^2(1-c)+c, yz(1-c)-xs, 0 ]]
//    [[ xz(1-c)-ys, yz(1-c)+xs, z^2(1-c)+c, 0 ]]
//    [[ 0         , 0         , 0         , 1 ]]
func HomogRotate3D[T Float](angle T, axis Vec3[T]) Mat4[T] {
	x, y, z := axis[0], axis[1], axis[2]
	s, c := T(math.Sin(float64(angle))), T(math.Cos(float64(angle)))
	k := 1 - c

	return Mat4[T]{x*x*k + c, x*y*k + z*s, x*z*k - y*s, 0, x*y*k - z*s, y*y*k + c, y*z*k + x*s, 0, x*z*k + y*s, y*z*k - x*s, z*z*k + c, 0, 0, 0, 0, 1}
}

// Extract3DScale extracts the 3d scaling from a homogeneous matrix
func Extract3DScale[T Float](m Mat4[T]) (x, y, z float64) {
	return float64(math.Sqrt(float64(m[0]*m[0] + m[1]*m[1] + m[2]*m[2]))),
		float64(math.Sqrt(float64(m[4]*m[4] + m[5]*m[5] + m[6]*m[6]))),
		float64(math.Sqrt(float64(m[8]*m[8] + m[9]*m[9] + m[10]*m[10])))
}

// ExtractMaxScale extracts the maximum scaling from a homogeneous matrix
func ExtractMaxScale[T Float](m Mat4[T]) float64 {
	scaleX := float64(m[0]*m[0] + m[1]*m[1] + m[2]*m[2])
	scaleY := float64(m[4]*m[4] + m[5]*m[5] + m[6]*m[6])
	scaleZ := float64(m[8]*m[8] + m[9]*m[9] + m[10]*m[10])

	return float64(math.Sqrt(math.Max(scaleX, math.Max(scaleY, scaleZ))))
}

// Mat4Normal calculates the Normal of the Matrix (aka the inverse transpose)
func Mat4Normal[T Float](m Mat4[T]) Mat3[T] {
	n := m.Inv().Transpose()
	return Mat3[T]{n[0], n[1], n[2], n[4], n[5], n[6], n[8], n[9], n[10]}
}

// TransformCoordinate multiplies a 3D vector by a transformation given by
// the homogeneous 4D matrix m, applying any translation.
// If this transformation is non-affine, it will project this
// vector onto the plane w=1 before returning the result.
//
// This is similar to saying you're transforming and projecting a point.
//
// This is effectively equivalent to the GLSL code
//     vec4 r = (m * vec4(v,1.));
//     r = r/r.w;
//     vec3 newV = r.xyz;
func TransformCoordinate[T Float](v Vec3[T], m Mat4[T]) Vec3[T] {
	t := v.Vec4(1)
	t = m.Mul4x1(t)
	t = t.Mul(1 / t[3])

	return t.Vec3()
}

// TransformNormal multiplies a 3D vector by a transformation given by
// the homogeneous 4D matrix m, NOT applying any translations.
//
// This is similar to saying you're applying a transformation
// to a direction or normal. Rotation still applies (as does scaling),
// but translating a direction or normal is meaningless.
//
// This is effectively equivalent to the GLSL code
//    vec4 r = (m * vec4(v,0.));
//    vec3 newV = r.xyz
func TransformNormal[T Float](v Vec3[T], m Mat4[T]) Vec3[T] {
	t := v.Vec4(0)
	t = m.Mul4x1(t)

	return t.Vec3()
}
