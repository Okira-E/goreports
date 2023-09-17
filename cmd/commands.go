package cmd

import (
	"github.com/Okira-E/goreports/datasource"
	"github.com/Okira-E/goreports/internalDb"
	"github.com/Okira-E/goreports/server"
	"github.com/Okira-E/goreports/types"
	"github.com/Okira-E/goreports/utils"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var rootCmd = &cobra.Command{
	Use:   "goreports",
	Short: "GoReports is a report generation tool",
	Long:  `GoReports is a report generation tool that allows you to build and generate dynamic reports in many formats.`,
}
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Lists the current version of GoReports",
	Long:  `Lists the current version of GoReports`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Log("0.5.0")
	},
}

var runInit = &cobra.Command{
	Use:   "init",
	Short: "Initializes GoReports on your system",
	Long:  `Initializes GoReports on your system`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the database information from the flags, if they exist.
		dbDialect, err := cmd.Flags().GetString("db-dialect")
		if err != nil {
			log.Fatalf("error while getting the db-dialect flag: %v", err)
		}
		dbUser, err := cmd.Flags().GetString("db-username")
		if err != nil {
			log.Fatalf("error while getting the db-username flag: %v", err)
		}
		dbPassword, err := cmd.Flags().GetString("db-password")
		if err != nil {
			log.Fatalf("error while getting the db-password flag: %v", err)
		}
		dbHost, err := cmd.Flags().GetString("db-host")
		if err != nil {
			log.Fatalf("error while getting the db-host flag: %v", err)
		}
		dbPortStr, err := cmd.Flags().GetString("db-port")
		if err != nil {
			log.Fatalf("error while getting the db-port flag: %v", err)
		}
		dbPort, err := strconv.Atoi(dbPortStr)
		if err != nil {
			dbPort = 0
		}
		dbName, err := cmd.Flags().GetString("db-name")
		if err != nil {
			log.Fatalf("error while getting the db-name flag: %v", err)
		}

		dbConfig := types.DbConfig{
			Dialect:  dbDialect,
			Host:     dbHost,
			Port:     dbPort,
			Username: dbUser,
			Password: dbPassword,
			Database: dbName,
		}

		// Create the config files if they don't exist.
		found, errOpt := utils.DoesConfigFileExists()
		if errOpt.IsSome() {
			log.Fatalf("error while checking if the config file exists: %v", errOpt.Unwrap())
		}

		if found {
			utils.Log("The config file already exists.")
			return
		}

		errOpt = utils.CreateConfigFile()
		if errOpt.IsSome() {
			log.Fatalf("error while creating the config file: %v", errOpt.Unwrap())
		}

		// Prompt the user to enter their database information.
		errOpt = utils.PromptForDbConfig(&dbConfig)
		if errOpt.IsSome() {
			log.Fatalf("error while prompting for the database config: %v", errOpt.Unwrap())
		}

		// Create the configData.json file data.
		configData := types.Config{
			DbConfig: dbConfig,
		}

		// Write the config.json file.
		configPath, errOpt := utils.GetConfigFilePathBasedOnOS()
		if errOpt.IsSome() {
			log.Fatalf("error while getting the config file path: %v", errOpt.Unwrap())
		}

		errOpt = utils.WriteToJSONFile(configPath, configData)

		dataDir, errOpt := utils.GetDataDirBasedOnOS()
		if errOpt.IsSome() {
			log.Fatalf("error while getting the data directory: %v", errOpt.Unwrap())
		}

		// Create the database directory.
		errOpt = utils.CreateDir(dataDir)
		if errOpt.IsSome() {
			log.Fatalf("error while creating the data directory: %v", errOpt.Unwrap())
		}
		// Create the database file.
		errOpt = utils.CreateFile(dataDir + "/internal.db")
		if errOpt.IsSome() {
			log.Fatalf("error while creating the database file: %v", errOpt.Unwrap())
		}

		// Build the database.
		errOpt = internalDb.BuildInternalDb(dataDir)
		if errOpt.IsSome() {
			log.Fatalf("error while building the database: %v", errOpt.Unwrap())
		}

		utils.Log("Initialized the config directory at " + configPath)
	},
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the program",
	Long:  "Starts accepting requests and generating reports",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the config file exists.
		found, errOpt := utils.DoesConfigFileExists()
		if errOpt.IsSome() {
			log.Fatalf("error while checking if the config file exists: %v", errOpt.Unwrap())
		}

		if !found {
			utils.Log("The config file does not exist. Running the `init` command...")
			runInit.Run(cmd, args)
		}

		// Start the server.
		utils.Log("Starting the server...")
		server.StartServer()
	},
}

var listReportsCmd = &cobra.Command{
	Use:   "list-reports",
	Short: "List all reports",
	Long:  "Lists all reports",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the config file exists.
		found, errOpt := utils.DoesConfigFileExists()
		if errOpt.IsSome() {
			log.Fatalf("error while checking if the config file exists: %v", errOpt.Unwrap())
		}

		if !found {
			utils.Log("The config file does not exist. Running the `init` command...")
			runInit.Run(cmd, args)
		}

		// Establish a connection with internal database.
		dataDir, errOpt := utils.GetDataDirBasedOnOS()
		if errOpt.IsSome() {
			log.Fatalf("error while getting the data directory: %v", errOpt.Unwrap())
		}

		var internalDbConn datasource.DataSource
		internalDbConn = datasource.NewSqliteDb(dataDir + "/internal.db")

		errOpt = internalDbConn.Connect()
		if errOpt.IsSome() {
			log.Fatalf("error while connecting to the database: %v", errOpt.Unwrap())
		}
		defer internalDbConn.Disconnect()

		// List all reports.
		utils.Log("Listing all reports...")
		reports, errOpt := internalDb.ListReports(&internalDbConn)
		if errOpt.IsSome() {
			log.Fatalf("error while listing all reports: %v", errOpt.Unwrap())
		}

		for i, report := range reports {
			utils.Log(strconv.Itoa(i+1) + ": " + report.Name)
		}
	},
}

func Execute() {
	// Add the flags to the runInit command.
	runInit.Flags().StringP("db-dialect", "d", "", "The dialect of the database")
	runInit.Flags().StringP("db-username", "u", "", "The username for the database")
	runInit.Flags().StringP("db-password", "p", "", "The password for the database")
	runInit.Flags().StringP("db-host", "H", "", "The host for the database")
	runInit.Flags().StringP("db-port", "P", "", "The port for the database")
	runInit.Flags().StringP("db-name", "D", "", "The database name")

	// Add the commands to the root command.
	rootCmd.AddCommand(
		versionCmd,
		runInit,
		startServerCmd,
		listReportsCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
