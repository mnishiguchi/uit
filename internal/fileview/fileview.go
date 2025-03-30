package fileview

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mnishiguchi/command-line-go/uit/internal/gitutil"
)

// FileViewWithLines prints the content of a single file to the writer with line numbers.
func FileViewWithLines(path string, w io.Writer, maxLines int) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("cannot render directory as file: %s", absPath)
	}

	if isBin, err := isBinaryFile(absPath); err == nil && isBin {
		printFileHeader(w, path)
		fmt.Fprintln(w, "[binary file omitted]")
		printFileFooter(w)
		return nil
	}

	file, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	printFileHeader(w, path)

	if err := printFileBodyWithLines(file, w, maxLines); err != nil {
		return err
	}

	printFileFooter(w)

	return nil
}

func printFileHeader(w io.Writer, path string) {
	// Try to get path relative to Git root if available
	gitRoot, err := gitutil.GetGitRoot(path)
	if err != nil {
		relPath, relErr := filepath.Rel(".", path)
		if relErr != nil {
			relPath = path
		}

		fmt.Fprintf(w, "/%s:\n", filepath.ToSlash(relPath))
	} else {
		relToGitRoot, err := filepath.Rel(gitRoot, path)
		if err != nil {
			relToGitRoot = path
		}

		fmt.Fprintf(w, "\n\n/%s:\n", filepath.ToSlash(relToGitRoot))
	}

	fmt.Fprintln(w, strings.Repeat("-", 80))
}

func printFileBodyWithLines(file *os.File, w io.Writer, maxLines int) error {
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		if maxLines > 0 && lineNum > maxLines {
			break
		}
		fmt.Fprintf(w, "%4d | %s\n", lineNum, scanner.Text())
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func printFileFooter(w io.Writer) {
	fmt.Fprintln(w, "\n\n"+strings.Repeat("-", 80))
}

// isBinaryFile returns true if the file contains a null byte in the first 8000 bytes.
func isBinaryFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	const maxBytes = 8000
	buf := make([]byte, maxBytes)

	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	return bytes.IndexByte(buf[:n], 0) >= 0, nil
}
