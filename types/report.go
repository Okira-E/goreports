package types

import "database/sql"

type Report struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Template    string `json:"template" validate:"required"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type ReportWithNullableFields struct {
	ID          uint32         `json:"id"`
	Name        sql.NullString `json:"name" validate:"required"`
	Title       sql.NullString `json:"title" validate:"required"`
	Description sql.NullString `json:"description"`
	Template    sql.NullString `json:"template" validate:"required"`
	CreatedAt   sql.NullInt64  `json:"createdAt"`
	UpdatedAt   sql.NullInt64  `json:"updatedAt"`
}
