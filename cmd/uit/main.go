package main

import (
	"fmt"
	"os"

	"github.com/mnishiguchi/command-line-go/uit/internal/cli"
)

var version = "v-dev"

func main() {
	app := cli.NewApp(version)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
