package f32

import (
	"ui/gioui/mat"
)

type Vec2 = mat.Vec2[float32]
type Vec3 = mat.Vec3[float32]
type Vec4 = mat.Vec4[float32]
type VecN = mat.VecN[float32]
type Mat2 = mat.Vec4[float32]
type Mat2x3 = mat.Mat2x3[float32]
type Mat2x4 = mat.Mat2x4[float32]
type Mat3x2 = mat.Mat3x2[float32]
type Mat3 = mat.Mat3[float32]
type Mat3x4 = mat.Mat3x4[float32]
type Mat4x2 = mat.Mat4x2[float32]
type Mat4x3 = mat.Mat4x3[float32]
type Mat4 = mat.Mat4[float32]
type MatMxN = mat.MatMxN[float32]
type Affine2D = mat.Affine2D[float32]
type Point = mat.Point[float32]
type Rectangle = mat.Rectangle[float32]

var (
	NewAffine2D = mat.NewAffine2D[float32]
)

var (
	RectangleToRectangle = mat.RectangleToRectangle[float32]
	CircleToCircle       = mat.CircleToCircle[float32]
	CircleToSpot         = mat.CircleToSpot[float32]
	RectangleToSpot      = mat.RectangleToSpot[float32]
	CircleToRectangle    = mat.CircleToRectangle[float32]
)

var (
	CartesianToSpherical   = mat.CartesianToSpherical[float32]
	CartesianToCylindical  = mat.CartesianToCylindical[float32]
	SphericalToCartesian   = mat.SphericalToCartesian[float32]
	SphericalToCylindrical = mat.SphericalToCylindrical[float32]
	CylindircalToSpherical = mat.CylindircalToSpherical[float32]
	CylindricalToCartesian = mat.CylindricalToCartesian[float32]
	DegToRad               = mat.DegToRad[float32]
	RadToDeg               = mat.RadToDeg[float32]
)

var (
	NewMatrix         = mat.NewMatrix[float32]
	NewMatrixFromData = mat.NewMatrixFromData[float32]
	CopyMatMN         = mat.CopyMatMN[float32]
	IdentN            = mat.IdentN[float32]
	DiagN             = mat.DiagN[float32]
)

var (
	Ident2         = mat.Ident2[float32]
	Diag2          = mat.Diag2[float32]
	Mat2FromRows   = mat.Mat2FromRows[float32]
	Mat2FromCols   = mat.Mat2FromCols[float32]
	Mat2x3FromRows = mat.Mat2x3FromRows[float32]
	Mat2x3FromCols = mat.Mat2x3FromCols[float32]
	Mat2x4FromRows = mat.Mat2x4FromRows[float32]
	Mat2x4FromCols = mat.Mat2x4FromCols[float32]
	Mat3x2FromRows = mat.Mat3x2FromRows[float32]
	Mat3x2FromCols = mat.Mat3x2FromCols[float32]
	Ident3         = mat.Ident3[float32]
	Diag3          = mat.Diag3[float32]
	Mat3FromRows   = mat.Mat3FromRows[float32]
	Mat3FromCols   = mat.Mat3FromCols[float32]
	Mat3x4FromRows = mat.Mat3x4FromRows[float32]
	Mat3x4FromCols = mat.Mat3x4FromCols[float32]
	Mat4x2FromRows = mat.Mat4x2FromRows[float32]
	Mat4x2FromCols = mat.Mat4x2FromCols[float32]
	Mat4x3FromRows = mat.Mat4x3FromRows[float32]
	Mat4x3FromCols = mat.Mat4x3FromCols[float32]
	Ident4         = mat.Ident4[float32]
	Diag4          = mat.Diag4[float32]
	Mat4FromRows   = mat.Mat4FromRows[float32]
	Mat4FromCols   = mat.Mat4FromCols[float32]
)
var (
	Pt    = mat.Pt[float32]
	Rect  = mat.Rect[float32]
	FRect = mat.FRect[float32]
	FPt   = mat.FPt[float32]
)

