package main

type UploadRequest struct {
	FilePath string `json:"filepath"`
}

type MessageRequest struct {
	Message string `json:"message"`
}
