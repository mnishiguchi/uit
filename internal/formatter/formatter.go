package formatter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// TreeNode represents a node in the directory tree.
type TreeNode struct {
	Name     string
	IsFile   bool
	Children map[string]*TreeNode
}

// RenderGitTree builds and prints a Git-tracked file tree starting from the user-specified path.
func RenderGitTree(inputPath string, w io.Writer) error {
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve input path: %w", err)
	}

	gitRoot, err := FindGitRoot(absInput)
	if err != nil {
		return fmt.Errorf("failed to find git root: %w", err)
	}

	cmd := exec.Command("git", "-C", gitRoot, "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run git ls-files: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	relInputPath, err := filepath.Rel(gitRoot, absInput)
	if err != nil {
		return fmt.Errorf("failed to get relative input path: %w", err)
	}

	var relevantPaths [][]string
	for _, line := range lines {
		if relInputPath == "." || strings.HasPrefix(line, relInputPath+"/") || line == relInputPath {
			trimmed := strings.TrimPrefix(line, relInputPath+"/")
			relevantPaths = append(relevantPaths, strings.Split(trimmed, "/"))
		}
	}

	tree := &TreeNode{
		Name:     filepath.Base(absInput),
		IsFile:   false,
		Children: make(map[string]*TreeNode),
	}

	for _, parts := range relevantPaths {
		addPath(tree, parts)
	}

	printTree(tree, w)
	return nil
}

// addPath inserts a file path (split into parts) into the tree recursively.
func addPath(node *TreeNode, parts []string) {
	if len(parts) == 0 {
		return
	}

	name := parts[0]
	child, exists := node.Children[name]
	if !exists {
		child = &TreeNode{
			Name:     name,
			IsFile:   len(parts) == 1,
			Children: make(map[string]*TreeNode),
		}
		node.Children[name] = child
	}

	addPath(child, parts[1:])
}

// printTree prints the tree starting from the root node.
func printTree(node *TreeNode, w io.Writer) {
	fmt.Fprintf(w, "%s\n", node.Name)
	printChildren(node, "", true, w)
}

// printChildren prints child nodes of the given tree node recursively.
func printChildren(node *TreeNode, prefix string, isLast bool, w io.Writer) {
	_ = isLast // Reserved for future enhancements

	var keys []string
	for k := range node.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, key := range keys {
		child := node.Children[key]

		connector := "├──"
		nextPrefix := prefix + "│   "
		if i == len(keys)-1 {
			connector = "└──"
			nextPrefix = prefix + "    "
		}

		fmt.Fprintf(w, "%s%s %s\n", prefix, connector, child.Name)

		if !child.IsFile {
			printChildren(child, nextPrefix, i == len(keys)-1, w)
		}
	}
}

// RenderFileContent prints the content of a single file to the writer with line numbers.
func RenderFileContent(path string, w io.Writer, headLines int) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	if isBin, err := isBinary(absPath); err == nil && isBin {
		printFileContentHeader(w, path)
		fmt.Fprintln(w, "[binary file omitted]")
		printFileContentFooter(w)
		return nil
	}

	file, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	printFileContentHeader(w, path)

	if err := printFileContentBody(file, w, headLines); err != nil {
		return err
	}

	printFileContentFooter(w)

	return nil
}

func printFileContentHeader(w io.Writer, path string) {
	// Try to get path relative to Git root if available
	gitRoot, err := FindGitRoot(path)
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

func printFileContentBody(file *os.File, w io.Writer, headLines int) error {
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		if headLines > 0 && lineNum > headLines {
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

func printFileContentFooter(w io.Writer) {
	fmt.Fprintln(w, "\n\n"+strings.Repeat("-", 80))
}

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
	var files []string
	for _, line := range lines {
		files = append(files, filepath.Join(dir, line))
	}
	return files, nil
}

// isBinary checks if the file content is likely binary.
func isBinary(path string) (bool, error) {
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

	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true, nil
		}
	}

	return false, nil
}
