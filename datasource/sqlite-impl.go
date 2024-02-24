package datasource

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/okira-e/goreports/safego"
)

// SqliteDb is a wrapper for the sql.DB type.
type SqliteDb struct {
	db     *sql.DB
	dbPath string
}

// NewSqliteDb returns a new SqliteDb instance.
func NewSqliteDb(dbPath string) *SqliteDb {
	return &SqliteDb{
		db:     nil,
		dbPath: dbPath,
	}
}

// Connect connects to the database.
func (self *SqliteDb) Connect() safego.Option[error] {
	db, err := sql.Open("sqlite3", self.dbPath)
	if err != nil {
		safego.Some(err)
	}

	self.db = db

	return safego.None[error]()
}

// Disconnect disconnects from the database.
func (self *SqliteDb) Disconnect() safego.Option[error] {
	err := self.db.Close()
	if err != nil {
		safego.Some(err)
	}

	return safego.None[error]()
}

// Ping pings the database.
func (self *SqliteDb) Ping() safego.Option[error] {
	err := self.db.Ping()
	if err != nil {
		safego.Some(err)
	}

	return safego.None[error]()
}

// Query executes a query that returns rows.
func (self *SqliteDb) Query(query string, args ...any) (*sql.Rows, safego.Option[error]) {
	rows, err := self.db.Query(query, args...)
	if err != nil {
		return nil, safego.Some(err)
	}

	return rows, safego.None[error]()
}

// Exec executes a query that returns a single row.
func (self *SqliteDb) Exec(query string, args ...any) safego.Option[error] {
	_, err := self.db.Exec(query, args...)
	if err != nil {
		safego.Some(err)
	}

	return safego.None[error]()
}
