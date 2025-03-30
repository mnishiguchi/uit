package cli_test

import (
	"bytes"
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/mnishiguchi/command-line-go/uit/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestRun(t *testing.T) {
	cases := map[string]struct {
		maxLines   int
		noTree     bool
		noContent  bool
		copyToClip bool
		useFZF     bool
		filter     string
	}{
		"default": {
			maxLines: 500,
		},
		"max-lines": {
			maxLines: 3,
		},
		"no-tree": {
			noTree: true,
		},
		"no-content": {
			noContent: true,
		},
		"filter": {
			filter: "a\\.txt$",
		},
		"copy": {
			copyToClip: true,
		},
		"binary": {
			maxLines: 500,
		},
	}

	for label, tt := range cases {
		// Skip if clipboard isn't supported
		if err := clipboard.WriteAll("test"); err != nil {
			t.Skip("Skipping clipboard test: clipboard not available")
		}

		t.Run(label, func(t *testing.T) {
			inputDir := filepath.Join("testdata", "input", label)
			goldenFile := filepath.Join("testdata", "golden", label)

			var buf bytes.Buffer
			err := cli.Run(
				inputDir,
				tt.maxLines,
				tt.noTree,
				tt.noContent,
				tt.copyToClip,
				tt.useFZF,
				tt.filter,
				&buf,
			)
			require.NoError(t, err)

			actual := buf.String()

			if *updateGolden {
				err := os.WriteFile(goldenFile, []byte(actual), 0644)
				require.NoError(t, err)
			}

			expected, err := os.ReadFile(goldenFile)
			require.NoError(t, err)

			assert.Equal(t, string(expected), actual)
		})
	}
}

func TestCopyConfirmationMessage(t *testing.T) {
	// Skip if clipboard isn't supported
	if err := clipboard.WriteAll("test"); err != nil {
		t.Skip("Skipping clipboard test: clipboard not available")
	}

	inputDir := filepath.Join("testdata", "input", "copy")
	var stdoutBuf bytes.Buffer

	// Capture stderr
	stderrReader, stderrWriter, err := os.Pipe()
	require.NoError(t, err)
	originalStderr := os.Stderr
	os.Stderr = stderrWriter
	defer func() {
		os.Stderr = originalStderr
		stderrWriter.Close()
	}()

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, stderrReader)
		done <- buf.String()
	}()

	err = cli.Run(
		inputDir,
		500,   // maxLines
		false, // noTree
		false, // noContent
		true,  // copyToClip
		false, // useFZF
		"",    // filter
		&stdoutBuf,
	)
	require.NoError(t, err)

	stderrWriter.Close()
	stderrOutput := <-done

	assert.Contains(t, stderrOutput, "Copied to clipboard")
}
