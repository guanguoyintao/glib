// Package uexp for expression
package uexp

// Iif Inline if 三元表达式
func Iif(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
