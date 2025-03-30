package fileview_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mnishiguchi/command-line-go/uit/internal/fileview"
	"github.com/mnishiguchi/command-line-go/uit/internal/gitutil"
	"github.com/mnishiguchi/command-line-go/uit/internal/treeview"
	"github.com/stretchr/testify/assert"
)

func TestRenderGitTree(t *testing.T) {
	t.Run("prints tree including known file", func(t *testing.T) {
		var buf bytes.Buffer

		cwd, err := os.Getwd()
		assert.NoError(t, err)

		err = treeview.RenderGitTree(cwd, &buf)
		assert.NoError(t, err)

		output := buf.String()

		t.Run("includes root directory name", func(t *testing.T) {
			expectedRoot := filepath.Base(cwd)
			assert.Contains(t, output, expectedRoot)
		})

		t.Run("includes known file", func(t *testing.T) {
			assert.Contains(t, output, "fileview.go", "expected tree output to contain a known file")
		})
	})
}

func TestRenderFileContent(t *testing.T) {
	t.Run("renders file with line limit", func(t *testing.T) {
		tmpDir := t.TempDir()
		textFile := filepath.Join(tmpDir, "sample-head.txt")
		content := `line 1
line 2
line 3
line 4
line 5`
		err := os.WriteFile(textFile, []byte(content), 0644)
		assert.NoError(t, err)

		var buf bytes.Buffer
		err = fileview.RenderFileContent(textFile, &buf, 3)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "   1 | line 1")
		assert.Contains(t, output, "   2 | line 2")
		assert.Contains(t, output, "   3 | line 3")
		assert.NotContains(t, output, "line 4")
		assert.NotContains(t, output, "line 5")
	})

	t.Run("renders full file when no limit", func(t *testing.T) {
		tmpDir := t.TempDir()
		textFile := filepath.Join(tmpDir, "sample-full.txt")
		content := `line A
line B
line C`
		err := os.WriteFile(textFile, []byte(content), 0644)
		assert.NoError(t, err)

		var buf bytes.Buffer
		err = fileview.RenderFileContent(textFile, &buf, 0)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "   1 | line A")
		assert.Contains(t, output, "   2 | line B")
		assert.Contains(t, output, "   3 | line C")
	})
}

func TestFindGitRoot(t *testing.T) {
	t.Run("returns error for non-Git directory", func(t *testing.T) {
		tmp := t.TempDir()
		_, err := gitutil.FindGitRoot(tmp)
		assert.Error(t, err)
	})
}

func TestListGitFilesUnder(t *testing.T) {
	t.Run("returns error for non-Git directory", func(t *testing.T) {
		tmp := t.TempDir()
		_, err := gitutil.ListGitFilesUnder(tmp)
		assert.Error(t, err)
	})
}

func TestRenderFileContent_RejectsDirectory(t *testing.T) {
	t.Run("returns error when given a directory", func(t *testing.T) {
		dir := t.TempDir()
		var buf bytes.Buffer

		err := fileview.RenderFileContent(dir, &buf, 0)
		assert.ErrorContains(t, err, "cannot render directory as file")
	})
}
