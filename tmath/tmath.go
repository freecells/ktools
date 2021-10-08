package tmath

import (
	"errors"
	"math"
)

//CalcAngle 计算两点之间夹角
func CalcAngle(x1, y1, x2, y2 float64) (angle float64) {

	x := math.Abs(x1 - x2)
	y := math.Abs(y1 - y2)

	z := math.Sqrt(x*x + y*y)

	angle = Round((math.Asin(y/z) / math.Pi * 180), 2)

	if y1 > y2 {
		angle = -angle
	}

	return
}

//PointRotateNext 旋转点位
func PointRotateNext(dx, dy, x1, y1 float64, angle float64) (x, y float64) {

	/* 	平面中，一个点(x,y)绕任意点(dx,dy)顺时针旋转a度后的坐标

	xx= (x - dx)*cos(-a) - (y - dy)*sin(-a) + dx ;

	yy= (x - dx)*sin(-a) + (y - dy)*cos(-a) +dy ;

	平面中，一个点(x,y)绕任意点(dx,dy)逆时针旋转a度后的坐标

	xx= (x - dx)*cos(a) - (y - dy)*sin(a) + dx ;

	yy= (x - dx)*sin(a) + (y - dy)*cos(a) +dy ; */

	// if dy < y1 {
	// 	angle = -angle //负角度 为顺时针旋转
	// }

	x = (x1-dx)*math.Cos(angle*math.Pi/180) - (y1-dy)*math.Sin(angle*math.Pi/180) + dx
	y = (x1-dx)*math.Cos(angle*math.Pi/180) + (y1-dy)*math.Sin(angle*math.Pi/180) + dy

	// x = (x1-dx)*math.Cos(angle) - (y1-dy)*math.Sin(angle) + dx
	// y = (x1-dx)*math.Cos(angle) + (y1-dy)*math.Sin(angle) + dy

	return
}

// 根据原点0。0 坐标系旋转 angle
func XytoX1Y1(tx, ty, angle float64) (x, y float64) {

	x = tx*math.Cos(angle*math.Pi/180) + ty*math.Sin(angle*math.Pi/180)

	y = ty*math.Cos(angle*math.Pi/180) - tx*math.Sin(angle*math.Pi/180)

	return
}

//余弦计算
func CosineSimilar(a []float64, b []float64) (cosine float64, err error) {
	fangDa := 1.0
	count := 0
	length_a := len(a)
	length_b := len(b)
	if length_a > length_b {
		count = length_a
	} else {
		count = length_b
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= length_a {
			s2 += math.Pow(b[k], 2)
			continue
		}
		if k >= length_b {
			s1 += math.Pow(a[k], 2)
			continue
		}
		sumA += a[k] * b[k]
		s1 += math.Pow(a[k], 2)
		s2 += math.Pow(b[k], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, errors.New("Vectors should not be null (all zeros)")
	}

	calcRes := (sumA / (math.Sqrt(s1) * math.Sqrt(s2))) * fangDa
	return calcRes, nil
	// return sumA / (math.Sqrt(s1 * s2)), nil
}

//Distance 二维空间两点距离
func Distance(x1, y1, x2, y2 float64) (dt float64) {

	dt = math.Sqrt(math.Pow((x1-x2), 2) + math.Pow((y1-y2), 2))

	//percent distance
	// dt = 1 / (1 + dt)

	return
}

//ItoFloat64 整数转换到float64
func ItoFloat64(some interface{}) (res float64) {

	switch val := some.(type) {
	case float32:
		res = float64(val)
		return
	case float64:
		res = val
		return
	case int:
		res = float64(val)
		return

	default:
		res = 0
		return
	}

	return
}
