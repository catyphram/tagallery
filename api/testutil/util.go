package testutil

import "strings"

// FormatTestError creates a formatted string to pass to log functions.
// The function returns a interpolated string and the values for this string.
func FormatTestError(desc string, args map[string]interface{}) (string, []interface{}) {
	var params []interface{}
	logString := []string{desc}

	for k, v := range args {
		params = append(params, v)
		logString = append(logString, "\n", k, ": %v")
	}

	return strings.Join(logString, ""), params
}
