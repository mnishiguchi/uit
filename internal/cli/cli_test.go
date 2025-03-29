package cli_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/mnishiguchi/command-line-go/uit/internal/cli"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("copies rendered output to clipboard", func(t *testing.T) {
		// Skip if clipboard isn't supported (e.g. in CI)
		if err := clipboard.WriteAll("test"); err != nil {
			t.Skip("Skipping clipboard test: clipboard not available")
		}

		// Create a temporary text file with some content
		tmpFile, err := os.CreateTemp("", "uit-clipboard-test-*.txt")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		content := "This is a test for clipboard output."
		_, err = tmpFile.WriteString(content)
		assert.NoError(t, err)
		tmpFile.Close()

		var buf bytes.Buffer

		err = cli.Run(
			tmpFile.Name(), // inputPath
			0,              // maxLines
			true,           // noTree
			false,          // noContent
			true,           // copyToClipboard
			false,          // useFzf
			"",             // filterRegex
			&buf,
		)
		assert.NoError(t, err)

		clip, err := clipboard.ReadAll()
		assert.NoError(t, err)
		assert.Contains(t, clip, content)
		assert.Contains(t, buf.String(), content)
	})
}
