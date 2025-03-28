package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mnishiguchi/command-line-go/uit/internal/formatter"
	"github.com/urfave/cli/v2"
)

// NewApp returns a CLI app instance for uit.
func NewApp(version string) *cli.App {
	return &cli.App{
		Name:      "uit",
		Usage:     "Render directory tree and file contents from a Git repo",
		UsageText: "uit [options] [path]",
		Version:   version,
		Authors: []*cli.Author{
			{
				Name: "Masatoshi Nishiguchi",
			},
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "max-lines",
				Usage: "limit the number of lines printed per file",
			},
			&cli.BoolFlag{
				Name:  "no-tree",
				Usage: "do not render the tree view",
			},
			&cli.BoolFlag{
				Name:  "no-content",
				Usage: "do not render file contents",
			},
		},
		Action: func(c *cli.Context) error {
			inputPath := "."

			// Use argument as path if provided
			if c.Args().Len() > 0 {
				inputPath = c.Args().First()
			}

			return Run(
				inputPath,
				c.Int("max-lines"),
				c.Bool("no-tree"),
				c.Bool("no-content"),
			)
		},
	}
}

// Run executes the main logic using the given config.
func Run(inputPath string, maxLines int, noTree, noContent bool) error {
	// Print Git-aware tree structure rooted at given path
	if !noTree {
		if err := formatter.RenderGitTree(inputPath, os.Stdout); err == nil {
			// two blank lines after tree if tree was printed
			fmt.Fprintln(os.Stdout)
			fmt.Fprintln(os.Stdout)
		}
	}

	// Skip file content rendering entirely
	if noContent {
		return nil
	}

	// Check if the input path is a file or directory
	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if info.IsDir() {
		// Render all Git-tracked files under the directory
		files, err := formatter.ListGitFilesUnder(inputPath)
		if err != nil {
			if strings.Contains(err.Error(), "not a Git repository") {
				return fmt.Errorf("this directory is not inside a Git repository: %s", inputPath)
			}
			return fmt.Errorf("failed to list files: %w", err)
		}

		for _, f := range files {
			if err := formatter.RenderFileContent(f, os.Stdout, maxLines); err != nil {
				return fmt.Errorf("failed to render file %s: %w", f, err)
			}
		}
	} else {
		// Render a single file
		if err := formatter.RenderFileContent(inputPath, os.Stdout, maxLines); err != nil {
			return fmt.Errorf("failed to render file %s: %w", inputPath, err)
		}
	}

	return nil
}
