package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/zpnst/containy/linux"
)

const (
	helpText = `Usage: 
	containy [OPTIONS] COMMAND
Name: 
	containy - minimal Docker-like container runtime

Description:
	Containy is a minimal Docker-like container runtime for 
	academic purposes. It creates containers from bundles, similar to runC

Options:
	-b, --bundle   option for specifying the path to the bundle, 
		without this option configy.json will be taken from the working directory.

Common Commands:
	run         Create and run a new container from, using configy.json file
	version     Show the containy version information
	help        Shows you this help

For more help on how to use containy, head to https://github.com/zpnst/containy`
)

func main() {
	if len(os.Args) < 2 {
		showHelpExit()
	}

	command := os.Args[1]

	switch command {
	case "run":
		var pathToTheBundle string
		var pathToTheBundleWithConfig string
		var defaultConfigPath string = "configy.json"
		if len(os.Args) == 4 {
			if (os.Args[2] == "-b") || (os.Args[2] == "--bundle") {
				pathToTheBundle = os.Args[3]
				pathToTheBundleWithConfig = fmt.Sprintf("%s/%s", os.Args[3], defaultConfigPath)
			}
		} else {
			pathToTheBundle = ""
			pathToTheBundleWithConfig = defaultConfigPath
		}

		if _, err := os.Stat(pathToTheBundleWithConfig); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Invalid path to the bundle: %s\n", pathToTheBundleWithConfig)
			os.Exit(1)
		}

		configy, err := linux.ParseConfigy(pathToTheBundleWithConfig)
		if err != nil {
			log.Panicln(err)
		}
		contaiery := linux.NewContaiery(*configy, pathToTheBundle)

		if os.Getenv("IN_CONTAINER") == "TRUE" {
			contaiery.ContainerRuntime()
			return
		}

		fmt.Println("Running container from bundle")
		contaiery.CreateContainer()
		return

	case "help":
		showHelp()

	case "version":
		showVersion()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		showHelpExit()
	}

}

func showVersion() {
	fmt.Println("containy version 0.1.0")
}

func showHelp() {
	fmt.Println(helpText)
}

func showHelpExit() {
	fmt.Println(helpText)
	os.Exit(1)
}
