package main

import (
	"fmt"
	"os"

	"github.com/cbsinteractive/tfc-cli/cmd"
)

func main() {
	options := cmd.ExecuteOpts{
		AppName: "tfc-cli",
	}
	if err := cmd.Execute(options); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
