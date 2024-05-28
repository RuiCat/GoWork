package f64

import (
	"ui/gioui/mat"
)

type Vec2 = mat.Vec2[float64]
type Vec3 = mat.Vec3[float64]
type Vec4 = mat.Vec4[float64]
type VecN = mat.VecN[float64]
type Mat2 = mat.Vec4[float64]
type Mat2x3 = mat.Mat2x3[float64]
type Mat2x4 = mat.Mat2x4[float64]
type Mat3x2 = mat.Mat3x2[float64]
type Mat3 = mat.Mat3[float64]
type Mat3x4 = mat.Mat3x4[float64]
type Mat4x2 = mat.Mat4x2[float64]
type Mat4x3 = mat.Mat4x3[float64]
type Mat4 = mat.Mat4[float64]
type MatMxN = mat.MatMxN[float64]
type Affine2D = mat.Affine2D[float64]
type Point = mat.Point[float64]
type Rectangle = mat.Rectangle[float64]

var (
	NewAffine2D = mat.NewAffine2D[float64]
)

var (
	RectangleToRectangle = mat.RectangleToRectangle[float64]
	CircleToCircle       = mat.CircleToCircle[float64]
	CircleToSpot         = mat.CircleToSpot[float64]
	RectangleToSpot      = mat.RectangleToSpot[float64]
	CircleToRectangle    = mat.CircleToRectangle[float64]
)

var (
	CartesianToSpherical   = mat.CartesianToSpherical[float64]
	CartesianToCylindical  = mat.CartesianToCylindical[float64]
	SphericalToCartesian   = mat.SphericalToCartesian[float64]
	SphericalToCylindrical = mat.SphericalToCylindrical[float64]
	CylindircalToSpherical = mat.CylindircalToSpherical[float64]
	CylindricalToCartesian = mat.CylindricalToCartesian[float64]
	DegToRad               = mat.DegToRad[float64]
	RadToDeg               = mat.RadToDeg[float64]
)

var (
	NewMatrix         = mat.NewMatrix[float64]
	NewMatrixFromData = mat.NewMatrixFromData[float64]
	CopyMatMN         = mat.CopyMatMN[float64]
	IdentN            = mat.IdentN[float64]
	DiagN             = mat.DiagN[float64]
)

var (
	Ident2         = mat.Ident2[float64]
	Diag2          = mat.Diag2[float64]
	Mat2FromRows   = mat.Mat2FromRows[float64]
	Mat2FromCols   = mat.Mat2FromCols[float64]
	Mat2x3FromRows = mat.Mat2x3FromRows[float64]
	Mat2x3FromCols = mat.Mat2x3FromCols[float64]
	Mat2x4FromRows = mat.Mat2x4FromRows[float64]
	Mat2x4FromCols = mat.Mat2x4FromCols[float64]
	Mat3x2FromRows = mat.Mat3x2FromRows[float64]
	Mat3x2FromCols = mat.Mat3x2FromCols[float64]
	Ident3         = mat.Ident3[float64]
	Diag3          = mat.Diag3[float64]
	Mat3FromRows   = mat.Mat3FromRows[float64]
	Mat3FromCols   = mat.Mat3FromCols[float64]
	Mat3x4FromRows = mat.Mat3x4FromRows[float64]
	Mat3x4FromCols = mat.Mat3x4FromCols[float64]
	Mat4x2FromRows = mat.Mat4x2FromRows[float64]
	Mat4x2FromCols = mat.Mat4x2FromCols[float64]
	Mat4x3FromRows = mat.Mat4x3FromRows[float64]
	Mat4x3FromCols = mat.Mat4x3FromCols[float64]
	Ident4         = mat.Ident4[float64]
	Diag4          = mat.Diag4[float64]
	Mat4FromRows   = mat.Mat4FromRows[float64]
	Mat4FromCols   = mat.Mat4FromCols[float64]
)
var (
	Pt    = mat.Pt[float64]
	Rect  = mat.Rect[float64]
	FRect = mat.FRect[float64]
	FPt   = mat.FPt[float64]
)

