# internal/cli/testdata

This directory contains test input and expected output files used for verifying
the `uit` CLI tool via golden file tests.

## Structure

```
testdata/
├── input/       # Input files organized by scenario
└── golden/      # Expected output for each corresponding scenario
```

Each subdirectory under `input/` corresponds to a test case. The file with the
same name under `golden/` contains the expected CLI output for that scenario.

## Updating Golden Files

To regenerate golden files based on the current CLI output:

```bash
go test ./internal/cli -update
```

This will overwrite files in the `golden/` directory with new output.

## Regenerating Input Files

To rebuild the test input directory structure from scratch:

```bash
scripts/mk-testdata.sh
```

This script creates consistent input files for all predefined scenarios under
`input/`.
