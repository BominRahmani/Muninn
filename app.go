package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	rt "github.com/wailsapp/wails/v2/pkg/runtime"
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
	a.scheduleSend()
	go a.listenHotkeys()
}

// LaunchNvimWithNote launches nvim in a terminal with the provided note content
func (a *App) LaunchNvimWithNote(noteContent string) error {
	// Resolve absolute path to nvim
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		return fmt.Errorf("nvim not found in PATH. Please install nvim first")
	}

	// Create a temporary file with the note content
	tempDir, err := os.MkdirTemp("", "muninn-note-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	tempFile := filepath.Join(tempDir, "note.md")
	if err := os.WriteFile(tempFile, []byte(noteContent), 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		if _, err := exec.LookPath("kitty"); err == nil {
			// First try to open in a new tab of an existing kitty instance
			remoteCmd := exec.Command("kitty", "@", "launch", "--type=tab", nvimPath, tempFile)
			if err := remoteCmd.Run(); err != nil {
				// If that fails, start a new kitty window
				cmd = exec.Command("kitty", nvimPath, tempFile)
			} else {
				// Successfully launched in an existing kitty tab
				// Minimize the Muninn window after launching nvim
				// rt.WindowMinimise(a.ctx)
				return nil
			}
		} else {
			// Fallback to other terminals
			terminals := []struct {
				name string
				args []string
			}{
				{"gnome-terminal", []string{"--", nvimPath, tempFile}},
				{"konsole", []string{"-e", nvimPath, tempFile}},
				{"xterm", []string{"-e", nvimPath, tempFile}},
				{"alacritty", []string{"-e", nvimPath, tempFile}},
				{"xfce4-terminal", []string{"-e", nvimPath, tempFile}},
				{"tilix", []string{"-e", nvimPath, tempFile}},
			}

			for _, term := range terminals {
				if _, err := exec.LookPath(term.name); err == nil {
					cmd = exec.Command(term.name, term.args...)
					break
				}
			}
			if cmd == nil {
				return fmt.Errorf("no supported terminal found")
			}
		}

	case "windows":
		// On Windows, "start" needs to be run via cmd
		cmd = exec.Command("cmd", "/c", "start", nvimPath, tempFile)

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Launch the terminal with nvim (if we got here, kitty remote failed or fallback was chosen)
	if cmd != nil {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to launch nvim: %w", err)
		}
	}

	// Minimize the Muninn window after launching nvim
	rt.WindowMinimise(a.ctx)

	return nil
}
