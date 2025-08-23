package main

type Attachment struct {
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FilePath string `json:"filePath"`
	Data     []byte `json:"data"`
}

type Thought struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}



