// Package distance https://gist.github.com/andrei-m/982927#gistcomment-1931258
package udistance

import (
	"context"
	"unicode/utf8"
)

// 最小长度阈值，用于小字符串优化
const minLengthThreshold = 32

// LevenshteinDistance 计算两个字符串之间的编辑距离
func LevenshteinDistance(ctx context.Context, a, b string) int {
	if len(a) == 0 {
		return utf8.RuneCountInString(b)
	}

	if len(b) == 0 {
		return utf8.RuneCountInString(a)
	}

	if a == b {
		return 0
	}

	// 如果字符串非ASCII，需要转换为 []rune
	// 这可以通过使用 utf8.RuneCountInString 避免，
	// 然后通过操作 rune 索引进行一些处理，
	// 但会导致更多的索引检查。这是一个合理的权衡。
	s1 := []rune(a)
	s2 := []rune(b)

	// 为了节省一些内存，交换字符串以减少内存占用，O(min(a,b)) 而不是 O(a)
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}
	lenS1 := len(s1)
	lenS2 := len(s2)

	// 初始化行
	var x []uint16
	if lenS1+1 > minLengthThreshold {
		x = make([]uint16, lenS1+1)
	} else {
		// 对于小字符串，我们在这里进行了一些小优化。
		// 因为具有恒定长度的切片实际上是一个数组，
		// 它不会分配内存。因此，我们可以将其重新切片到所需的长度，
		// 只要它在所需阈值以下。
		x = make([]uint16, minLengthThreshold)
		x = x[:lenS1+1]
	}

	// 我们从 1 开始，因为索引 0 已经是 0。
	for i := 1; i < len(x); i++ {
		x[i] = uint16(i)
	}

	// 做一个虚拟的索引检查，以防止下面的两次索引检查。
	// 循环内部的检查特别昂贵。
	_ = x[lenS1]
	// 填充其余部分
	for i := 1; i <= lenS2; i++ {
		prev := uint16(i)
		for j := 1; j <= lenS1; j++ {
			current := x[j-1] // 匹配
			if s2[i-1] != s1[j-1] {
				current = min(min(x[j-1]+1, prev+1), x[j]+1)
			}
			x[j-1] = prev
			prev = current
		}
		x[lenS1] = prev
	}
	return int(x[lenS1])
}

// min 返回两个 uint16 中的较小值
func min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}
