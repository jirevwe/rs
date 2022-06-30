package main

import (
	"embed"
	"errors"
	"strings"
)

//go:embed VERSION
var f embed.FS

func ReadVersion() ([]byte, error) {
	data, err := f.ReadFile("VERSION")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetVersion() string {
	v := "0.1.0"

	f, err := ReadVersion()
	if err != nil {
		return v
	}

	v = strings.TrimSuffix(string(f), "\n")
	return v
}

var (
	dbPath                  = ""
	distro                  = ""
	version                 = ""
	base                    = "https://fastdl.mongodb.org"
	ErrInvalidVersionFormat = errors.New("please pass a valid mongodb version; version must be in x.x.x format")
)
