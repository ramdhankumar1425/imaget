package model

import "image"

type JobMeta struct {
	ID          string
	Status      string
	Type        string
	FileSizeMB  float64
	FileType    string
	FileExt     string
	ResultURL   string
	Error       string
	QueuedAt    int64
	StartedAt   int64
	CompletedAt int64
}

type Job struct {
	Meta   JobMeta
	Image  image.Image
	Params TransformParams
}
