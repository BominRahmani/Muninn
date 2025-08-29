package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// getBaseDir returns the base directory for storing notes and attachments
func getBaseDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	dirName := ".muninn"
	if runtime.GOOS == "windows" {
		dirName = "Muninn"
	}
	return filepath.Join(homeDir, dirName), nil
}

func (a *App) GetFilePath() (string, error) {
	baseDir, err := getBaseDir()
	if err != nil {
		return "", err
	}
	currentDateStr := time.Now().Format("2006-01-02")
	return filepath.Join(baseDir, currentDateStr+".json"), nil
}

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
