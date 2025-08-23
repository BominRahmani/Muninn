package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"runtime"
	"time"
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
	if err := InitDirectory(); err != nil {
		fmt.Printf("Failed to initialize directory: %v\n", err)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// InitDirectory ensures the temporary directory used for buffering the files before they are sent to the backend is initialized
func InitDirectory() error {
	// Use user's home directory instead of config directory for better permissions
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("The following error was encountered when checking for User Home Directory: %w", err)
	}

	// Use platform-appropriate directory name
	dirName := "Muninn"
	if runtime.GOOS == "windows" {
		dirName = "Muninn"
	} else {
		dirName = ".muninn"
	}

	dirPath := path.Join(homeDir, dirName)

	_, err = os.Stat(dirPath)
	if errors.Is(err, fs.ErrNotExist) {
		err = os.Mkdir(dirPath, 0755) // 0755 read, write, execute for owner, read/execute for group/others
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}
	}
	return nil
}

// CreateFile will create a file where the thoughts get appended onto. If no file is presented it will create one with the date as the name.
func (a *App) CreateFile() (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("The following error was encountered when checking for User Home Directory: %w", err)
	}

	// Use platform-appropriate directory name
	dirName := "Muninn"
	if runtime.GOOS == "windows" {
		dirName = "Muninn"
	} else {
		dirName = ".muninn"
	}

	dirPath := path.Join(homeDir, dirName)

	currentDate := time.Now()
	currentDateStr := currentDate.Format("2006-01-02")

	filePath := path.Join(dirPath, currentDateStr+".txt")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	return file, nil
}

func (a *App) SaveNote(thought Thought) error {
	file, err := a.CreateFile()
	if err != nil {
		return fmt.Errorf("The following error was encountered when creating a file to append thoughts to: %w", err)
	}
	defer file.Close()

	// Format the thought with delimiter and timestamp
	timestamp := thought.Timestamp.Format("15:04:05")
	noteText := fmt.Sprintf("%s [%s] %s\n", THOUGHT_DELIMITER, timestamp, thought.Text)

	if _, err := file.WriteString(noteText); err != nil {
		return fmt.Errorf("error writing thought: %w", err)
	}

	// Handle attachments if any
	if len(thought.Attachments) > 0 {
		for _, attachment := range thought.Attachments {
			// TODO: Implement attachment saving logic
			// For now, just log that we have attachments
			fmt.Printf("Attachment found: %s (%s)\n", attachment.FileName, attachment.FileType)
		}
	}

	return nil
}

// SaveTextNote is a convenience function to save a simple text note
func (a *App) SaveTextNote(text string) error {
	thought := NewThought(text)
	return a.SaveNote(thought)
}
