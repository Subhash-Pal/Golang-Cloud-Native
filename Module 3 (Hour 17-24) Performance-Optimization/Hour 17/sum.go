package benchmark

import "strings"

func SumInts(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

func JoinWithPlus(parts []string) string {
	var builder strings.Builder
	for i, part := range parts {
		if i > 0 {
			builder.WriteString(" + ")
		}
		builder.WriteString(part)
	}
	return builder.String()
}
