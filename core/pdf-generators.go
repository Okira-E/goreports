package core

import (
	"bytes"
	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/okira-e/goreports/safego"
	"github.com/okira-e/goreports/types"
	"os"
	"strings"
)

// GeneratePDFFromHtml generates a PDF from HTML.
// It takes in HTML as a string and printing options of type types.PrintingOptions
// It returns a buffer and an error.
func GeneratePDFFromHtml(reportParams types.ReportAttributesForPdfGenerator, printingOptions types.PrintingOptions) (*bytes.Buffer, safego.Option[error]) {
	// Create new PDF generator
	pdfGenerator, err := pdf.NewPDFGenerator()
	if err != nil {
		return &bytes.Buffer{}, safego.Some(err)
	}

	pdfGenerator.SetStderr(&bytes.Buffer{})

	// Set global options
	pdfGenerator.Dpi.Set(300)
	pdfGenerator.MarginLeft.Set(uint(printingOptions.MarginLeft))
	pdfGenerator.MarginRight.Set(uint(printingOptions.MarginRight))
	pdfGenerator.MarginTop.Set(uint(printingOptions.MarginTop))
	pdfGenerator.MarginBottom.Set(uint(printingOptions.MarginBottom))

	pdfGenerator.PageSize.Set(printingOptions.PaperSize)

	if printingOptions.Landscape {
		pdfGenerator.Orientation.Set("Landscape")
	} else {
		pdfGenerator.Orientation.Set("Portrait")
	}

	pdfGenerator.Grayscale.Set(true)

	pdfGenerator.Title.Set(reportParams.Title)

	page := pdf.NewPageReader(strings.NewReader(reportParams.Body))

	// Setup repeating header if provided.
	if reportParams.Header.IsSome() && (!printingOptions.PageNumbers.Enabled || !strings.Contains(printingOptions.PageNumbers.Position, "top")) {
		err = os.WriteFile("./core/header_temp.html", []byte("<!doctype html>"+reportParams.Header.Unwrap()), 0644)
		if err != nil {
			return &bytes.Buffer{}, safego.Some(err)
		}

		page.HeaderHTML.Set("file:///" + os.Getenv("PWD") + "/core/header_temp.html")
	}
	// Setup repeating footer if provided and page numbers are not enabled.
	if reportParams.Footer.IsSome() && (!printingOptions.PageNumbers.Enabled || !strings.Contains(printingOptions.PageNumbers.Position, "bottom")) {
		err = os.WriteFile("./core/footer_temp.html", []byte("<!doctype html>"+reportParams.Footer.Unwrap()), 0644)
		if err != nil {
			return &bytes.Buffer{}, safego.Some(err)
		}

		page.FooterHTML.Set("file:///" + os.Getenv("PWD") + "/core/footer_temp.html")
	}

	// Setup page numbers if enabled.
	if printingOptions.PageNumbers.Enabled {
		if printingOptions.PageNumbers.Position == "top-left" {
			pdfGenerator.MarginTop.Set(7)
			page.HeaderLeft.Set("[page]")
		} else if printingOptions.PageNumbers.Position == "top-center" {
			pdfGenerator.MarginTop.Set(7)
			page.HeaderCenter.Set("[page]")
		} else if printingOptions.PageNumbers.Position == "top-right" {
			pdfGenerator.MarginTop.Set(7)
			page.HeaderRight.Set("[page]")
		} else if printingOptions.PageNumbers.Position == "bottom-left" {
			pdfGenerator.MarginBottom.Set(7)
			page.FooterLeft.Set("[page]")
		} else if printingOptions.PageNumbers.Position == "bottom-center" {
			pdfGenerator.MarginBottom.Set(7)
			page.FooterCenter.Set("[page]")
		} else if printingOptions.PageNumbers.Position == "bottom-right" {
			pdfGenerator.MarginBottom.Set(7)
			page.FooterRight.Set("[page]")
		}
	}

	// Create PDF document in internal buffer
	pdfGenerator.AddPage(page)
	err = pdfGenerator.Create()
	if err != nil {
		return &bytes.Buffer{}, safego.Some(err)
	}

	os.Remove("./core/header_temp.html")
	os.Remove("./core/footer_temp.html")

	// Send file in response
	return pdfGenerator.Buffer(), safego.None[error]()
}
