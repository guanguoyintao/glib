package udistance

import (
	"context"
	"math"
	"strings"
)

func JaroWinklerDistance(ctx context.Context, a, b string) float64 {
	// 将字符串转换为小写
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	// 计算Jaro距离和匹配前缀长度
	jaro, prefixLen := jaroAndPrefix(a, b)

	// Jaro-Winkler常数
	jaroWinklerScalingFactor := 0.1
	jaroWinklerPrefixBoostThreshold := 0.7

	// 如果匹配前缀长度超过阈值，则应用前缀加权
	if float64(prefixLen)/float64(len(a)) > jaroWinklerPrefixBoostThreshold {
		jaro += float64(prefixLen) * jaroWinklerScalingFactor * (1.0 - jaro)
	}

	return jaro
}

func jaroAndPrefix(a, b string) (float64, int) {
	if a == b {
		return 1.0, len(a)
	}

	matchDistance := int(math.Floor(float64(math.Max(float64(len(a)), float64(len(b))))/2.0)) - 1
	aMatches := make([]bool, len(a))
	bMatches := make([]bool, len(b))
	matches := 0

	// 统计匹配并计算转置数
	for i := 0; i < len(a); i++ {
		start := int(math.Max(0, float64(i-matchDistance)))
		end := int(math.Min(float64(i+matchDistance+1), float64(len(b))))

		for j := start; j < end; j++ {
			if !bMatches[j] && a[i] == b[j] {
				aMatches[i] = true
				bMatches[j] = true
				matches++
				break
			}
		}
	}

	if matches == 0 {
		return 0.0, 0
	}

	transpositions := 0
	k := 0
	for i := 0; i < len(a); i++ {
		if aMatches[i] {
			for !bMatches[k] {
				k++
			}
			if a[i] != b[k] {
				transpositions++
			}
			k++
		}
	}

	// 计算Jaro距离
	jaro := (float64(matches)/float64(len(a)) + float64(matches)/float64(len(b)) + (float64(matches)-float64(transpositions)/2.0)/float64(matches)) / 3.0

	// 计算匹配前缀长度
	prefixLen := 0
	maxPrefixLen := int(math.Min(4, math.Min(float64(len(a)), float64(len(b)))))
	for i := 0; i < maxPrefixLen; i++ {
		if a[i] == b[i] {
			prefixLen++
		} else {
			break
		}
	}

	return jaro, prefixLen
}