var (
	Ortho       = mat.Ortho[float64]
	Ortho2D     = mat.Ortho2D[float64]
	Perspective = mat.Perspective[float64]
	Frustum     = mat.Frustum[float64]
	LookAt      = mat.LookAt[float64]
	LookAtV     = mat.LookAtV[float64]
	Project     = mat.Project[float64]
	UnProject   = mat.UnProject[float64]
)

var (
	QuatIdent          = mat.QuatIdent[float64]
	QuatRotate         = mat.QuatRotate[float64]
	QuatSlerp          = mat.QuatSlerp[float64]
	QuatLerp           = mat.QuatLerp[float64]
	QuatNlerp          = mat.QuatNlerp[float64]
	AnglesToQuat       = mat.AnglesToQuat[float64]
	Mat4ToQuat         = mat.Mat4ToQuat[float64]
	QuatLookAtV        = mat.QuatLookAtV[float64]
	QuatBetweenVectors = mat.QuatBetweenVectors[float64]
)

var (
	Circle                    = mat.Circle[float64]
	RectVec2                  = mat.RectVec2[float64]
	QuadraticBezierCurve2D    = mat.QuadraticBezierCurve2D[float64]
	QuadraticBezierCurve3D    = mat.QuadraticBezierCurve3D[float64]
	CubicBezierCurve2D        = mat.CubicBezierCurve2D[float64]
	CubicBezierCurve3D        = mat.CubicBezierCurve3D[float64]
	BezierCurve2D             = mat.BezierCurve2D[float64]
	BezierCurve3D             = mat.BezierCurve3D[float64]
	MakeBezierCurve2D         = mat.MakeBezierCurve2D[float64]
	MakeBezierCurve3D         = mat.MakeBezierCurve3D[float64]
	BezierSurface             = mat.BezierSurface[float64]
	BezierSplineInterpolate2D = mat.BezierSplineInterpolate2D[float64]
	BezierSplineInterpolate3D = mat.BezierSplineInterpolate3D[float64]
	ReticulateSplines         = mat.ReticulateSplines[float64]
	ScreenToGLCoords          = mat.ScreenToGLCoords[float64]
	GLToScreenCoords          = mat.GLToScreenCoords[float64]
)

var (
	Abs                 = mat.Abs[float64]
	FloatEqual          = mat.FloatEqual[float64]
	FloatEqualFunc      = mat.FloatEqualFunc[float64]
	FloatEqualThreshold = mat.FloatEqualThreshold[float64]
	Clamp               = mat.Clamp[float64]
	ClampFunc           = mat.ClampFunc[float64]
	IsClamped           = mat.IsClamped[float64]
	SetMin              = mat.SetMin[float64]
	SetMax              = mat.SetMax[float64]
	Round               = mat.Round[float64]
)

var (
	NewVecNFromData = mat.NewVecNFromData[float64]
	NewVecN         = mat.NewVecN[float64]
)

var (
	Rotate2D            = mat.Rotate2D[float64]
	Rotate3DX           = mat.Rotate3DX[float64]
	Rotate3DY           = mat.Rotate3DY[float64]
	Rotate3DZ           = mat.Rotate3DZ[float64]
	Translate2D         = mat.Translate2D[float64]
	Translate3D         = mat.Translate3D[float64]
	HomogRotate2D       = mat.HomogRotate2D[float64]
	HomogRotate3DX      = mat.HomogRotate3DX[float64]
	HomogRotate3DY      = mat.HomogRotate3DY[float64]
	HomogRotate3DZ      = mat.HomogRotate3DZ[float64]
	Scale3D             = mat.Scale3D[float64]
	Scale2D             = mat.Scale2D[float64]
	ShearX2D            = mat.ShearX2D[float64]
	ShearY2D            = mat.ShearY2D[float64]
	ShearX3D            = mat.ShearX3D[float64]
	ShearY3D            = mat.ShearY3D[float64]
	ShearZ3D            = mat.ShearZ3D[float64]
	HomogRotate3D       = mat.HomogRotate3D[float64]
	Extract3DScale      = mat.Extract3DScale[float64]
	ExtractMaxScale     = mat.ExtractMaxScale[float64]
	Mat4Normal          = mat.Mat4Normal[float64]
	TransformCoordinate = mat.TransformCoordinate[float64]
	TransformNormal     = mat.TransformNormal[float64]
)
