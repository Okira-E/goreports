package types

import "github.com/Okira-E/goreports/safego"

type Config struct {
	DbConfig DbConfig `json:"db_config"`
}

type DbConfig struct {
	Dialect  string `json:"dialect"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// GetConnectionString returns the connection string for the database.
// If the dialect is not supported, it returns an error.
func (d *DbConfig) GetConnectionString() (string, safego.Option[string]) {
	if d.Dialect == "mysql" || d.Dialect == "mariadb" {
		return d.Username + ":" + d.Password + "@(" + d.Host + ":" + string(rune(d.Port)) + ")/" + d.Database, safego.None[string]()
	} else if d.Dialect == "postgres" {
		return "postgres://" + d.Username + ":" + d.Password + "@" + d.Host + "/" + d.Database + "?sslmode=disable", safego.None[string]()
	} else if d.Dialect == "mssql" {
		return "sqlserver://" + d.Username + ":" + d.Password + "@" + d.Host, safego.None[string]()
	}

	return "", safego.Some[string]("Invalid dialect")
}
