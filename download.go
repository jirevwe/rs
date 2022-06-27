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
	base                    = "https://fastdl.mongodb.org"
	ErrInvalidVersionFormat = errors.New("version must be in x.x.x format")
	ErrMissingVersionArg    = errors.New("please pass a valid mongodb version")
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads and configures a mongodb version",
	Long:  "Downloads and configures a mongodb version",
	RunE: func(cmd *cobra.Command, args []string) error {
		// default mongodb version
		version := "4.2.21"

		if len(args) > 0 {
			version = args[0]
		}

		major, minor, err := parseVersionNumber(version)
		if err != nil {
			return ErrMissingVersionArg
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
	// extract the package
	tarCmd := exec.Command("tar", "-zxvf", file)
	err := tarCmd.Run()
	if err != nil {
		return err
	}
	println("extrated the binaries to", dir)

	// move the folder contents to the home dir
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	home := fmt.Sprintf("%s/.rs/%s", homedir, version)
	makeCmd := exec.Command("mkdir", "-p", home)
	err = makeCmd.Run()
	if err != nil {
		return err
	}
	println("init home dir at", home)

	mvCmd := exec.Command("mv", "-v", fmt.Sprintf("./%s/bin", dir), home)
	err = mvCmd.Run()
	if err != nil {
		return err
	}
	println("moved the binaries to", home)

	// delete the package extracted
	rmCmd := exec.Command("rm", "-rf", dir)
	err = rmCmd.Run()
	if err != nil {
		return err
	}
	println("cleaned up done")

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

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVar(&distro, "distro", "ubuntu1804", "allows you specify the linux distro")
}
