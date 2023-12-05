package uconvert

import (
	"strconv"
	"time"
	"unicode"
)

func BatchConvertStringToInt64(ins []string) ([]int64, error) {
	outs := make([]int64, 0, len(ins))
	for _, in := range ins {
		out, err := strconv.ParseInt(in, 10, 64)
		if err != nil {
			return nil, err
		}
		outs = append(outs, out)
	}
	return outs, nil
}

func BatchConvertStringToUint64(ins []string) ([]uint64, error) {
	outs := make([]uint64, 0, len(ins))
	for _, in := range ins {
		out, err := strconv.ParseUint(in, 10, 64)
		if err != nil {
			return nil, err
		}
		outs = append(outs, out)
	}
	return outs, nil
}

// ConvertTimeToUint32 受unix time的限制
func ConvertTimeToUint32(time time.Time) uint32 {
	ts := time.Unix()
	if ts < 0 {
		ts = 0
	}
	return uint32(ts)
}

func CovertStringToTimeDuration(t string) (time.Duration, error) {
	timeDuration, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}

	return timeDuration, nil
}

// ConvertCamelToSnake 将驼峰命名法字符串转换为蛇形命名法字符串。
// 驼峰命名法：MyVariableName
// 蛇形命名法：my_variable_name
func ConvertCamelToSnake(s string) string {
	var result []rune

	for i, r := range s {
		if unicode.IsUpper(r) {
			// 如果当前字符是大写字母，且前一个字符不是大写字母，
			// 并且下一个字符是小写字母或已经是字符串的最后一个字符，
			// 在前面添加下划线，并将当前字符转换为小写字母。
			if i > 0 && ((i+1 < len(s) && unicode.IsLower(rune(s[i+1]))) || (i+1 == len(s))) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			// 如果当前字符不是大写字母，则直接添加到结果中。
			result = append(result, r)
		}
	}

	return string(result)
}

func ConvertCamelWithPlaceholder(s, placeholder string) string {
	var result []rune

	for i, r := range s {
		if unicode.IsUpper(r) {
			// 如果当前字符是大写字母，且前一个字符不是大写字母，
			// 并且下一个字符是小写字母或已经是字符串的最后一个字符，
			// 在前面添加下划线，并将当前字符转换为小写字母。
			if i > 0 && ((i+1 < len(s) && unicode.IsLower(rune(s[i+1]))) || (i+1 == len(s))) {
				result = append(result, []rune(placeholder)...)
			}
			result = append(result, unicode.ToLower(r))
		} else {
			// 如果当前字符不是大写字母，则直接添加到结果中。
			result = append(result, r)
		}
	}

	return string(result)
}
