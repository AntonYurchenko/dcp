package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/AntonYurchenko/dcp/internal/dev"
)

const (
	dcpEnvName        string = "DCP_RULE"
	dcpEnvRuleSep     string = ";"
	dcpEnvRuleNameSep string = ":"
	dcpEnvArgsSep     string = ">"

	helpMessage string = `this is a tool for copy raw data to device.

    Using:
        - serialization example:     dcp [path to file] [path to device]
        - deserialization example:   dcp [path to device] [path to file] 
	`
)

func main() {

	// Reading and validation arguments.
	first, second, err := readArgs()
	switch {
	case err != nil:
		fmt.Printf("\ndcp: %s\n", err.Error())
		os.Exit(1)
	case strings.HasPrefix(first, "/dev/") && strings.HasPrefix(second, "/dev/"):
		fmt.Print("\ndcp: only one arguments may be device.")
		os.Exit(1)
	}

	if strings.HasPrefix(first, "/dev/") {

		// Reading from device.
		size, err := devToFile(second, first)
		if err != nil {
			fmt.Printf("dcp: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("Downloaded %d bytes to file %q from device %q\n", size, second, first)

	} else {

		// Writing to device.
		size, err := fileToDev(second, first)
		if err != nil {
			fmt.Printf("dcp: %v\n", err)
			os.Exit(3)
		}
		fmt.Printf("Uploaded %d bytes from file %q to device %q\n", size, first, second)

	}
}

// parseEnv extract custom rule from environment.
func parseEnv(ruleName, env string) (first, second string, err error) {

	if ruleName == "" || env == "" {
		return "", "", errors.New("invalid arguments")
	}

	var arr []string
	for _, rule := range strings.Split(env, dcpEnvRuleSep) {

		arr = strings.SplitN(rule, dcpEnvRuleNameSep, 2)
		if len(arr) == 2 && arr[0] == ruleName {

			arr = strings.SplitN(arr[1], dcpEnvArgsSep, 2)
			if len(arr) == 2 && arr[0] != "" && arr[1] != "" {
				return arr[0], arr[1], nil
			}

		}

	}

	return "", "", errors.New("arguments not found")
}

func readArgs() (first, second string, err error) {

	flag.Parse()
	args := flag.Args()
	switch {
	case len(args) == 1 && args[0] == "help":
		return "", "", errors.New(helpMessage)
	case len(args) == 1:
		return parseEnv(args[0], os.Getenv(dcpEnvName))
	case len(args) == 2 && args[0] != "" && args[1] != "":
		return args[0], args[1], nil
	}

	return "", "", errors.New("arguments not found")
}

func devToFile(filePath, devPath string) (size int64, err error) {

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, fmt.Errorf("open file %q error: %v", filePath, err)
	}
	defer file.Close()

	size, err = dev.CopyFrom(file, devPath)
	if err != nil {
		return 0, fmt.Errorf("read file %q from device %q error: %v", filePath, devPath, err)
	}
	return size, nil
}

func fileToDev(devPath, filePath string) (size int64, err error) {

	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return 0, fmt.Errorf("open file %q error: %v", filePath, err)
	}
	defer file.Close()

	size, err = dev.CopyTo(devPath, file)
	if err != nil {
		return 0, fmt.Errorf("write data from file %q to device %q error: %v", filePath, devPath, err)
	}
	return size, nil
}
