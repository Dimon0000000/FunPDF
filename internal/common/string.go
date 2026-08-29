package common

import "strings"

// ConcatenatingStrings remove all newlines from the string
func ConcatenatingStrings(str string) string {
	parts := strings.Split(str, "\n")
	return strings.Join(parts, "")
}
