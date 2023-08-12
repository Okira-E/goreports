package cmd

import (
	"github.com/Okira-E/goreports/internalDb"
	"github.com/Okira-E/goreports/server"
	"github.com/Okira-E/goreports/types"
	"github.com/Okira-E/goreports/utils"
	"github.com/spf13/cobra"
	"log"
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
		dbConfig, errOpt := utils.PromptForDbConfig()
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

var startAppCmd = &cobra.Command{
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

func Execute() {
	rootCmd.AddCommand(
		versionCmd,
		runInit,
		startAppCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
