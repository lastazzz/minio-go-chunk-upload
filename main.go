package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"minio-go-chunk-upload/handler"
	"minio-go-chunk-upload/oss"
)

func main() {
	s, err := oss.NewOSS()
	if err != nil {
		log.Println(err)
		return
	}
	h := handler.NewHandler(s)

	svr := gin.Default()
	svr.POST("prepare_multipart_upload", h.PrepareMultipartUpload)
	svr.POST("complete_multipart_upload", h.CompleteMultipartUpload)
	svr.POST("abort_multipart_upload", h.AbortMultipartUpload)
	svr.POST("list_uploaded_parts_number", h.ListUploadedPartsNumber)

	if err := svr.Run(":8080"); err != nil {
		log.Println(err)
		return
	}
}
