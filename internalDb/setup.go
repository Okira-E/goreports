package internalDb

import (
	"database/sql"
	"github.com/Okira-E/goreports/safego"
	"github.com/Okira-E/goreports/utils"
	_ "github.com/mattn/go-sqlite3"
)

// BuildInternalDb builds the internal database for the first time.
func BuildInternalDb(dataDir string) safego.Option[error] {
	errOpt := utils.CreateFile(dataDir + "/internal.db")
	if errOpt.IsSome() {
		return errOpt
	}

	// Connect to the database.
	connection, err := sql.Open("sqlite3", dataDir+"/internal.db")
	if err != nil {
		return safego.Some(err)
	}

	// Create the tables.
	errOpt = createInternalDbTables(connection)
	if errOpt.IsSome() {
		return errOpt
	}

	return safego.None[error]()
}

// createInternalDbTables creates all internal database tables.
func createInternalDbTables(connection *sql.DB) safego.Option[error] {
	const (
		createReportsTableSQL = `
			CREATE TABLE IF NOT EXISTS reports (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name VARCHAR(255) NOT NULL UNIQUE,
				title VARCHAR(255) NOT NULL,
				description TEXT NULL,
				body TEXT NOT NULL,
				header TEXT NULL,
				footer TEXT NULL,
				created_at TEXT NOT NULL,
				updated_at TEXT NOT NULL
        );`
	)

	// Create the reports table.
	_, err := connection.Exec(createReportsTableSQL)
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}
