package core

import (
	"bytes"
	"github.com/Okira-E/goreports/safego"
	"github.com/Okira-E/goreports/types"
	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"strings"
)

// GeneratePDFFromHtml generates a PDF from HTML.
// It takes in HTML as a string and printing options of type types.PrintingOptions
// It returns a buffer and an error.
func GeneratePDFFromHtml(reportParams types.ReportGeneratorParams, printingOptions types.PrintingOptions) (*bytes.Buffer, safego.Option[error]) {
	// Create new PDF generator
	pdfGenerator, err := pdf.NewPDFGenerator()
	if err != nil {
		return &bytes.Buffer{}, safego.Some[error](err)
	}

	// Set global options
	pdfGenerator.Dpi.Set(300)
	pdfGenerator.MarginLeft.Set(0)
	pdfGenerator.MarginRight.Set(0)
	pdfGenerator.MarginTop.Set(0)
	pdfGenerator.MarginBottom.Set(0)

	pdfGenerator.PageSize.Set(printingOptions.PaperSize)

	if printingOptions.Landscape {
		pdfGenerator.Orientation.Set("Landscape")
	} else {
		pdfGenerator.Orientation.Set("Portrait")
	}

	pdfGenerator.Grayscale.Set(true)

	pdfGenerator.Title.Set(reportParams.Title)
	pdfGenerator.AddPage(pdf.NewPageReader(strings.NewReader(reportParams.Html)))

	// Create PDF document in internal buffer
	err = pdfGenerator.Create()
	if err != nil {
		return &bytes.Buffer{}, safego.Some[error](err)
	}

	// Send file in response
	return pdfGenerator.Buffer(), safego.None[error]()
}
