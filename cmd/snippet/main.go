package main

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/hayzedd2/snippet-cli/utils"
	"github.com/spf13/cobra"
	"os"
	"time"
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved snippets",
	Run: func(cmd *cobra.Command, args []string) {
		snippets := loadSnippets()
		if len(snippets.Snippets) == 0 {
			fmt.Println("No snippets found")
			return
		}
		fmt.Println("\nAvailable Snippets:")
		fmt.Println("------------------")
		for i, snippet := range snippets.Snippets {
			fmt.Printf("%d. Tag: %s\n   Created: %s\n",
				i+1,
				snippet.Tag,
				snippet.CreatedAt.Format("2006-01-02 15:04:05"),
			)
		}
		fmt.Printf("\nTotal snippets: %d\n", len(snippets.Snippets))
		fmt.Println("\nUse 'snippet get --tag <tag_name>' to view a specific snippet")
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific snippet",
	Run: func(cmd *cobra.Command, args []string) {
		tag, _ := cmd.Flags().GetString("tag")
		if tag == "" {
			fmt.Println("Tag is required")
			return
		}
		snippet, err := getSingleSnippet(tag)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Tag: %s\nCode: %s\nCreated: %s\n",
			snippet.Tag,
			snippet.Code,
			snippet.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	},
}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a new snippet",
	Run: func(cmd *cobra.Command, args []string) {
		opts, err := utils.ParseSaveOptions(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
		code, err := utils.GetCodeFromFile(opts.FilePath, opts.StartLine, opts.EndLine)
		if err != nil {
			fmt.Println(err)
			return
		}
		saveSnippet(opts.Tag, code)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a saved snippet",
	Run: func(cmd *cobra.Command, args []string) {
		tag, _ := cmd.Flags().GetString("tag")
		if tag == "" {
			fmt.Println("Tag is required")
			return
		}
		err := deleteSnippet(tag)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Snippet with tag '%s' deleted successfully\n", tag)
	},
}

func init() {
	saveCmd.Flags().StringP("tag", "t", "", "Tag to identify the snippet")
	saveCmd.Flags().StringP("filepath", "f", "", "File to save code from")
	saveCmd.Flags().StringP("startline", "s", "", "Line to start saving code")
	saveCmd.Flags().StringP("endline", "e", "", "Line to end saving code")
	copyCmd.Flags().StringP("tag", "t", "", "Tag to identify the snippet")
	deleteCmd.Flags().StringP("tag", "t", "", "Tag to identify snippet to delete")
	getCmd.Flags().StringP("tag", "t", "", "Tag to view single snippet")
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(getCmd)
}

func loadSnippets() SnippetStore {
	file := utils.GetStorageFile()
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
			fmt.Printf("Tag '%s' already exists. Please use a different tag.\n", tag)
			return
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
	err = os.WriteFile(utils.GetStorageFile(), data, 0644)
	if err != nil {
		fmt.Println("Error writing snippets to file:", err)
		return
	}
	fmt.Printf("Snippet saved successfully with tag '%s'\n", tag)
}

func copySnippet(tag string) (Snippet, error) {
	store := loadSnippets()
	for _, s := range store.Snippets {
		if s.Tag == tag {
			return s, nil
		}
	}
	return Snippet{}, fmt.Errorf("\nsnippet with tag '%s' not found.\nUse 'snippet list' to see saved snippet tags", tag)
}

func deleteSnippet(tag string) error {
	store := loadSnippets()
	index := -1
	for i, s := range store.Snippets {
		if s.Tag == tag {
			index = i
			break
		}
	}

	// If tag wasn't found
	if index == -1 {
		return fmt.Errorf("snippet with tag '%s' not found", tag)
	}
	store.Snippets = append(store.Snippets[:index], store.Snippets[index+1:]...)
	data, err := json.MarshalIndent(store, "", "    ")
	if err != nil {
		return fmt.Errorf("error encoding snippets: %w", err)
	}

	if err := os.WriteFile(utils.GetStorageFile(), data, 0644); err != nil {
		return fmt.Errorf("error saving snippets: %w", err)
	}
	return nil
}


func getSingleSnippet(tag string) (Snippet, error) {
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
