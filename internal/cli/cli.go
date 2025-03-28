package cli

import (
	"fmt"
	"os"

	"github.com/mnishiguchi/command-line-go/uit/internal/formatter"
	"github.com/urfave/cli/v2"
)

// Config holds CLI options.
type Config struct {
	Path       string
	ShowBinary bool
	HeadLines  int
	NoTree     bool
	NoContent  bool
}

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
			&cli.BoolFlag{
				Name:  "show-binary",
				Usage: "show binary file contents",
				Value: false,
			},
			&cli.IntFlag{
				Name:  "head",
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
			// Parse CLI flags into config.
			cfg := Config{
				Path:       ".",
				ShowBinary: c.Bool("show-binary"),
				NoTree:     c.Bool("no-tree"),
				NoContent:  c.Bool("no-content"),
				HeadLines:  c.Int("head"),
			}

			// Use argument as path if provided
			if c.Args().Len() > 0 {
				cfg.Path = c.Args().First()
			}

			return Run(cfg)
		},
	}
}

// Run executes the main logic using the given config.
func Run(cfg Config) error {
	// Print Git-aware tree structure rooted at given path
	if !cfg.NoTree {
		if err := formatter.RenderGitTree(cfg.Path, os.Stdout); err == nil {
			fmt.Println() // spacer if tree was printed
		}
	}

	// Skip file content rendering entirely
	if cfg.NoContent {
		return nil
	}

	// Check if the input path is a file or directory
	info, err := os.Stat(cfg.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if info.IsDir() {
		// Render all Git-tracked files under the directory
		files, err := formatter.ListGitFilesUnder(cfg.Path)
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		for _, f := range files {
			if err := formatter.RenderFileContent(f, os.Stdout, cfg.ShowBinary, cfg.HeadLines); err != nil {
				return fmt.Errorf("failed to render file %s: %w", f, err)
			}
		}
	} else {
		// Render a single file
		if err := formatter.RenderFileContent(cfg.Path, os.Stdout, cfg.ShowBinary, cfg.HeadLines); err != nil {
			return fmt.Errorf("failed to render file %s: %w", cfg.Path, err)
		}
	}

	return nil
}
