package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Thought struct {
	ID          string       `json:"id"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
	Timestamp   time.Time    `json:"timestamp"`
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
