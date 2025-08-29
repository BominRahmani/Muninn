package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
)

func (a *App) SendFile() error {
	baseDir, err := getBaseDir()
	if err != nil {
		return err
	}

	currentDateStr := time.Now().Format("2006-01-02")
	attachmentsDir := filepath.Join(baseDir, "attachments", currentDateStr)
	thoughtFile := filepath.Join(baseDir, currentDateStr+".json")

	if _, err := os.Stat(attachmentsDir); os.IsNotExist(err) {
		return fmt.Errorf("no attachments found for %s", currentDateStr)
	}

	pr, pw := io.Pipe()

	go func() {
		gzWriter := gzip.NewWriter(pw)
		tarWriter := tar.NewWriter(gzWriter)

		err := filepath.WalkDir(attachmentsDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			return addFileToTar(tarWriter, attachmentsDir, path)
		})
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		if err := addFileToTar(tarWriter, baseDir, thoughtFile); err != nil {
			pw.CloseWithError(err)
			return
		}

		if err := tarWriter.Close(); err != nil {
			pw.CloseWithError(err)
			return
		}
		if err := gzWriter.Close(); err != nil {
			pw.CloseWithError(err)
			return
		}

		pw.Close()
	}()

	resp, err := http.Post("http://127.0.0.1:8000/upload", "application/octet-stream", pr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed: %s", resp.Status)
	}
	return nil
}

// addFileToTar is a helper function that adds files into tar bundle
func addFileToTar(tw *tar.Writer, baseDir, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = relPath

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(tw, file)
	return err
}

// scheduleSend is a cron job function that will send the attachments and thought to the server at the stroke of midnight
func (a *App) scheduleSend() error {
	c := cron.New()

	_, err := c.AddFunc("0 0 0 * * *", func() {
		if err := a.SendFile(); err != nil {
			fmt.Printf("[%s] failed to send file: %v\n", time.Now().Format(time.RFC3339), err)
		} else {
			fmt.Printf("[%s] successfully sent file\n", time.Now().Format(time.RFC3339))
		}
	})
	if err != nil {
		return err
	}

	c.Start()
	return nil
}
