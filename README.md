# uit

Render directory tree and file contents from a Git repository.

## Features

- Displays Git-tracked files as a tree
- Shows file content with line numbers
- Supports partial output via `--head`
- Skips binary files by default (override with `--show-binary`)
- Toggle visibility of tree and contents

## Installation

Download a binary from [Releases](https://github.com/mnishiguchi/uit/releases)  
or build from source:

```sh
make release
```

## Usage

```sh
uit [options] [path]
```

### Options

| Option          | Description                           |
|-----------------|---------------------------------------|
| `--head N`      | Limit the number of lines per file    |
| `--show-binary` | Show contents of binary files         |
| `--no-tree`     | Suppress tree view                    |
| `--no-content`  | Suppress file content view            |

## Example

```sh
uit --head 5 ./internal
```

## License

MIT
