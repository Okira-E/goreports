package internalDb

import (
	"github.com/okira-e/goreports/datasource"
	"github.com/okira-e/goreports/safego"
	"github.com/okira-e/goreports/types"
)

func ListReports(internalDb *datasource.DataSource) ([]types.Report, safego.Option[error]) {
	rows, errOpt := (*internalDb).Query("SELECT * FROM reports")
	if errOpt.IsSome() {
		return []types.Report{}, safego.Some(errOpt.Unwrap())
	}

	var reportsWithNullableFields []types.ReportWithNullableFields

	for rows.Next() {
		report := types.ReportWithNullableFields{}
		err := rows.Scan(&report.ID, &report.Name, &report.Title, &report.Description, &report.Body, &report.Header, &report.Footer, &report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			return []types.Report{}, safego.Some(err)
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

	return reports, safego.None[error]()
}
