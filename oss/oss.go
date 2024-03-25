package oss

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// 用于创建minio客户端
const (
	Endpoint        = "play.min.io"
	AccessKeyID     = "Q3AM3UQ867SPQQA43P2F"
	SecretAccessKey = "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	Secure          = true
)

// 简单封装了一下minio-go

type OSS struct {
	*minio.Core
}

func NewOSS() (*OSS, error) {
	minioCore, err := minio.NewCore(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKeyID, SecretAccessKey, ""),
		Secure: Secure,
	})
	if err != nil {
		return nil, err
	}
	return &OSS{
		Core: minioCore,
	}, nil
}

// PresignedPutObjectPart minio-go没有提供，故参考PutObjectPart()方法实现了一个
func (s *OSS) PresignedPutObjectPart(c context.Context,
	bucketName string,
	objectName string,
	uploadId string,
	partNumber int,
	expires time.Duration) (*url.URL, error) {
	param := url.Values{
		"uploadId":   []string{uploadId},
		"partNumber": []string{strconv.Itoa(partNumber)},
	}
	return s.Presign(c, http.MethodPut, bucketName, objectName, expires, param)
}

func (s *OSS) ListObjectParts(c context.Context,
	bucketName string,
	objectName string,
	uploadId string) ([]minio.ObjectPart, error) {

	objectParts := make([]minio.ObjectPart, 0)
	partNumberMarker := 0
	for {
		partsResult, err := s.Core.ListObjectParts(c, bucketName, objectName, uploadId, partNumberMarker, 1000)
		if err != nil {
			return nil, err
		}
		objectParts = append(objectParts, partsResult.ObjectParts...)
		if !partsResult.IsTruncated {
			break
		}
	}
	return objectParts, nil
}

func (s *OSS) CompleteMultipartUpload(c context.Context, bucketName string, objectName string, uploadId string) (minio.UploadInfo, error) {
	// completeParts的ETag必须要填，所以先获取一下已上传的ObjectParts，再拿到ETag
	// 当然，在使用url上传时返回的响应头部中也能获取到分片的ETag
	// ETag包含了分片文件的MD5和partNumber
	objectParts, err := s.ListObjectParts(c, bucketName, objectName, uploadId)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	completeParts := make([]minio.CompletePart, 0)
	for _, part := range objectParts {
		completeParts = append(completeParts, minio.CompletePart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
			// 下面这些不填好像也没影响
			ChecksumCRC32:  part.ChecksumCRC32,
			ChecksumCRC32C: part.ChecksumCRC32C,
			ChecksumSHA1:   part.ChecksumSHA1,
			ChecksumSHA256: part.ChecksumSHA256,
		})
	}
	return s.Core.CompleteMultipartUpload(c, bucketName, objectName, uploadId, completeParts, minio.PutObjectOptions{})
}

func (s *OSS) IsObjectExisted(c context.Context, bucketName string, objectName string) (bool, error) {
	_, err := s.StatObject(c, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		var e minio.ErrorResponse
		if errors.As(err, &e) && e.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
