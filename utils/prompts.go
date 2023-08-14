package utils

import (
	"fmt"
	"github.com/Okira-E/goreports/safego"
	"github.com/Okira-E/goreports/types"
	"github.com/Okira-E/goreports/vars"
	"github.com/manifoldco/promptui"
	"strconv"
)

// PromptForDbConfig prompts the user for their database configuration.
// It returns a DbConfig and an error.
func PromptForDbConfig(dbConfig *types.DbConfig) safego.Option[error] {
	// Dialect
	if dbConfig.Dialect == "" {
		dialectPrmpt := promptui.Select{
			Label: "Dialect",
			Items: vars.SupportedDatabases,
		}
		_, dialect, err := dialectPrmpt.Run()
		if err != nil {
			return safego.Some(err)
		}

		dbConfig.Dialect = dialect
	}

	// Host
	if dbConfig.Host == "" {
		HostPrmpt := promptui.Prompt{
			Label: "Host",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("host cannot be empty")
				}

				return nil
			},
		}
		host, err := HostPrmpt.Run()
		if err != nil {
			return safego.Some(err)
		}

		dbConfig.Host = host
	}

	// Port
	if dbConfig.Port == 0 {
		PortPrmpt := promptui.Prompt{
			Label: "Port",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("port cannot be empty")
				}
				// Check if the port is a number.
				intPort, err := strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("port must be a number")
				}

				if intPort < 0 || intPort > 65535 {
					return fmt.Errorf("port must be between 0 and 65535")
				}

				return nil
			},
		}
		stringPort, err := PortPrmpt.Run()
		if err != nil {
			return safego.Some(err)
		}
		port, err := strconv.Atoi(stringPort)
		if err != nil {
			return safego.Some(err)
		}

		dbConfig.Port = port
	}

	// User
	if dbConfig.Username == "" {
		userPrmpt := promptui.Prompt{
			Label: "User",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("user cannot be empty")
				}

				return nil
			},
		}
		user, err := userPrmpt.Run()
		if err != nil {
			return safego.Some(err)
		}

		dbConfig.Username = user
	}

	// Password
	if dbConfig.Password == "" {
		passwordPrmpt := promptui.Prompt{
			Label: "Password",
			Mask:  '*',
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("password cannot be empty")
				}

				return nil
			},
		}
		password, err := passwordPrmpt.Run()
		if err != nil {
			return safego.Some(err)
		}

		dbConfig.Password = password
	}

	// Database name
	if dbConfig.Database == "" {
		databasePrmpt := promptui.Prompt{
			Label: "Database name",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("database name cannot be empty")
				}

				return nil
			},
		}
		databaseName, err := databasePrmpt.Run()
		if err != nil {
			return safego.Some(err)
		}

		dbConfig.Database = databaseName
	}

	return safego.None[error]()
}
