package minioStorage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client *minio.Client
	bucket string
}

func New(endpoint, accessKey, secretKey, bucket string, ctx context.Context) (*MinIOStorage, error) {
	const op = "minio.New"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		if err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &MinIOStorage{
		client: minioClient,
		bucket: bucket,
	}, nil
}

func (m *MinIOStorage) SaveStringAsFile(ctx context.Context, hash, content string) error {
	const op = "minio.SaveStringAsFile"

	data := bytes.NewReader([]byte(content))

	_, err := m.client.PutObject(ctx, m.bucket, hash+".txt", data, int64(data.Len()), minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
