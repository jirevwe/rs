package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetDownloadUrl(t *testing.T) {
	type Args struct {
		distro  string
		version string
		os      string
		major   int64
		minor   int64
	}

	type TestData struct {
		name     string
		args     *Args
		wantUrl  string
		wantDir  string
		wantFile string
	}

	testData := []TestData{
		{
			name: "linux - 4.0.6",
			args: &Args{
				version: "4.0.6",
				os:      "linux",
				distro:  "ubuntu2004",
				major:   4,
				minor:   0,
			},
			wantUrl:  "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-4.0.6.tgz",
			wantDir:  "mongodb-linux-x86_64-4.0.6",
			wantFile: "mongodb-linux-x86_64-4.0.6.tgz",
		},
		{
			name: "linux - 4.2.0",
			args: &Args{
				version: "4.2.0",
				os:      "linux",
				distro:  "ubuntu2004",
				major:   4,
				minor:   2,
			},
			wantUrl:  "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu2004-4.2.0.tgz",
			wantDir:  "mongodb-linux-x86_64-ubuntu2004-4.2.0",
			wantFile: "mongodb-linux-x86_64-ubuntu2004-4.2.0.tgz",
		},
		{
			name: "osx - 4.0.6",
			args: &Args{
				version: "4.0.6",
				os:      "darwin",
				major:   4,
				minor:   0,
			},
			wantUrl:  "https://fastdl.mongodb.org/osx/mongodb-osx-ssl-x86_64-4.0.6.tgz",
			wantDir:  "mongodb-osx-ssl-x86_64-4.0.6",
			wantFile: "mongodb-osx-ssl-x86_64-4.0.6.tgz",
		},
		{
			name: "osx - 4.2.0",
			args: &Args{
				version: "4.2.0",
				os:      "darwin",
				major:   4,
				minor:   2,
			},
			wantUrl:  "https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-4.2.0.tgz",
			wantDir:  "mongodb-macos-x86_64-4.2.0",
			wantFile: "mongodb-macos-x86_64-4.2.0.tgz",
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			url, dir, file, err := getDownloadUrl(
				tt.args.version,
				tt.args.os,
				tt.args.distro,
				tt.args.major,
				tt.args.minor,
			)
			require.NoError(t, err)

			require.Equal(t, tt.wantUrl, url)
			require.Equal(t, tt.wantDir, dir)
			require.Equal(t, tt.wantFile, file)
		})
	}
}

func Test_ParseVersionNumber(t *testing.T) {
	type Args struct {
		version string
	}

	type TestData struct {
		name string
		args *Args

		wantMajor int64
		wantMinor int64
		wantErr   bool
	}

	testData := []TestData{
		{
			name:      "valid semvar - 4.2.0",
			args:      &Args{version: "4.2.0"},
			wantMajor: 4,
			wantMinor: 2,
		},
		{
			name:      "valid semvar - 4.0.6",
			args:      &Args{version: "4.0.6"},
			wantMajor: 4,
			wantMinor: 0,
		},
		{
			name:      "valid semvar - 4.0",
			args:      &Args{version: "4.0"},
			wantMajor: 4,
			wantMinor: 0,
			wantErr:   true,
		},
		{
			name:      "valid semvar - 4",
			args:      &Args{version: "4"},
			wantMajor: 4,
			wantMinor: 0,
			wantErr:   true,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			major, minor, err := parseVersionNumber(tt.args.version)

			if tt.wantErr {
				require.Error(t, ErrInvalidVersionFormat, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tt.wantMajor, major)
			require.Equal(t, tt.wantMinor, minor)
		})
	}
}
