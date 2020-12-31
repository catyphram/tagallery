package util

import "strings"

// ContainsString checks if a given string exists in a slice.
func ContainsString(s []string, str string, cs bool) bool {
	for _, v := range s {
		if cs {
			if v == str {
				return true
			}
		} else if strings.EqualFold(v, str) {
			return true
		}
	}

	return false
}

// IntPtr creates an int and returns a pointer to it. Useful for struct inits.
func IntPtr(i int) *int {
	return &i
}

// StringPtr creates a string and returns a pointer to it. Useful for struct inits.
func StringPtr(s string) *string {
	return &s
}
