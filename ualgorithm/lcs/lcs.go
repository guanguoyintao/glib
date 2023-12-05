package ulcs

import "unicode"

func Lcs(a, b string) (string, int) {
	r1, r2 := []rune(a), []rune(b)
	m, n := len(r1), len(r2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i, c1 := range r1 {
		for j, c2 := range r2 {
			if unicode.ToLower(c1) == unicode.ToLower(c2) {
				dp[i+1][j+1] = dp[i][j] + 1
			} else {
				dp[i+1][j+1] = max(dp[i][j+1], dp[i+1][j])
			}
		}
	}
	// 构建最长公共子序列的字符串
	lcsString := make([]rune, dp[m][n])
	i, j, k := m, n, dp[m][n]
	for i > 0 && j > 0 {
		if unicode.ToLower(r1[i-1]) == unicode.ToLower(r2[j-1]) {
			lcsString[k-1] = r1[i-1]
			i--
			j--
			k--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	return string(lcsString), dp[m][n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