var (
	Ortho       = mat.Ortho[float32]
	Ortho2D     = mat.Ortho2D[float32]
	Perspective = mat.Perspective[float32]
	Frustum     = mat.Frustum[float32]
	LookAt      = mat.LookAt[float32]
	LookAtV     = mat.LookAtV[float32]
	Project     = mat.Project[float32]
	UnProject   = mat.UnProject[float32]
)

var (
	QuatIdent          = mat.QuatIdent[float32]
	QuatRotate         = mat.QuatRotate[float32]
	QuatSlerp          = mat.QuatSlerp[float32]
	QuatLerp           = mat.QuatLerp[float32]
	QuatNlerp          = mat.QuatNlerp[float32]
	AnglesToQuat       = mat.AnglesToQuat[float32]
	Mat4ToQuat         = mat.Mat4ToQuat[float32]
	QuatLookAtV        = mat.QuatLookAtV[float32]
	QuatBetweenVectors = mat.QuatBetweenVectors[float32]
)

var (
	Circle                    = mat.Circle[float32]
	RectVec2                  = mat.RectVec2[float32]
	QuadraticBezierCurve2D    = mat.QuadraticBezierCurve2D[float32]
	QuadraticBezierCurve3D    = mat.QuadraticBezierCurve3D[float32]
	CubicBezierCurve2D        = mat.CubicBezierCurve2D[float32]
	CubicBezierCurve3D        = mat.CubicBezierCurve3D[float32]
	BezierCurve2D             = mat.BezierCurve2D[float32]
	BezierCurve3D             = mat.BezierCurve3D[float32]
	MakeBezierCurve2D         = mat.MakeBezierCurve2D[float32]
	MakeBezierCurve3D         = mat.MakeBezierCurve3D[float32]
	BezierSurface             = mat.BezierSurface[float32]
	BezierSplineInterpolate2D = mat.BezierSplineInterpolate2D[float32]
	BezierSplineInterpolate3D = mat.BezierSplineInterpolate3D[float32]
	ReticulateSplines         = mat.ReticulateSplines[float32]
	ScreenToGLCoords          = mat.ScreenToGLCoords[float32]
	GLToScreenCoords          = mat.GLToScreenCoords[float32]
)

var (
	Abs                 = mat.Abs[float32]
	FloatEqual          = mat.FloatEqual[float32]
	FloatEqualFunc      = mat.FloatEqualFunc[float32]
	FloatEqualThreshold = mat.FloatEqualThreshold[float32]
	Clamp               = mat.Clamp[float32]
	ClampFunc           = mat.ClampFunc[float32]
	IsClamped           = mat.IsClamped[float32]
	SetMin              = mat.SetMin[float32]
	SetMax              = mat.SetMax[float32]
	Round               = mat.Round[float32]
)

var (
	NewVecNFromData = mat.NewVecNFromData[float32]
	NewVecN         = mat.NewVecN[float32]
)

var (
	Rotate2D            = mat.Rotate2D[float32]
	Rotate3DX           = mat.Rotate3DX[float32]
	Rotate3DY           = mat.Rotate3DY[float32]
	Rotate3DZ           = mat.Rotate3DZ[float32]
	Translate2D         = mat.Translate2D[float32]
	Translate3D         = mat.Translate3D[float32]
	HomogRotate2D       = mat.HomogRotate2D[float32]
	HomogRotate3DX      = mat.HomogRotate3DX[float32]
	HomogRotate3DY      = mat.HomogRotate3DY[float32]
	HomogRotate3DZ      = mat.HomogRotate3DZ[float32]
	Scale3D             = mat.Scale3D[float32]
	Scale2D             = mat.Scale2D[float32]
	ShearX2D            = mat.ShearX2D[float32]
	ShearY2D            = mat.ShearY2D[float32]
	ShearX3D            = mat.ShearX3D[float32]
	ShearY3D            = mat.ShearY3D[float32]
	ShearZ3D            = mat.ShearZ3D[float32]
	HomogRotate3D       = mat.HomogRotate3D[float32]
	Extract3DScale      = mat.Extract3DScale[float32]
	ExtractMaxScale     = mat.ExtractMaxScale[float32]
	Mat4Normal          = mat.Mat4Normal[float32]
	TransformCoordinate = mat.TransformCoordinate[float32]
	TransformNormal     = mat.TransformNormal[float32]
)
