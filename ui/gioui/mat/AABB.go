package mat

import "math"

// RectangleToRectangle 检测矩阵与矩阵碰撞
//  oneX 长方形左上角X坐标
//  oneY 长方形右上角Y坐标
//  oneW 长方形宽
//  oneH 长方形高
//  twoX 长方形左上角X坐标
//  twoY 长方形右上角Y坐标
//  twoW 长方形宽
//  twoH 长方形高
func RectangleToRectangle[T Float](oneX, oneY, oneW, oneH, twoX, twoY, twoW, twoH T) bool {
	// x轴方向碰撞？
	collisionX := oneX+oneW >= twoX && twoX+twoW >= oneX
	// y轴方向碰撞？
	collisionY := oneY+oneH >= twoY && twoY+twoH >= oneY
	// 只有两个轴向都有碰撞时才碰撞
	return collisionX && collisionY
}

// CircleToCircle 检测圆与圆碰撞
//  oneX 圆心X坐标
//  oneY 圆心Y坐标
//  oneR 圆形半径
//  twoX 圆心X坐标
//  twoY 圆心Y坐标
//  twoR 圆形半径
func CircleToCircle[T Float](oneX, oneY, oneR, twoX, twoY, twoR T) bool {
	return math.Sqrt(math.Pow(float64(oneX-twoX), 2)+math.Pow(float64(oneY-twoY), 2)) > float64(oneR+twoR)
}

// CircleToSpot 检测圆与点碰撞
//  oneX 圆心X坐标
//  oneY 圆心Y坐标
//  oneR 圆形半径
//  X, Y 点坐标
func CircleToSpot[T Float](oneX, oneY, oneR, X, Y T) bool {
	return math.Sqrt(math.Pow(float64(oneX-X), 2)+math.Pow(float64(oneY-Y), 2)) <= float64(oneR)
}

// RectangleToSpot 检测矩阵与点碰撞
//  oneX 长方形左上角X坐标
//  oneY 长方形右上角Y坐标
//  oneW 长方形宽
//  oneH 长方形高
//  X, Y 点坐标
func RectangleToSpot[T Float](oneX, oneY, oneW, oneH, X, Y T) bool {
	return X >= oneX && X <= oneX+oneH && Y >= oneY && Y <= oneY+oneH
}

// CircleToRectangle 检测圆与矩阵碰撞
//  arcOx  圆心X坐标
//  arcOy  圆心Y坐标
//  arcR   圆形半径
//  rectX  长方形左上角X坐标
//  rectY  长方形右上角Y坐标
//  rectW  长方形宽
//  rectH  长方形高
func CircleToRectangle[T Float](arcOx, arcOy, arcR, rectX, rectY, rectW, rectH T) bool {
	// 分别判断矩形4个顶点与圆心的距离是否<=圆半径；如果<=，说明碰撞成功
	if ((rectX-arcOx)*(rectX-arcOx) + (rectY-arcOy)*(rectY-arcOy)) <= arcR*arcR {
		return true
	}
	if ((rectX+rectW-arcOx)*(rectX+rectW-arcOx) + (rectY-arcOy)*(rectY-arcOy)) <= arcR*arcR {
		return true
	}
	if ((rectX-arcOx)*(rectX-arcOx) + (rectY+rectH-arcOy)*(rectY+rectH-arcOy)) <= arcR*arcR {
		return true
	}
	if ((rectX+rectW-arcOx)*(rectX+rectW-arcOx) + (rectY+rectH-arcOy)*(rectY+rectH-arcOy)) <= arcR*arcR {
		return true
	}
	// 判断当圆心的Y坐标进入矩形内时X的位置，如果X在(rectX-arcR)到(rectX+rectW+arcR)这个范围内，则碰撞成功
	minDisX := T(0)
	if arcOy >= rectY && arcOy <= rectY+rectH {
		if arcOx < rectX {
			minDisX = rectX - arcOx
		} else if arcOx > rectX+rectW {
			minDisX = arcOx - rectX - rectW
		} else {
			return true
		}
		if minDisX <= arcR {
			return true
		}
	}
	// 判断当圆心的X坐标进入矩形内时Y的位置，如果X在(rectY-arcR)到(rectY+rectH+arcR)这个范围内，则碰撞成功
	minDisY := T(0)
	if arcOx >= rectX && arcOx <= rectX+rectW {
		if arcOy < rectY {
			minDisY = rectY - arcOy
		} else if arcOy > rectY+rectH {
			minDisY = arcOy - rectY - rectH
		} else {
			return true
		}
		if minDisY <= arcR {
			return true
		}
	}
	return false
}
