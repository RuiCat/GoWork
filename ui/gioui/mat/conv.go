// This file is generated from mgl32/conv.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"math"
)

// CartesianToSpherical converts 3-dimensional cartesian coordinates (x,y,z) to spherical
// coordinates with radius r, inclination theta, and azimuth phi.
//
// All angles are in radians.
func CartesianToSpherical[T Float](coord Vec3[T]) (r, theta, phi T) {
	r = coord.Len()
	theta = T(math.Acos(float64(coord[2] / r)))
	phi = T(math.Atan2(float64(coord[1]), float64(coord[0])))

	return
}

// CartesianToCylindical converts 3-dimensional cartesian coordinates (x,y,z) to
// cylindrical coordinates with radial distance r, azimuth phi, and height z.
//
// All angles are in radians.
func CartesianToCylindical[T Float](coord Vec3[T]) (rho, phi, z T) {
	rho = T(math.Hypot(float64(coord[0]), float64(coord[1])))

	phi = T(math.Atan2(float64(coord[1]), float64(coord[0])))

	z = coord[2]

	return
}

// SphericalToCartesian converts spherical coordinates with radius r, inclination theta,
// and azimuth phi to cartesian coordinates (x,y,z).
//
// Angles are in radians.
func SphericalToCartesian[T Float](r, theta, phi T) Vec3[T] {
	st, ct := math.Sincos(float64(theta))
	sp, cp := math.Sincos(float64(phi))

	return Vec3[T]{r * T(st*cp), r * T(st*sp), r * T(ct)}
}

// SphericalToCylindrical converts spherical coordinates with radius r,
// inclination theta, and azimuth phi to cylindrical coordinates with radial
// distance r, azimuth phi, and height z.
//
// Angles are in radians
func SphericalToCylindrical[T Float](r, theta, phi T) (rho, phi2, z T) {
	s, c := math.Sincos(float64(theta))

	rho = r * T(s)
	z = r * T(c)
	phi2 = phi

	return
}

// CylindircalToSpherical converts cylindrical coordinates with radial distance
// r, azimuth phi, and height z to spherical coordinates with radius r,
// inclination theta, and azimuth phi.
//
// Angles are in radians
func CylindircalToSpherical[T Float](rho, phi, z T) (r, theta, phi2 T) {
	r = T(math.Hypot(float64(rho), float64(z)))
	phi2 = phi
	theta = T(math.Atan2(float64(rho), float64(z)))

	return
}

// CylindricalToCartesian converts cylindrical coordinates with radial distance
// r, azimuth phi, and height z to cartesian coordinates (x,y,z)
//
// Angles are in radians.
func CylindricalToCartesian[T Float](rho, phi, z T) Vec3[T] {
	s, c := math.Sincos(float64(phi))

	return Vec3[T]{rho * T(c), rho * T(s), z}
}

// DegToRad converts degrees to radians
func DegToRad[T Float](angle T) T {
	return angle * T(math.Pi) / 180
}

// RadToDeg converts radians to degrees
func RadToDeg[T Float](angle T) T {
	return angle * 180 / T(math.Pi)
}
