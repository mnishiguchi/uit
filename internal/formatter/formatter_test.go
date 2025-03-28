package formatter_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mnishiguchi/command-line-go/uit/internal/formatter"
	"github.com/stretchr/testify/assert"
)

func TestRenderGitTree(t *testing.T) {
	var buf bytes.Buffer

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	err = formatter.RenderGitTree(cwd, &buf)
	assert.NoError(t, err)

	output := buf.String()

	t.Run("prints root directory name", func(t *testing.T) {
		expectedRoot := filepath.Base(cwd)
		assert.Contains(t, output, expectedRoot)
	})

	t.Run("includes known file", func(t *testing.T) {
		assert.Contains(t, output, "formatter.go")
	})
}

func TestRenderFileContent(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("renders file with line limit", func(t *testing.T) {
		textFile := filepath.Join(tmpDir, "sample-head.txt")
		content := `line 1
line 2
line 3
line 4
line 5`
		err := os.WriteFile(textFile, []byte(content), 0644)
		assert.NoError(t, err)

		var buf bytes.Buffer
		err = formatter.RenderFileContent(textFile, &buf, 3)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "   1 | line 1")
		assert.Contains(t, output, "   2 | line 2")
		assert.Contains(t, output, "   3 | line 3")
		assert.NotContains(t, output, "line 4")
		assert.NotContains(t, output, "line 5")
	})

	t.Run("renders full file when no limit", func(t *testing.T) {
		textFile := filepath.Join(tmpDir, "sample-full.txt")
		content := `line A
line B
line C`
		err := os.WriteFile(textFile, []byte(content), 0644)
		assert.NoError(t, err)

		var buf bytes.Buffer
		err = formatter.RenderFileContent(textFile, &buf, 0)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "   1 | line A")
		assert.Contains(t, output, "   2 | line B")
		assert.Contains(t, output, "   3 | line C")
	})
}
