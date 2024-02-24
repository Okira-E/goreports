package datasource

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/okira-e/goreports/safego"
	"strings"
)

// ExternalDb is a wrapper for the sql.DB type.
type ExternalDb struct {
	db            *sql.DB
	connectionStr string
}

// NewExternalDb returns a new ExternalDb instance.
func NewExternalDb(connectionStr string) *ExternalDb {
	return &ExternalDb{
		db:            nil,
		connectionStr: connectionStr,
	}
}

// Connect connects to the database.
func (self *ExternalDb) Connect() safego.Option[error] {
	dialect := self.connectionStr[:strings.Index(self.connectionStr, ":")]

	db, err := sql.Open(dialect, self.connectionStr)
	if err != nil {
		return safego.Some(err)
	}

	self.db = db

	return safego.None[error]()
}

// Disconnect disconnects from the database.
func (self *ExternalDb) Disconnect() safego.Option[error] {
	err := self.db.Close()
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}

// Ping pings the database.
func (self *ExternalDb) Ping() safego.Option[error] {
	err := self.db.Ping()
	if err != nil {
		safego.Some(err)
	}

	return safego.None[error]()
}

// Query executes a query that returns rows.
func (self *ExternalDb) Query(query string, args ...any) (*sql.Rows, safego.Option[error]) {
	rows, err := self.db.Query(query, args...)
	if err != nil {
		return nil, safego.Some(err)
	}

	return rows, safego.None[error]()
}

// Exec executes a query that returns a single row.
func (self *ExternalDb) Exec(query string, args ...any) safego.Option[error] {
	_, err := self.db.Exec(query, args...)
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}
