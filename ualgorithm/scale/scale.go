package uscale

//todo: z-score标准化(或零－均值标准化)

// MinMaxScale 最小－最大规范化(线性变换)
func MinMaxScale(value, sourceMin, sourceMax, targetMin, targetMax float64) float64 {
	if sourceMin == sourceMax {
		return targetMin
	}
	factor := (targetMax - targetMin) / (sourceMax - sourceMin)
	return factor*(value-sourceMin) + targetMin
}

//todo: 小数定标规范化：通过移动X的小数位置来进行规范化

//todo: 对数Logistic模式

//todo: 模糊量化模式
