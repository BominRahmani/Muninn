package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Attachment struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FilePath string `json:"filePath"`
	Data     []byte `json:"data"`
}

func (a *App) SaveAttachment(thought *Thought) error {
	baseDir, err := getBaseDir()
	if err != nil {
		return err
	}

	currentDateStr := time.Now().Format("2006-01-02")
	dirPath := filepath.Join(baseDir, "attachments", currentDateStr, thought.ID)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	for i, att := range thought.Attachments {
		if att.FileName == "" {
			att.FileName = fmt.Sprintf("attachment_%d", i)
		}

		destPath := filepath.Join(dirPath, att.FileName)

		if len(att.Data) > 0 {
			if err := os.WriteFile(destPath, att.Data, 0644); err != nil {
				return fmt.Errorf("failed to save attachment %s: %w", att.FileName, err)
			}
		} else if att.FilePath != "" {
			if err := copyFile(att.FilePath, destPath); err != nil {
				return fmt.Errorf("failed to copy attachment %s: %w", att.FileName, err)
			}
		}

		thought.Attachments[i].FilePath = filepath.Join("attachments", currentDateStr, thought.ID, att.FileName)
		thought.Attachments[i].Data = nil
	}

	return nil
}
