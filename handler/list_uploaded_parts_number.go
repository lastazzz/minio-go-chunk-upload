package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"minio-go-chunk-upload/handler/request"
	"net/http"
)

func (h *Handler) ListUploadedPartsNumber(c *gin.Context) {
	var req request.ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		return
	}
	ids := make([]int, 0)
	objectParts, err := h.oss.ListObjectParts(c, BucketName, req.Filename, req.UploadId)
	if err != nil {
		log.Println(err)
		return
	}
	for _, part := range objectParts {
		ids = append(ids, part.PartNumber)
	}
	log.Println("获取到", len(objectParts), "条")

	c.JSON(http.StatusOK, gin.H{
		"ids": ids,
	})
}
