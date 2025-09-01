package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Thought struct {
	ID          string       `json:"id"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
	Timestamp   time.Time    `json:"timestamp"`
}

type SearchResult struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Text    string `json:"text"`
}

func (a *App) SaveNote(thought Thought) error {
	if err := a.SaveAttachment(&thought); err != nil {
		return fmt.Errorf("failed to save attachments: %w", err)
	}

	filePath, err := a.GetFilePath()
	if err != nil {
		return err
	}

	var thoughts []Thought
	if data, err := os.ReadFile(filePath); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &thoughts); err != nil {
			return fmt.Errorf("failed to parse existing thoughts: %w", err)
		}
	}

	thoughts = append(thoughts, thought)

	thoughtsJSON, err := json.MarshalIndent(thoughts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal thoughts: %w", err)
	}

	return os.WriteFile(filePath, thoughtsJSON, 0o644)
}

func (a *App) SearchNotes(query string) ([]SearchResult, error) {
	filePath, err := a.GetFilePath()
	if err != nil {
		return nil, err
	}

	var thoughts []Thought
	if data, err := os.ReadFile(filePath); err != nil {
		// If file doesn't exist, return empty results
		return []SearchResult{}, nil
	} else if len(data) > 0 {
		if err := json.Unmarshal(data, &thoughts); err != nil {
			return nil, fmt.Errorf("failed to parse thoughts: %w", err)
		}
	}

	query = strings.ToLower(query)
	var results []SearchResult

	for _, thought := range thoughts {
		if strings.Contains(strings.ToLower(thought.Text), query) {
			// Truncate content for display
			content := thought.Text
			if len(content) > 200 {
				content = content[:200] + "..."
			}

			results = append(results, SearchResult{
				ID:      thought.ID,
				Content: content,
				Text:    thought.Text,
			})
		}
	}

	return results, nil
}
