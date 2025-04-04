package main

import (
	"fmt"
	"os"

	"github.com/mnishiguchi/uit/internal/cli"
)

var version = "v-dev"

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "--help")
	}

	app := cli.NewApp(version)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
