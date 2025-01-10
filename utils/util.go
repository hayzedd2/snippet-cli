package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type saveOptions struct {
	Tag       string
	FilePath  string
	StartLine int64
	EndLine   int64
}

type saveAllOptions struct {
	Tag      string
	FilePath string
}

func ParseSaveOptions(cmd *cobra.Command) (*saveOptions, error) {
	opts := &saveOptions{}
	var err error
	opts.Tag, _ = cmd.Flags().GetString("tag")
	if opts.Tag == "" {
		return nil, returnError("tag", "some")
	}
	opts.FilePath, _ = cmd.Flags().GetString("filepath")
	if opts.FilePath == "" {
		return nil, returnError("filepath", "some")
	}
	startLine, _ := cmd.Flags().GetString("startline")
	if startLine == "" {
		return nil, returnError("startline", "some")
	}
	opts.StartLine, err = strconv.ParseInt(startLine, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("start line must be a number")
	}
	endLine, err := cmd.Flags().GetString("endline")
	if err != nil {
		return nil, fmt.Errorf("invalid end line: %w", err)
	}

	// If endLine is not provided, use startLine
	if endLine == "" {
		opts.EndLine = opts.StartLine
	} else {
		opts.EndLine, err = strconv.ParseInt(endLine, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("end line must be a number")
		}
	}

	if opts.EndLine < opts.StartLine {
		return nil, fmt.Errorf("end line cannot be less than start line")
	}

	return opts, nil
}

func ParseSaveAllOptions(cmd *cobra.Command) (*saveAllOptions, error) {
	opts := &saveAllOptions{}
	opts.FilePath, _ = cmd.Flags().GetString("filepath")
	if opts.FilePath == "" {
		return nil, returnError("filepath", "all")
	}
	opts.Tag, _ = cmd.Flags().GetString("tag")
	if opts.Tag == "" {
		return nil, returnError("tag", "all")
	}
	return opts, nil
}
func GetStorageFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	snippetDir := filepath.Join(homeDir, ".snippets")
	if err := os.MkdirAll(snippetDir, 0755); err != nil {
		fmt.Println("Error creating snippets directory:", err)
		os.Exit(1)
	}
	return filepath.Join(snippetDir, "snippets.json")
}

func GetCodeFromFile(filePath string, startLine, endLine int64) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	if startLine < 1 {
		return "", fmt.Errorf("start line must be greater than 0")
	}
	if endLine < startLine {
		return "", fmt.Errorf("end line must be greater than or equal to start line")
	}
	var selectedLines []string
	scanner := bufio.NewScanner(file)
	var currentLine int64 = 1
	for scanner.Scan() {
		if currentLine >= startLine && currentLine <= endLine {
			selectedLines = append(selectedLines, scanner.Text())
		}
		if currentLine > endLine {
			break
		}
		currentLine++
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file '%s': %v", filePath, err)
	}
	// Handle case where file has fewer lines than requested
	if currentLine <= startLine {
		return "", fmt.Errorf("file has only %d lines, requested start line was %d", currentLine-1, startLine)
	}
	code := strings.Join(selectedLines, "\n")
	return code, nil
}

func returnError(s, t string) error {
	if t == "some" {
		return fmt.Errorf("%v is required \nUse `snippet save -t tag -f filepath -startline startline -endline? endline` to save a snippet", s)
	} else {
		return fmt.Errorf("%v is required \nUse `snippet save all -f filepath -t tag` to save all content in a file", s)
	}
}

func GetAllContentFromFile(filePath string) (string, error) {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return "", fmt.Errorf("error reading file: %w", err)
    }
    return string(content), nil
}

