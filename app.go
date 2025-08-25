package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	hook "github.com/robotn/gohook"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	THOUGHT_DELIMITER = "!"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.listenHotkeys()
}

func (a *App) listenHotkeys() {
	fmt.Println("Listening for hotkeys...")

	hook.Register(hook.KeyDown, []string{"ctrl", "o"}, func(e hook.Event) {
		rt.WindowShow(a.ctx)
		rt.WindowSetAlwaysOnTop(a.ctx, true)
		rt.WindowSetAlwaysOnTop(a.ctx, false)
	})

	hook.Register(hook.KeyDown, []string{"esc"}, func(e hook.Event) {
		rt.WindowHide(a.ctx)
	})

	s := hook.Start()
	<-hook.Process(s)
}

func (a *App) GetFilePath() string {
	homeDir, _ := os.UserHomeDir()
	dirName := ".muninn"
	if runtime.GOOS == "windows" {
		dirName = "Muninn"
	}
	currentDateStr := time.Now().Format("2006-01-02")
	return path.Join(homeDir, dirName, currentDateStr+".json")
}

func (a *App) SaveNote(thought Thought) error {
	if err := a.SaveAttachment(&thought); err != nil {
		return fmt.Errorf("failed to save attachments: %w", err)
	}

	filePath := a.GetFilePath()

	// ReadFile doesn't error if file doesn't exist, just returns empty
	var thoughts []Thought
	if data, err := os.ReadFile(filePath); err == nil && len(data) > 0 {
		json.Unmarshal(data, &thoughts)
	}

	thoughts = append(thoughts, thought)

	thoughtsJSON, err := json.MarshalIndent(thoughts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal thoughts: %w", err)
	}

	// WriteFile creates the file if it doesn't exist
	return os.WriteFile(filePath, thoughtsJSON, 0644)
}

func (a *App) SaveAttachment(thought *Thought) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	dirName := ".muninn"
	if runtime.GOOS == "windows" {
		dirName = "Muninn"
	}

	dirPath := filepath.Join(homeDir, dirName, "attachments")
	currentDateStr := time.Now().Format("2006-01-02")
	filePath := filepath.Join(dirPath, currentDateStr, thought.ID)

	if err := os.MkdirAll(filePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", filePath, err)
	}

	for i, att := range thought.Attachments {
		if att.FileName == "" {
			att.FileName = fmt.Sprintf("attachment_%d", i)
		}

		destPath := filepath.Join(filePath, att.FileName)

		// Save file to disk
		if len(att.Data) > 0 {
			if err := os.WriteFile(destPath, att.Data, 0644); err != nil {
				return fmt.Errorf("failed to save attachment %s: %w", att.FileName, err)
			}
		} else if att.FilePath != "" {
			if err := copyFile(att.FilePath, destPath); err != nil {
				return fmt.Errorf("failed to copy attachment %s: %w", att.FileName, err)
			}
		}

		// Update the attachment with file path and clear binary data
		thought.Attachments[i].FilePath = filepath.Join("attachments", currentDateStr, thought.ID, att.FileName)
		thought.Attachments[i].Data = nil // Clear binary data before saving to JSON
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return destFile.Sync()
}
