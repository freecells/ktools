package tmath

import (
	"fmt"
	"math"
	"strconv"
)

//Round 返回浮点数 指定长度小数位
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

//Round2 返回浮点数 指定长度小数位 use fmt.sprintf
func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}

//4/5
func FourOrFive(x float64) float64 {
	return math.Floor(x + 0.5)
}
