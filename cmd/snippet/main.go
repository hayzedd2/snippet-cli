package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/atotto/clipboard"
)

type Snippet struct {
	Tag       string    `json:"tag"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

type SnippetStore struct {
	Snippets []Snippet `json:"snippets"`
}

var rootCmd = &cobra.Command{
	Use:   "snippet",
	Short: "Snippet manager is a CLI tool for saving code snippets",
	Long:  `A CLI tool that helps developers save and retrieve commonly used code snippets.`,
}

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy an existing snippet",
	Run: func(cmd *cobra.Command, args []string) {
		tag, _ := cmd.Flags().GetString("tag")
		if tag == "" {
			fmt.Println("Tag is required")
			return
		}
		s, err := copySnippet(tag)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = clipboard.WriteAll(s.Code)
		if err != nil {
			fmt.Println("Error copying to clipboard:", err)
			return
		}
		fmt.Printf("Snippet with tag '%s' copied to clipboard.\n", tag)
	},
}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a new snippet",
	Run: func(cmd *cobra.Command, args []string) {
		tag, _ := cmd.Flags().GetString("tag")
		filepath, _ := cmd.Flags().GetString("filepath")
		startLine, _ := cmd.Flags().GetString("startline")
		endLine, _ := cmd.Flags().GetString("endline")
		if tag == "" {
			fmt.Println("Tag is required")
			return
		}
		startLineInt, err := strconv.ParseInt(startLine, 10, 64)
		if err != nil {
			fmt.Println("Invalid start line")
		}
		endLineInt, err := strconv.ParseInt(endLine, 10, 64)
		if err != nil {
			fmt.Println("Invalid end line")
		}
		code, err := getCodeFromFile(filepath, startLineInt, endLineInt)
		if err != nil {
			fmt.Println(err)
			return
		}
		saveSnippet(tag, code)
	},
}

func init() {
	saveCmd.Flags().StringP("tag", "t", "", "Tag to identify the snippet")
	saveCmd.Flags().StringP("filepath", "f", "", "File to save code from")
	saveCmd.Flags().StringP("startline", "s", "", "Line to start saving code")
	saveCmd.Flags().StringP("endline", "e", "", "Line to end saving code")
	copyCmd.Flags().StringP("tag", "t", "", "Tag to identify the snippet")
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(copyCmd)
}

func getStorageFile() string {
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

func loadSnippets() SnippetStore {
	file := getStorageFile()
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return SnippetStore{Snippets: []Snippet{}}
		}
		fmt.Println("Error reading snippets file:", err)
		os.Exit(1)
	}
	var store SnippetStore
	err = json.Unmarshal(data, &store)
	if err != nil {
		fmt.Println("Error parsing snippets file:", err)
		os.Exit(1)
	}
	return store
}

func saveSnippet(tag, code string) {
	existingSnippets := loadSnippets()
	for _, s := range existingSnippets.Snippets {
		if s.Tag == tag {
			if s.Tag == tag {
				fmt.Printf("Tag '%s' already exists. Please use a different tag.\n", tag)
				return
			}
		}
	}
	newSnippet := Snippet{
		Tag:       tag,
		Code:      code,
		CreatedAt: time.Now(),
	}

	existingSnippets.Snippets = append(existingSnippets.Snippets, newSnippet)
	data, err := json.MarshalIndent(existingSnippets, "", "  ")
	if err != nil {
		fmt.Println("Error encoding snippets:", err)
		return
	}

	err = os.WriteFile(getStorageFile(), data, 0644)
	if err != nil {
		fmt.Println("Error writing snippets to file:", err)
		return
	}
	fmt.Printf("Snippet saved successfully with tag '%s'\n", tag)
}

func getCodeFromFile(filePath string, startLine, endLine int64) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
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
	code := strings.Join(selectedLines, "\n")
	return code, nil
}

func copySnippet(tag string) (Snippet, error) {
	store := loadSnippets()
	for _, s := range store.Snippets {
		if s.Tag == tag {
			return s, nil
		}
	}
	return Snippet{}, fmt.Errorf("snippet with tag '%s' not found", tag)
}
func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
