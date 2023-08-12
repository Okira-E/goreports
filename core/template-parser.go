package core

import (
	"encoding/json"
	"fmt"
	"github.com/Okira-E/goreports/datasource"
	"github.com/Okira-E/goreports/safego"
	"github.com/Okira-E/goreports/utils"
	"strconv"
	"strings"
)

// ParseTemplate takes in a template in Handlebars format with custom directives and returns the template string with
// the directives replaced with the actual values. It also returns a map of the queries and their results.
// The template can contain parameters and queries. Parameters are evaluated first, then queries.
// If a parameter is not provided, the function returns an error message.
func ParseTemplate(template string, params map[string]any, ds *datasource.DataSource) (string, map[string]any, safego.Option[string]) {
	//// Parameters evaluation ////

	// Extract every [P[...]] expression.
	regexExpr := `\[P\[(.+?)\]\]`
	parameters := utils.ExtractExpressions(template, regexExpr)

	for _, parameter := range parameters {
		// Check if the parameter is provided.
		if _, ok := params[parameter]; !ok {
			return "", map[string]any{}, safego.Some(fmt.Sprintf("Parameter %s is not provided.", parameter))
		}

		// Replace the parameter in the original template with the generated name.
		template = strings.ReplaceAll(template, fmt.Sprintf("[P[%s]]", parameter), fmt.Sprintf("%v", params[parameter]))
	}

	//// Queries evaluation ////

	// Extract every [Q[...]] expression.
	regexExpr = `\[Q\[(.+?)\]\]`
	res := utils.ExtractExpressions(template, regexExpr)

	// Evaluate every query.
	queries := map[string]any{}
	queryCounter := 0
	for _, query := range res {
		rows, errOpt := (*ds).Query(query)
		if errOpt.IsSome() {
			return "", map[string]any{}, safego.Some(errOpt.Unwrap().Error())
		}

		// The result for a single query in the template.
		// Can be a single row or multiple rows.
		queryResults := utils.JsonifyQueryData(rows)

		// Unmarshal the query result(s).
		if len(queryResults) == 1 {
			var data map[string]any

			err := json.Unmarshal([]byte(queryResults[0]), &data)
			if err != nil {
				return "", map[string]any{}, safego.Some(err.Error())
			}

			// Get the first (and only) value from the map.
			var actualResult string
			for _, value := range data {
				actualResult = fmt.Sprintf("%v", value)
				break
			}

			// Replace the query in the original template with the generated name.
			template = strings.ReplaceAll(template, fmt.Sprintf("[Q[%s]]", query), actualResult)
		} else {
			// If the query returned multiple rows.
			// Initialize the outer map.
			queryKeyName := "data_" + strconv.Itoa(queryCounter)
			queries[queryKeyName] = []map[string]any{}

			// Iterate through the query results.
			// I can't use `i` because of the `continue` statement.
			rowCounter := 0
			for _, queryResult := range queryResults {
				var data map[string]any
				// Skip the line if it's a comma.
				if strings.Fields(queryResult)[0] == "," {
					continue
				}

				err := json.Unmarshal([]byte(queryResult), &data)
				if err != nil {
					return "", map[string]any{}, safego.Some(err.Error())
				}

				// Append the data to the outer map.
				queries[queryKeyName] = append(queries[queryKeyName].([]map[string]any), data)

				rowCounter++
			}

			// Replace the query in the original template with the generated name.
			theQueryInTheTemplate := fmt.Sprintf("[Q[%s]]", query)
			template = strings.ReplaceAll(template, theQueryInTheTemplate, queryKeyName)
		}

		queryCounter++
	}

	return template, queries, safego.None[string]()
}
