package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/urfave/cli/v2"

	"github.com/mnishiguchi/uit/internal/fileview"
	"github.com/mnishiguchi/uit/internal/gitutil"
	"github.com/mnishiguchi/uit/internal/treeview"
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
				Usage: "copy output to clipboard",
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

			return Execute(
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

func Execute(
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
	out := io.MultiWriter(writer, &clipboardBuf)

	if !noTree {
		if err := treeview.TreeViewFromGit(inputPath, out); err == nil {
			fmt.Fprintln(out)
			fmt.Fprintln(out)
		}
	}

	if noContent {
		return finalizeOutput(clipboardBuf, copyToClipboard)
	}

	if err := renderFiles(inputPath, maxLines, useFZF, filterPattern, out); err != nil {
		return err
	}

	return finalizeOutput(clipboardBuf, copyToClipboard)
}

func renderFiles(path string, maxLines int, useFZF bool, filter string, out io.Writer) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if !info.IsDir() {
		if err := fileview.FileViewWithLines(path, out, maxLines); err != nil {
			return fmt.Errorf("failed to render file %s: %w", path, err)
		}
		return nil
	}

	files, err := listGitFiles(path)
	if err != nil {
		return err
	}

	files = filterFiles(files, filter)

	if useFZF && isFZFInstalled() {
		files, err = selectFilesWithFZF(files)
		if err != nil {
			return err
		}
	}

	for _, f := range files {
		if err := fileview.FileViewWithLines(f, out, maxLines); err != nil {
			return fmt.Errorf("failed to render file %s: %w", f, err)
		}
	}

	return nil
}

func listGitFiles(path string) ([]string, error) {
	files, err := gitutil.ListGitTrackedFiles(path)
	if err != nil {
		if errors.Is(err, gitutil.ErrNotGitRepo) {
			return nil, fmt.Errorf("this directory is not inside a Git repository: %s", path)
		}

		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no Git-tracked files found in: %s", path)
	}

	return files, nil
}

func filterFiles(files []string, pattern string) []string {
	if pattern == "" {
		return files
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return []string{} // Or log invalid regex warning?
	}

	var filtered []string
	for _, f := range files {
		if re.MatchString(f) {
			filtered = append(filtered, f)
		}
	}

	return filtered
}

func finalizeOutput(buf bytes.Buffer, enabled bool) error {
	if enabled {
		if err := clipboard.WriteAll(buf.String()); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}

		fmt.Fprintln(os.Stderr, "✔️ Copied to clipboard.")
	}

	return nil
}

func isFZFInstalled() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

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
