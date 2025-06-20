package main

import (
	"bufio"
	"encoding/json"
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

var todoPattern = regexp.MustCompile(`// TODO\[(\w+)\]: (.+)`)

func scanTodos(root string) ([]TodoItem, error) {
	var todos []TodoItem
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || (!strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".js")) {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
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
	return todos, err
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
	root := "./"
	trackerPath := "TODO_TRACKER.json"
	// set now to tomorrow's date
	now := time.Now().Format("2006-01-02")

	found, err := scanTodos(root)
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
}
