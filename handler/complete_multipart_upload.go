package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"minio-go-chunk-upload/handler/request"
	"net/http"
)

func (h *Handler) CompleteMultipartUpload(c *gin.Context) {
	var req request.CompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		return
	}
	_, err := h.oss.CompleteMultipartUpload(c, BucketName, req.Filename, req.UploadId)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
