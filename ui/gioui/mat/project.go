// This file is generated from mgl32/project.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"errors"
	"math"
)

// Ortho generates an Ortho Matrix.
func Ortho[T Float](left, right, bottom, top, near, far T) Mat4[T] {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)
	return Mat4[T]{T(2. / rml), 0, 0, 0, 0, T(2. / tmb), 0, 0, 0, 0, T(-2. / fmn), 0, T(-(right + left) / rml), T(-(top + bottom) / tmb), T(-(far + near) / fmn), 1}
}

// Ortho2D is equivalent to Ortho with the near and far planes being -1 and 1,
// respectively.
func Ortho2D[T Float](left, right, bottom, top T) Mat4[T] {
	return Ortho(left, right, bottom, top, -1, 1)
}

// Perspective generates a Perspective Matrix.
func Perspective[T Float](fovy, aspect, near, far T) Mat4[T] {
	// fovy = (fovy * math.Pi) / 180.0 // convert from degrees to radians
	nmf, f := near-far, T(1./math.Tan(float64(fovy)/2.0))

	return Mat4[T]{T(f / aspect), 0, 0, 0, 0, T(f), 0, 0, 0, 0, T((near + far) / nmf), -1, 0, 0, T((2. * far * near) / nmf), 0}
}

// Frustum generates a Frustum Matrix.
func Frustum[T Float](left, right, bottom, top, near, far T) Mat4[T] {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)
	A, B, C, D := (right+left)/rml, (top+bottom)/tmb, -(far+near)/fmn, -(2*far*near)/fmn

	return Mat4[T]{T((2. * near) / rml), 0, 0, 0, 0, T((2. * near) / tmb), 0, 0, T(A), T(B), T(C), -1, 0, 0, T(D), 0}
}

// LookAt generates a transform matrix from world space to the given eye space.
func LookAt[T Float](eyeX, eyeY, eyeZ, centerX, centerY, centerZ, upX, upY, upZ T) Mat4[T] {
	return LookAtV[T](Vec3[T]{eyeX, eyeY, eyeZ}, Vec3[T]{centerX, centerY, centerZ}, Vec3[T]{upX, upY, upZ})
}

// LookAtV generates a transform matrix from world space into the specific eye
// space.
func LookAtV[T Float](eye, center, up Vec3[T]) Mat4[T] {
	f := center.Sub(eye).Normalize()
	s := f.Cross(up.Normalize()).Normalize()
	u := s.Cross(f)

	M := Mat4[T]{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		0, 0, 0, 1,
	}

	return M.Mul4(Translate3D(T(-eye[0]), T(-eye[1]), T(-eye[2])))
}

// Project transforms a set of coordinates from object space (in obj) to window
// coordinates (with depth).
//
// Window coordinates are continuous, not discrete (well, as continuous as an
// IEEE Floating Point can be), so you won't get exact pixel locations without
// rounding or similar
func Project[T Float](obj Vec3[T], modelview, projection Mat4[T], initialX, initialY, width, height int) (win Vec3[T]) {
	obj4 := obj.Vec4(1)

	vpp := projection.Mul4(modelview).Mul4x1(obj4)
	vpp = vpp.Mul(1 / vpp.W())
	win[0] = T(initialX) + (T(width)*(vpp[0]+1))/2
	win[1] = T(initialY) + (T(height)*(vpp[1]+1))/2
	win[2] = (vpp[2] + 1) / 2

	return win
}

// UnProject transforms a set of window coordinates to object space. If your MVP
// (projection.Mul(modelview) matrix is not invertible, this will return an
// error.
//
// Note that the projection may not be perfect if you use strict pixel locations
// rather than the exact values given by Projectf. (It's still unlikely to be
// perfect due to precision errors, but it will be closer)
func UnProject[T Float](win Vec3[T], modelview, projection Mat4[T], initialX, initialY, width, height int) (obj Vec3[T], err error) {
	inv := projection.Mul4(modelview).Inv()
	var blank Mat4[T]
	if inv == blank {
		return Vec3[T]{}, errors.New("Could not find matrix inverse (projection times modelview is probably non-singular)")
	}

	obj4 := inv.Mul4x1(Vec4[T]{
		(2 * (win[0] - T(initialX)) / T(width)) - 1,
		(2 * (win[1] - T(initialY)) / T(height)) - 1,
		2*win[2] - 1,
		1.0,
	})
	obj = obj4.Vec3()

	//if obj4[3] > MinValue {}
	obj[0] /= obj4[3]
	obj[1] /= obj4[3]
	obj[2] /= obj4[3]

	return obj, nil
}
