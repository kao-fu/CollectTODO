package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type TodoItem struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Date        string `json:"date"`
}

type TodoTracker struct {
	Todos []TodoItem `json:"todos"`
}

var todoPattern = regexp.MustCompile(`TODO\[(\w+)\]: (.+)`)

const maxFileSize = 500 * 1024 // 500 KB

// formatSkippedFilesMarkdown returns a markdown string for skipped files
func formatSkippedFilesMarkdown(skippedFiles []string) string {
	if len(skippedFiles) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n# Skipped Files (larger than 500 KB)\n\n")
	for _, file := range skippedFiles {
		b.WriteString(fmt.Sprintf("- %s\n", file))
	}
	b.WriteString("\n")
	return b.String()
}

func shouldFileScan(path string) bool {
	// Normalize path to lowercase for case-insensitive comparison, and get base name and extension
	baseName := strings.ToLower(filepath.Base(path))
	fileExt := strings.ToLower(filepath.Ext(path))

	// --- 1. Exclude common hidden/system files and directories ---
	// These are typically not source code and can cause issues or are irrelevant for scanning.
	// We'll check for directories like .git earlier in the file walk, but this is a good
	// redundant check for files.
	excludedNames := map[string]bool{
		".git":        true, // .git folder content
		".vscode":     true, // VS Code config folder
		".idea":       true, // IntelliJ IDEA config folder
		".ds_store":   true, // macOS system file
		"thumbs.db":   true, // Windows system file
		"__pycache__": true, // Python cache directory
	}
	if excludedNames[baseName] || strings.HasPrefix(baseName, ".") {
		// If it's a dotfile (like .gitignore, .env) or an explicitly excluded name
		// For .gitignore, .env, etc., we generally don't expect TODOs.
		return false
	}

	// --- 2. Check for specific filenames without extensions ---
	// These are common source/config files that don't always have extensions.
	if baseName == "dockerfile" || baseName == "makefile" {
		return true
	}

	// --- 3. Check for common source code and configuration file extensions ---
	// This map allows for quick O(1) lookup.
	targetExtensions := map[string]bool{
		// Primary Source Code & Scripting Files
		".c": true, ".h": true, ".cpp": true, ".hpp": true,
		".java": true, ".js": true, ".jsx": true, ".ts": true,
		".tsx": true, ".py": true, ".go": true, ".cs": true,
		".rb": true, ".php": true, ".sh": true, ".swift": true,
		".kt": true, ".rs": true,

		// Configuration & Markup Files
		".html": true, ".htm": true, ".css": true, ".scss": true,
		".less": true, ".xml": true, ".yaml": true, ".yml": true,
		".md": true, ".sql": true,
	}

	return targetExtensions[fileExt]
}

// scanTodos walks the directory tree and collects TODOs, also recording skipped files
func scanTodos(root string) ([]TodoItem, []string, error) {
	var todos []TodoItem
	var skippedFiles []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || (!shouldFileScan(path)) {
			return nil
		}

		info, err := os.Stat(path)
		if err == nil && info.Size() > int64(maxFileSize) {
			skippedFiles = append(skippedFiles, path)
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, maxFileSize)
		scanner.Buffer(buf, maxFileSize)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if matches := todoPattern.FindStringSubmatch(line); matches != nil {
				todos = append(todos, TodoItem{
					Tag:         matches[1],
					Description: matches[2],
					File:        path,
					Line:        lineNum,
				})
			}
		}
		return scanner.Err()
	})
	return todos, skippedFiles, err
}

func loadTracker(path string) (TodoTracker, error) {
	var tracker TodoTracker
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return tracker, nil
		}
		return tracker, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	if err := dec.Decode(&tracker); err != nil {
		return tracker, nil // fallback to empty
	}
	return tracker, nil
}

func saveTracker(path string, tracker TodoTracker) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(tracker)
}

func updateTodos(old []TodoItem, found []TodoItem, now string) []TodoItem {
	// Map old todos by tag+desc+file+line
	oldMap := make(map[string]TodoItem)
	for _, t := range old {
		key := fmt.Sprintf("%s|%s|%s|%d", t.Tag, t.Description, t.File, t.Line)
		oldMap[key] = t
	}
	var updated []TodoItem
	for _, t := range found {
		key := fmt.Sprintf("%s|%s|%s|%d", t.Tag, t.Description, t.File, t.Line)
		if oldT, ok := oldMap[key]; ok {
			t.Date = oldT.Date
		} else {
			t.Date = now
		}
		updated = append(updated, t)
	}
	return updated
}

func writeMarkdownToStdout(todos []TodoItem) {
	var contentBuilder strings.Builder
	contentBuilder.WriteString("# TODO Summary\n\n")
	if len(todos) == 0 {
		contentBuilder.WriteString("No TODOs found.\n")
	} else {
		tagMap := make(map[string][]TodoItem)
		for _, t := range todos {
			tagMap[t.Tag] = append(tagMap[t.Tag], t)
		}
		tags := make([]string, 0, len(tagMap))
		for tag := range tagMap {
			tags = append(tags, tag)
		}
		sort.Strings(tags)
		for _, tag := range tags {
			contentBuilder.WriteString(fmt.Sprintf("## %s\n\n", tag))
			items := tagMap[tag]
			sort.Slice(items, func(i, j int) bool {
				return items[i].Date < items[j].Date
			})
			for _, t := range items {
				contentBuilder.WriteString(fmt.Sprintf("- **%s** (%s:%d, %s): %s\n", t.Date, filepath.Base(t.File), t.Line, t.File, t.Description))
			}
			contentBuilder.WriteString("\n")
		}
	}
	newSummary := contentBuilder.String()
	fmt.Println(newSummary)
}

func main() {
	// Define the --root flag
	root := flag.String("root", ".", "Root directory to scan")

	// Parse the flags from command line
	flag.Parse()

	trackerPath := "todo_tracker.json"
	now := time.Now().Format("2006-01-02")

	found, skippedFiles, err := scanTodos(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning todos: %v\n", err)
		os.Exit(1)
	}

	tracker, _ := loadTracker(trackerPath)
	updated := updateTodos(tracker.Todos, found, now)
	if err := saveTracker(trackerPath, TodoTracker{Todos: updated}); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving tracker: %v\n", err)
		os.Exit(1)
	}

	writeMarkdownToStdout(updated)

	fmt.Print(formatSkippedFilesMarkdown(skippedFiles))
}
