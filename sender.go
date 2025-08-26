package main

import (
    "archive/tar"
    "bytes"
    "compress/gzip"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "time"
)

func (a *App) SendFile() error {
    baseDir, err := getBaseDir()
    if err != nil {
        return err
    }
    currentDateStr := time.Now().Format("2006-01-02")
    attachmentsDir := filepath.Join(baseDir, "attachments", currentDateStr)

    if _, err := os.Stat(attachmentsDir); os.IsNotExist(err) {
        return fmt.Errorf("no attachments found for %s", currentDateStr)
    }

    var buf bytes.Buffer
    gzWriter := gzip.NewWriter(&buf)
    defer gzWriter.Close()
    tarWriter := tar.NewWriter(gzWriter)
    defer tarWriter.Close()

    err = filepath.Walk(attachmentsDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }

        relPath, err := filepath.Rel(attachmentsDir, path)
        if err != nil {
            return err
        }

        header, err := tar.FileInfoHeader(info, "")
        if err != nil {
            return err
        }
        header.Name = relPath

        if err := tarWriter.WriteHeader(header); err != nil {
            return err
        }

        file, err := os.Open(path)
        if err != nil {
            return err
        }
        if _, err := io.Copy(tarWriter, file); err != nil {
            file.Close()
            return err
        }
        return file.Close()
    })
    if err != nil {
        return err
    }

    resp, err := http.Post("http://127.0.0.1:8000/upload", "application/octet-stream", &buf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("upload failed: %s", resp.Status)
    }
    return nil
}
