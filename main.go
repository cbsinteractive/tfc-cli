package main

import (
	"os"

	"github.com/cbsinteractive/tfc-cli/cmd"
)

func main() {
	options := cmd.ExecuteOpts{
		AppName: "tfc-cli",
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
	if err := cmd.Execute(options); err != nil {
		// If an error took place, code within the testable part of the
		// application was expected to present output via os.Stderr
		os.Exit(1)
	}
}
