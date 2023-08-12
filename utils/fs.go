package utils

import (
	"github.com/Okira-E/goreports/safego"
	"os"
)

// CreateDir creates a directory.
func CreateDir(dir string) safego.Option[error] {
	// 0755 is the permission for the directory.
	// The difference between 0755 and 0777 is that 0755 is read and write for the owner and read for everyone else.
	err := os.Mkdir(dir, 0755)
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}

// CreateFile creates a file.
func CreateFile(path string) safego.Option[error] {
	var err error

	f, err := os.Create(path)
	if err != nil {
		return safego.Some(err)
	}
	defer func() {
		err = f.Close()
	}()
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}

// WriteFile writes data to a file.
func WriteFile(path string, data []byte) safego.Option[error] {
	var err error

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return safego.Some(err)
	}
	defer func() {
		err = f.Close()
	}()
	if err != nil {
		return safego.Some(err)
	}

	if _, err := f.Write(data); err != nil {
		return safego.Some(err)
	}

	if err := f.Sync(); err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}
