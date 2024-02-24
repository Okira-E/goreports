package utils

import (
	"database/sql"
	"encoding/json"
	"github.com/okira-e/goreports/safego"
	"os"
	"strconv"
	"strings"
)

// JsonifyQueryData converts the data from a sql query to json.
func JsonifyQueryData(rows *sql.Rows) []string {
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	dest := make([]interface{}, len(columns))

	desRefs := make([]interface{}, len(dest))
	for i := range dest {
		desRefs[i] = &dest[i]
	}

	c := 0
	results := make(map[string]interface{})
	data := []string{}

	for rows.Next() {
		if c > 0 {
			data = append(data, ",")
		}

		err = rows.Scan(desRefs...)
		if err != nil {
			panic(err.Error())
		}

		for i, value := range dest {
			switch value.(type) {
			case nil:
				results[columns[i]] = nil
			case []byte:
				s := string(value.([]byte))
				x, err := strconv.Atoi(s)

				if err != nil {
					results[columns[i]] = s
				} else {
					results[columns[i]] = x
				}
			default:
				results[columns[i]] = value
			}
		}

		b, _ := json.Marshal(results)
		data = append(data, strings.TrimSpace(string(b)))
		c++
	}

	return data
}

// ReadJSONFile reads a JSON file and unmarshals it into a value reference.
func ReadJSONFile(filePath string, valRef any) safego.Option[error] {
	fileContent, err := os.ReadFile(filePath)

	err = json.Unmarshal(fileContent, valRef)
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}

// WriteToJSONFile writes a value to a JSON file.
func WriteToJSONFile(filePath string, content any) safego.Option[error] {
	file, err := os.Create(filePath)
	defer file.Close()

	fileContent, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		return safego.Some(err)
	}

	_, err = file.Write(fileContent)
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}
