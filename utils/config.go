package utils

import (
	"errors"
	"github.com/Okira-E/goreports/safego"
	"github.com/Okira-E/goreports/types"
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigData extracts the database connection string from the config file.
func GetConfigData() (types.Config, safego.Option[error]) {
	var config types.Config

	filePath, errOpt := GetConfigFilePathBasedOnOS()
	if errOpt.IsSome() {
		return types.Config{}, errOpt
	}

	errOpt = ReadJSONFile(filePath, &config)
	if errOpt.IsSome() {
		return types.Config{}, errOpt
	}

	return config, safego.None[error]()
}

// GetConfigFilePathBasedOnOS returns the config file path based on the OS.
func GetConfigFilePathBasedOnOS() (string, safego.Option[error]) {
	var osUserName string

	if runtime.GOOS == "windows" {
		osUserName = os.Getenv("USERNAME")
		if osUserName == "" {
			return "", safego.Some(errors.New("USERNAME environment variable is not set"))
		}

		partialConfigPath := filepath.Join("AppData", "Roaming", "goreports", "config.json")
		return filepath.Join("C:\\Users", osUserName, partialConfigPath), safego.None[error]()
	} else if runtime.GOOS == "darwin" {
		osUserName = os.Getenv("USER")
		if osUserName == "" {
			return "", safego.Some(errors.New("USER environment variable is not set"))
		}

		partialConfigPath := filepath.Join("Library", "Application Support", "goreports", "config.json")
		return filepath.Join("/Users", osUserName, partialConfigPath), safego.None[error]()
	} else if runtime.GOOS == "linux" {
		osHomePath := os.Getenv("HOME")

		partialConfigPath := filepath.Join(".config", "goreports", "config.json")
		return filepath.Join(osHomePath, partialConfigPath), safego.None[error]()
	} else {
		err := errors.New("unsupported OS")
		return "", safego.Some(err)
	}
}

func GetDataDirBasedOnOS() (string, safego.Option[error]) {
	var osUserName string

	if runtime.GOOS == "windows" {
		osUserName = os.Getenv("USERNAME")
		if osUserName == "" {
			return "", safego.Some(errors.New("USERNAME environment variable is not set"))
		}

		partialConfigPath := filepath.Join("AppData", "Roaming", "goreports", "data", "")
		return filepath.Join("C:\\Users", osUserName, partialConfigPath), safego.None[error]()

	} else if runtime.GOOS == "darwin" {
		osUserName = os.Getenv("USER")
		if osUserName == "" {
			return "", safego.Some(errors.New("USER environment variable is not set"))
		}

		partialConfigPath := filepath.Join("Library", "Application Support", "goreports", "data", "")
		return filepath.Join("/Users", osUserName, partialConfigPath), safego.None[error]()

	} else if runtime.GOOS == "linux" {
		osHomePath := os.Getenv("HOME")

		partialConfigPath := filepath.Join(".config", "goreports", "data", "")
		return filepath.Join(osHomePath, partialConfigPath), safego.None[error]()
	} else {
		err := errors.New("unsupported OS")
		return "", safego.Some(err)
	}
}

// DoesConfigFileExists checks if the config file exists.
// It returns true if the file exists, false otherwise.
func DoesConfigFileExists() (bool, safego.Option[error]) {
	filePath, errOpt := GetConfigFilePathBasedOnOS()
	if errOpt.IsSome() {
		return false, errOpt
	}

	// Check if the file exists.
	if _, errOpt := os.Stat(filePath); os.IsNotExist(errOpt) {
		return false, safego.None[error]()
	}

	return true, safego.None[error]()
}

// CreateConfigFile creates the config file.
func CreateConfigFile() safego.Option[error] {
	filePath, errOpt := GetConfigFilePathBasedOnOS()
	if errOpt.IsSome() {
		return errOpt
	}

	// Create the directory.
	dirPath := filePath[:len(filePath)-len("/config.json")]
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return safego.Some(err)
	}

	// Create the file inside the directory.
	file, err := os.Create(filePath)
	defer file.Close()

	_, err = file.Write([]byte(`{}`))
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}
