package gitutil

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// FindGitRoot returns the absolute path of the Git repository root for the given path.
func FindGitRoot(path string) (string, error) {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a Git repository: %s", path)
	}

	return strings.TrimSpace(string(output)), nil
}

// ListGitFilesUnder returns a list of Git-tracked files under a given directory.
func ListGitFilesUnder(dir string) ([]string, error) {
	cmd := exec.Command("git", "-C", dir, "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	files := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue // Skip blank lines to avoid treating directory as file
		}
		files = append(files, filepath.Join(dir, line))
	}
	return files, nil
}
