package main

import "time"

type Attachment struct {
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FilePath string `json:"filePath"`
	Data     []byte `json:"data"`
}

type Thought struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
	Timestamp   time.Time    `json:"timestamp"`
}

// NewThought creates a new Thought with the given text and optional attachments
func NewThought(text string, attachments ...Attachment) Thought {
	return Thought{
		Text:        text,
		Attachments: attachments,
		Timestamp:   time.Now().Local(),
	}
}

// AddAttachment adds an attachment to the thought
func (t *Thought) AddAttachment(attachment Attachment) {
	t.Attachments = append(t.Attachments, attachment)
}
