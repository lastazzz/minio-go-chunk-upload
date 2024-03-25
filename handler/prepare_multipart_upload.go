package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"log"
	"minio-go-chunk-upload/handler/request"
	"net/http"
)

func (h *Handler) PrepareMultipartUpload(c *gin.Context) {
	var req request.PrepareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		return
	}

	// 检查一下是否已经上传过
	existed, err := h.oss.IsObjectExisted(c, BucketName, req.Filename)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if existed {
		c.JSON(http.StatusOK, gin.H{
			"msg": "该文件已上传",
		})
		return
	}

	log.Println("该文件未曾上传，或者未上传完成")
	// 获取该文件的分片上传ID: uploadId
	// 注意！ 一个Object可以同时存在多个MultipartUpload
	// 每个MultipartUploadId单独计算上传的分片，彼此不共享
	// 最后提交CompleteMultipartUpload时，后一个提交的会覆盖前一个提交的
	uploadId, err := h.oss.NewMultipartUpload(c, BucketName, req.Filename, minio.PutObjectOptions{})
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 分片总数
	partTotal := req.Size / ChunkPartSize
	log.Println("预计共", partTotal, "个分片")

	// 生成未上传分片的上传url
	// partNumber 作为分片的唯一标识，值不能等于0
	// 因此，这里从1开始分配
	urls := make([]string, 0)
	for i := 1; i <= partTotal; i++ {
		presignedUrl, err := h.oss.PresignedPutObjectPart(c, BucketName, req.Filename, uploadId, i, Expires)
		if err != nil {
			log.Println(err)
			//h.oss.AbortMultipartUpload(c, BucketName, req.Filename, uploadID)
			return
		}
		log.Printf("第 %d 个分片的上传url: %s\n", i, presignedUrl.String())
		urls = append(urls, presignedUrl.String())
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_id": uploadId,
		"urls":      urls,
	})
}
