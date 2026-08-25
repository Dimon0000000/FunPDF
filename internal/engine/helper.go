package engine

import "strings"

func ConcatenatingStrings(str string) string {
	parts := strings.Split(str, "\n")
	return strings.Join(parts, "")
}
