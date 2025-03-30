package treeview

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mnishiguchi/command-line-go/uit/internal/gitutil"
)

// TreeNode represents a node in the directory tree.
type TreeNode struct {
	Name     string
	IsFile   bool
	Children map[string]*TreeNode
}

// TreeViewFromGit builds and prints a Git-tracked file tree starting from the user-specified path.
func TreeViewFromGit(inputPath string, w io.Writer) error {
	tree, err := buildTreeFromGit(inputPath)
	if err != nil {
		return err
	}

	printTreeRoot(tree, w)

	return nil
}

func buildTreeFromGit(inputPath string) (*TreeNode, error) {
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve input path: %w", err)
	}

	gitRoot, err := gitutil.GetGitRoot(absInput)
	if err != nil {
		return nil, fmt.Errorf("failed to find git root: %w", err)
	}

	cmd := exec.Command("git", "-C", gitRoot, "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run git ls-files: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	relInputPath, err := filepath.Rel(gitRoot, absInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative input path: %w", err)
	}

	var relevantPaths [][]string
	for _, line := range lines {
		if relInputPath == "." || strings.HasPrefix(line, relInputPath+"/") || line == relInputPath {
			trimmed := strings.TrimPrefix(line, relInputPath+"/")
			relevantPaths = append(relevantPaths, strings.Split(trimmed, "/"))
		}
	}

	root := &TreeNode{
		Name:     filepath.Base(absInput),
		IsFile:   false,
		Children: make(map[string]*TreeNode),
	}

	for _, parts := range relevantPaths {
		insertPathParts(root, parts)
	}

	return root, nil
}

// insertPathParts inserts a file path (split into parts) into the tree recursively.
func insertPathParts(node *TreeNode, parts []string) {
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

	insertPathParts(child, parts[1:])
}

// printTreeRoot prints the tree starting from the root node.
func printTreeRoot(node *TreeNode, w io.Writer) {
	fmt.Fprintf(w, "%s\n", node.Name)
	printTreeChildren(node, "", true, w)
}

// printTreeChildren prints child nodes of the given tree node recursively.
func printTreeChildren(node *TreeNode, prefix string, isLast bool, w io.Writer) {
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
			printTreeChildren(child, nextPrefix, i == len(keys)-1, w)
		}
	}
}
