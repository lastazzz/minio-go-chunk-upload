package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"minio-go-chunk-upload/handler/request"
	"net/http"
)

func (h *Handler) AbortMultipartUpload(c *gin.Context) {
	var req request.AbortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		return
	}
	if err := h.oss.AbortMultipartUpload(c, BucketName, req.Filename, req.UploadId); err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
