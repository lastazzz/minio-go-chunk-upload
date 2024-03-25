package handler

import (
	"minio-go-chunk-upload/oss"
	"time"
)

const (
	BucketName    = "testbucket"
	ChunkPartSize = 5 * 1024 * 1024
	Expires       = 2 * time.Hour
)

type Handler struct {
	oss *oss.OSS
}

func NewHandler(s *oss.OSS) *Handler {
	return &Handler{oss: s}
}
