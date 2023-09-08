package routers

import (
	"database/sql"
	"github.com/Okira-E/goreports/core"
	"github.com/Okira-E/goreports/safego"
	"github.com/Okira-E/goreports/types"
	"github.com/Okira-E/goreports/utils"
	"github.com/gofiber/fiber/v2"
)

// ReportsRouter sets up the routes for reports.
// This function is called from server/routers/index.go.
func ReportsRouter(app *fiber.App) {
	const controllerName = "/report"

	app.Get(controllerName+"/list", func(ctx *fiber.Ctx) error {
		rows, errOpt := (*InternalDb).Query("SELECT * FROM reports")
		if errOpt.IsSome() {
			return ctx.Status(500).SendString(errOpt.Unwrap().Error())
		}

		reportsWithNullableFields := []types.ReportWithNullableFields{}

		for rows.Next() {
			report := types.ReportWithNullableFields{}
			err := rows.Scan(&report.ID, &report.Name, &report.Title, &report.Description, &report.Body, &report.CreatedAt, &report.UpdatedAt)
			if err != nil {
				return ctx.Status(500).SendString(err.Error())
			}
			reportsWithNullableFields = append(reportsWithNullableFields, report)
		}

		// Convert the nullable fields to non-nullable fields.
		reports := []types.Report{}
		for _, report := range reportsWithNullableFields {
			reports = append(reports, types.Report{
				ID:          report.ID,
				Name:        report.Name.String,
				Title:       report.Title.String,
				Description: report.Description.String,
				Body:        report.Body.String,
				CreatedAt:   report.CreatedAt.Int64,
				UpdatedAt:   report.UpdatedAt.Int64,
			})
		}

		return ctx.Status(200).JSON(reports)
	})

	// ------------------------------------------------------------

	app.Post(controllerName+"/save", func(ctx *fiber.Ctx) error {
		report := types.Report{}

		// Parse the request body.
		utils.ParseRequestBody(ctx, &report)

		// Validate the request body.
		if report.Name == "" {
			return ctx.Status(400).SendString("The report name is required.")
		}
		if report.Title == "" {
			return ctx.Status(400).SendString("The report title is required.")
		}
		if report.Body == "" {
			return ctx.Status(400).SendString("The report body is required.")
		}

		// Set the timestamps.
		report.CreatedAt = utils.GetTimestamp()
		report.UpdatedAt = 0

		// Save the report.
		var header, footer sql.NullString

		header = sql.NullString{
			String: report.Header,
			Valid:  true,
		}
		footer = sql.NullString{
			String: report.Footer,
			Valid:  true,
		}

		if len(header.String) == 0 {
			header = sql.NullString{}
		}
		if len(footer.String) == 0 {
			footer = sql.NullString{}
		}
		errOpt := (*InternalDb).Exec("INSERT INTO reports (name, title, description, body, header, footer, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", report.Name, report.Title, report.Description, report.Body, header, footer, report.CreatedAt, report.UpdatedAt)
		if errOpt.IsSome() {
			return ctx.Status(400).SendString(errOpt.Unwrap().Error())
		}

		// Return a response.
		return ctx.Status(201).JSON(map[string]string{
			"message":    "Report saved successfully.",
			"reportName": report.Name,
		})
	})

	// ------------------------------------------------------------

	app.Post(controllerName+"/render", func(ctx *fiber.Ctx) error {
		// Define the request renderBody.
		var renderBody struct {
			ReportName      string                `json:"reportName"`
			Params          map[string]any        `json:"params"`
			PrintingOptions types.PrintingOptions `json:"printingOptions"`
		}

		// Parse the request renderBody.
		utils.ParseRequestBody(ctx, &renderBody)

		// Validate the request renderBody.
		if renderBody.ReportName == "" {
			return ctx.Status(400).SendString("The report name is required.")
		}

		if renderBody.Params == nil {
			renderBody.Params = make(map[string]any)
		}

		// Get the report from the database.
		report := types.ReportWithNullableFields{}
		rows, errOpt := (*InternalDb).Query("SELECT * FROM reports WHERE name = ?", renderBody.ReportName)
		if errOpt.IsSome() {
			return ctx.Status(400).SendString(errOpt.Unwrap().Error())
		}
		// Extract the report from the rows.
		emptyResult := true
		for rows.Next() {
			emptyResult = false
			err := rows.Scan(&report.ID, &report.Name, &report.Title, &report.Description, &report.Body, &report.Header, &report.Footer, &report.CreatedAt, &report.UpdatedAt)
			if err != nil {
				return ctx.Status(500).SendString(err.Error())
			}
		}
		// If the report was not found, return a 404 response.
		if emptyResult {
			return ctx.Status(404).SendString("report was not found.")
		}

		handlebarsTemplate, queries, errMsgOpt := core.ParseTemplate(report.Body.String, renderBody.Params, ExternalDb)
		if errMsgOpt.IsSome() {
			return ctx.Status(400).SendString(errMsgOpt.Unwrap())
		}

		// Parse the template in handlebars.
		compiledTemplate, errOpt := utils.ParseHandleBars(handlebarsTemplate, queries)
		if errOpt.IsSome() {
			return ctx.Status(500).SendString(errOpt.Unwrap().Error())
		}

		// Generate the document
		header, footer := safego.None[string](), safego.None[string]()

		if report.Header.Valid {
			header = safego.Some(report.Header.String)
		}
		if report.Footer.Valid {
			footer = safego.Some(report.Footer.String)
		}

		reportGeneratorParams := types.ReportAttributesForPdfGenerator{
			Title:  report.Title.String,
			Body:   compiledTemplate,
			Header: header,
			Footer: footer,
		}

		generatedPDFBuffer, errOpt := core.GeneratePDFFromHtml(reportGeneratorParams, renderBody.PrintingOptions)
		if errOpt.IsSome() {
			return ctx.Status(500).SendString(errOpt.Unwrap().Error())
		}

		// Return a response.
		return ctx.Status(200).Send(generatedPDFBuffer.Bytes())
	})

	// ------------------------------------------------------------

	app.Delete(controllerName+"/delete", func(ctx *fiber.Ctx) error {
		// Parse the request body.
		var body struct {
			ReportName string `json:"reportName"`
		}

		utils.ParseRequestBody(ctx, &body)

		// Validate the request body.
		if body.ReportName == "" {
			return ctx.Status(400).SendString("reportName is required.")
		}

		// Delete the report.
		errOpt := (*InternalDb).Exec("DELETE FROM reports WHERE name = ?", body.ReportName)
		if errOpt.IsSome() {
			return ctx.Status(400).SendString(errOpt.Unwrap().Error())
		}

		// Return a response.
		return ctx.Status(200).JSON(map[string]string{
			"message": "Report " + body.ReportName + " deleted successfully.",
		})
	})
}
