package utils

import "strings"

func InSlice(arr []string, search string, strict bool) bool {
	element := ""
	for i := 0; i < len(arr); i++ {
		element = arr[i]
		if strict {
			if strings.EqualFold(element, search) {
				return true
			}
		} else {
			if element == search {
				return true
			}
		}
	}
	return false
}
func RepeatString(value string, n int) []string {
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		arr[i] = value
	}
	return arr
}
