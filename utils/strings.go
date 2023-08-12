package utils

import "regexp"

// ExtractExpressions extracts expressions from a string.
func ExtractExpressions(str string, expr string) []string {
	// Extract the query from the template.
	regEx := regexp.MustCompile(expr)
	query := regEx.FindAllSubmatch([]byte(str), -1)

	result := make([]string, len(query))

	for i, q := range query {
		result[i] = string(q[1])
	}

	return result
}
