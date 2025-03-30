package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
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
				Value: 500,
			},
			&cli.BoolFlag{
				Name:  "no-tree",
				Usage: "do not render the tree view",
			},
			&cli.BoolFlag{
				Name:  "no-content",
				Usage: "do not render file contents",
			},
			&cli.BoolFlag{
				Name:  "copy",
				Usage: "copy output to clipboard instead of printing",
			},
			&cli.BoolFlag{
				Name:  "fzf",
				Usage: "interactively select files via fzf (if installed)",
			},
			&cli.StringFlag{
				Name:  "filter",
				Usage: "filter file paths with a regular expression",
			},
		},
		Action: func(c *cli.Context) error {
			inputPath := "."
			if c.Args().Len() > 0 {
				inputPath = c.Args().First()
			}

			return Run(
				inputPath,
				c.Int("max-lines"),
				c.Bool("no-tree"),
				c.Bool("no-content"),
				c.Bool("copy"),
				c.Bool("fzf"),
				c.String("filter"),
				c.App.Writer,
			)
		},
	}
}

// Run executes the main logic using the given config.
func Run(
	inputPath string,
	maxLines int,
	noTree bool,
	noContent bool,
	copyToClipboard bool,
	useFZF bool,
	filterPattern string,
	writer io.Writer,
) error {
	var clipboardBuf bytes.Buffer

	// Multi-writer sends output to both terminal/test and clipboard buffer
	out := io.MultiWriter(writer, &clipboardBuf)

	// Print Git-aware tree structure rooted at given path
	if !noTree {
		if err := formatter.RenderGitTree(inputPath, out); err == nil {
			// Two blank lines after tree if tree was printed
			fmt.Fprintln(out)
			fmt.Fprintln(out)
		}
	}

	// Skip file content rendering entirely
	if noContent {
		return outputResult(clipboardBuf, copyToClipboard)
	}

	// Check if the input path is a file or directory
	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if info.IsDir() {
		// List Git-tracked files under the directory
		files, err := formatter.ListGitFilesUnder(inputPath)
		if err != nil {
			if strings.Contains(err.Error(), "not a Git repository") {
				return fmt.Errorf("this directory is not inside a Git repository: %s", inputPath)
			}
			return fmt.Errorf("failed to list files: %w", err)
		}

		if len(files) == 0 {
			return fmt.Errorf("no Git-tracked files found in: %s", inputPath)
		}

		// If --filter is provided, filter file paths
		if filterPattern != "" {
			re, err := regexp.Compile(filterPattern)
			if err != nil {
				return fmt.Errorf("invalid regex pattern: %w", err)
			}

			var filtered []string
			for _, f := range files {
				if re.MatchString(f) {
					filtered = append(filtered, f)
				}
			}

			files = filtered
		}

		// If --fzf is enabled and fzf is available, let the user interactively select files
		if useFZF && isFZFInstalled() {
			files, err = selectFilesWithFZF(files)
			if err != nil {
				return err
			}
		}

		for _, f := range files {
			if err := formatter.RenderFileContent(f, out, maxLines); err != nil {
				return fmt.Errorf("failed to render file %s: %w", f, err)
			}
		}
	} else {
		// Render a single file
		if err := formatter.RenderFileContent(inputPath, out, maxLines); err != nil {
			return fmt.Errorf("failed to render file %s: %w", inputPath, err)
		}
	}

	return outputResult(clipboardBuf, copyToClipboard)
}

// outputResult handles writing the final output, either to stdout or clipboard.
func outputResult(buf bytes.Buffer, copyToClipboard bool) error {
	if copyToClipboard {
		if err := clipboard.WriteAll(buf.String()); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		fmt.Fprintln(os.Stderr, "✔️ Copied to clipboard.")
	}
	return nil
}

// isFZFInstalled checks if fzf is available in PATH.
func isFZFInstalled() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

// selectFilesWithFZF launches fzf with the given list of files and returns selected ones.
func selectFilesWithFZF(files []string) ([]string, error) {
	cmd := exec.Command("fzf", "--multi")
	cmd.Stdin = strings.NewReader(strings.Join(files, "\n"))

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("fzf failed: %w", err)
	}

	selection := strings.Split(strings.TrimSpace(string(out)), "\n")
	return selection, nil
}
