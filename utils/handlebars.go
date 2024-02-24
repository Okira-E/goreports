package utils

import (
	"github.com/aymerick/raymond"
	"github.com/okira-e/goreports/safego"
)

// ParseHandleBars parses a handlebars template with the given template string and data.
func ParseHandleBars(template string, data map[string]any) (string, safego.Option[error]) {
	result, err := raymond.Render(template, data)
	if err != nil {
		return "", safego.Some(err)
	}

	return result, safego.None[error]()
}
