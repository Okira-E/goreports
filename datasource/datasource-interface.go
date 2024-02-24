package datasource

import (
	"database/sql"
	"github.com/okira-e/goreports/safego"
)

type DataSource interface {
	Connect() safego.Option[error]
	Disconnect() safego.Option[error]
	Ping() safego.Option[error]
	Query(string, ...any) (*sql.Rows, safego.Option[error])
	Exec(string, ...any) safego.Option[error]
}
