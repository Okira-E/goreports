package types

import "github.com/okira-e/goreports/safego"

// ReportAttributesForPdfGenerator is the attributes needed to generate a PDF.
// It is used in core/pdf-generators.go.
type ReportAttributesForPdfGenerator struct {
	Title  string
	Body   string
	Header safego.Option[string]
	Footer safego.Option[string]
}
