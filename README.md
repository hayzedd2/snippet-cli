# Snippet Manager

A CLI tool for developers to save and retrieve commonly used code snippets.

## Installation

```bash
go install github.com/hayzedd2/snippet-cli/cmd/snippet@latest
```

## Features

- Save code snippets from files via CLI
- Save entire files as snippets
- Copy snippets to the clipboard for quick use
- List all saved snippets with their tags and creation timestamps
- Retrieve specific snippets by their tags
- Delete snippets when they're no longer needed

## Usage

### Save a Snippet

The save command has two modes: saving specific lines or saving an entire file.

#### Save Specific Lines
The standard save command allows you to store a snippet from specific lines in a file. You must specify the starting line number. For multi-line snippets, you can optionally specify an end line.

**Full Command:**
```bash
# For multiple lines
snippet save --tag <tag> --filepath <filepath> --startline <startline> --endline <endline>

# For a single line
snippet save --tag <tag> --filepath <filepath> --startline <startline>
```

**Short Form:**
```bash
# For multiple lines
snippet save -t <tag> -f <filepath> -s <startline> -e <endline>

# For a single line
snippet save -t <tag> -f <filepath> -s <startline>
```

**Parameters:**
- `tag`: A unique identifier for your snippet
- `filepath`: Path to the source file
- `startline`: Starting line number of the snippet (required)
- `endline`: Ending line number of the snippet (optional, for multi-line snippets)

#### Save Entire File
You can save the entire contents of a file as a snippet using the `save all` command. This doesn't require specifying line numbers.

**Full Command:**
```bash
snippet save all --file <filepath> --tag <tag>
```

**Short Form:**
```bash
snippet save all -f <filepath> -t <tag>
```

### List All Snippets

```bash
snippet list
```

### Get a Specific Snippet

```bash
snippet get --tag <tag>
```

### Copy Snippet to Clipboard

```bash
snippet copy --tag <tag>
```

### Delete a Snippet

```bash
snippet delete --tag <tag>
```

## Examples

Save a single line:
```bash
snippet save --tag "log-function" --filepath cmd/snippet/main.go --startline 15
```

Save multiple lines:
```bash
snippet save --tag "parse-options" --filepath cmd/snippet/main.go --startline 10 --endline 40
```

Save an entire file:
```bash
snippet save all -t "main-function" -f main.go
```

Using short form:
```bash
# Single line
snippet save -t "error-handler" -f errors.go -s 25

# Multiple lines
snippet save -t "middleware" -f server.go -s 30 -e 45
```

Get a saved snippet:
```bash
snippet get --tag "log-function"
```

## Tips

- Line numbers start at 1 (not 0)
- The `startline` parameter is always required when not using `save all`
- For single-line snippets, omit the `endline` parameter
- When `endline` is specified, both start and end lines are included in the snippet
- Use `save all` when you want to store the entire contents of a file