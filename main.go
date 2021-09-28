package main

import (
	"os"

	"github.com/cbsinteractive/tfc-cli/cmd"
)

func main() {
	options := cmd.ExecuteOpts{
		AppName: "tfc-cli",
		Writer:  os.Stdout,
	}
	if err := cmd.Execute(options); err != nil {
		// Command handlers are responsible for output
		os.Exit(1)
	}
}
