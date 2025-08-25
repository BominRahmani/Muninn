package main

import (
	"time"
)

type Attachment struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FilePath string `json:"filePath"`
	Data     []byte `json:"data"`
}

type Thought struct {
	ID          string       `json:"id"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
	Timestamp   time.Time    `json:"timestamp"`
}
