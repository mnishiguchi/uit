package cli_test

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/mnishiguchi/command-line-go/uit/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestRun(t *testing.T) {
	t.Run("matches golden output for single file", func(t *testing.T) {
		inputPath := filepath.Join("testdata", "input", "test-this.txt")
		goldenPath := filepath.Join("testdata", "golden", "test-this.txt")

		// Configurable CLI flags
		maxLines := 500
		noTree := true
		noContent := false
		copyToClipboard := false
		useFZF := false
		filter := ""

		var buf bytes.Buffer

		err := cli.Run(
			inputPath,
			maxLines,
			noTree,
			noContent,
			copyToClipboard,
			useFZF,
			filter,
			&buf,
		)
		require.NoError(t, err)

		actual := buf.String()

		if *updateGolden {
			err := os.WriteFile(goldenPath, []byte(actual), 0644)
			require.NoError(t, err)
		}

		expected, err := os.ReadFile(goldenPath)
		require.NoError(t, err)

		assert.Equal(t, string(expected), actual)
	})
}
