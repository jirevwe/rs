package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	distro                  = ""
	version                 = ""
	base                    = "https://fastdl.mongodb.org"
	ErrInvalidVersionFormat = errors.New("please pass a valid mongodb version; version must be in x.x.x format")
)

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVar(&distro, "distro", "ubuntu1804", "specify the linux distro")
	downloadCmd.Flags().StringVar(&version, "version", "4.2.21", "specify the mongodb version")
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads and configures a mongodb version",
	Long:  "Downloads and configures a mongodb version",
	RunE: func(cmd *cobra.Command, args []string) error {
		major, minor, err := parseVersionNumber(version)
		if err != nil {
			return err
		}

		url, dir, file, err := getDownloadUrl(version, runtime.GOOS, distro, major, minor)
		if err != nil {
			return err
		}

		err = downloadFile(file, url)
		if err != nil {
			return err
		}

		err = extract(file, dir, version)
		if err != nil {
			return err
		}

		return nil
	},
}

func extract(file, dir, version string) error {
	// move the folder contents to the home dir
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	extractDir := fmt.Sprintf("./%s", dir)
	homeDir := fmt.Sprintf("%s/.rs", homedir)
	homeVerDir := fmt.Sprintf("%s/.rs/%s", homedir, version)

	// extract the package
	tarCmd := exec.Command("tar", "-zxvf", file)
	err = tarCmd.Run()
	if err != nil {
		return err
	}
	println("Extrated the binaries to", dir)

	err = os.RemoveAll(homeVerDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(homeDir, os.ModePerm)
	if err != nil {
		return err
	}
	println("Init home dir at", homeDir)

	err = os.Rename(extractDir, homeVerDir)
	if err != nil {
		println(err.Error())
		return err
	}
	println("moved the binaries to", homeVerDir)

	// delete the package extracted
	err = os.RemoveAll(extractDir)
	if err != nil {
		return err
	}
	println("cleaned up files at", extractDir)

	return nil
}

func getDownloadUrl(version string, os string, distro string, major int64, minor int64) (string, string, string, error) {
	isBefore42 := major < 4 || (major == 4 && minor < 2)
	var file, dir string

	switch os {
	case "linux":
		if isBefore42 {
			dir = fmt.Sprintf("mongodb-linux-x86_64-%s", version)
		} else {
			dir = fmt.Sprintf("mongodb-linux-x86_64-%s-%s", distro, version)
		}
		break
	case "darwin":
		os = "osx"
		if isBefore42 {
			dir = fmt.Sprintf("mongodb-osx-ssl-x86_64-%s", version)
		} else {
			dir = fmt.Sprintf("mongodb-macos-x86_64-%s", version)
		}
		break
	default:
		return "", "", "", fmt.Errorf("Unrecognized os %s", os)
	}

	file = fmt.Sprintf("%s.tgz", dir)
	url := fmt.Sprintf("%s/%s/%s", base, os, file)

	return url, dir, file, nil
}

func parseVersionNumber(version string) (int64, int64, error) {
	r, err := regexp.Compile(`^(\d)\.(\d)\.(\d+)$`)
	if err != nil {
		return 0, 0, err
	}

	matched := r.FindAll([]byte(version), -1)

	if len(matched) == 0 {
		return 0, 0, ErrInvalidVersionFormat
	}

	split := strings.Split(version, ".")
	major, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	minor, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return major, minor, nil
}

func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
