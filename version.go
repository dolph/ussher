package main

import (
	"fmt"
)

var GitTag string
var GitCommit string
var GoVersion string
var BuildTimestamp string
var BuildOS string
var BuildArch string
var BuildTainted string

func PrintVersion() {
	if GitTag != "" {
		if BuildTainted == "true" {
			fmt.Printf("Version: %s (tainted)\n", GitTag)
		} else {
			fmt.Printf("Version: %s\n", GitTag)
		}
	}
	fmt.Printf("Commit SHA: %s\n", GitCommit)
	fmt.Printf("Go Version: %s\n", GoVersion)
	fmt.Printf("Built: %s\n", BuildTimestamp)
	fmt.Printf("OS/Arch: %s/%s\n", BuildOS, BuildArch)
}
