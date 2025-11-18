package storage

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioRepo struct {
	client     *minio.Client
	bucketName string
}

func NewMinioRepo(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioRepo, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioRepo{client: minioClient, bucketName: bucket}, nil
}

// PresignObject creates a presigned GET URL for the object
func (m *MinioRepo) PresignObject(ctx context.Context, objectKey string, expirySeconds int) (string, error) {
	reqParams := make(url.Values)
	expiry := time.Duration(expirySeconds) * time.Second
	u, err := m.client.PresignedGetObject(ctx, m.bucketName, objectKey, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
