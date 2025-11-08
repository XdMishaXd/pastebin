package minioStorage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client *minio.Client
	bucket string
}

func New(endpoint, accessKey, secretKey, bucket string, ctx context.Context, useSSL bool) (*MinIOStorage, error) {
	const op = "minio.New"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
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

// * SaveStringAsFile сохраняет текст в storage
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

// * GetString возвращает содержимое объекта как строку.
func (m *MinIOStorage) GetString(ctx context.Context, hash string) (string, error) {
	const op = "minio.GetString"

	obj, err := m.client.GetObject(ctx, m.bucket, hash+".txt", minio.GetObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, obj); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return buf.String(), nil
}

// * DeleteFile удаляет объект по хэшу.
func (m *MinIOStorage) DeleteFile(ctx context.Context, hash string) error {
	const op = "minio.DeleteFile"

	err := m.client.RemoveObject(ctx, m.bucket, hash+".txt", minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// * ListFiles возвращает все объекты в бакете.
func (m *MinIOStorage) ListFiles(ctx context.Context) ([]string, error) {
	const op = "minio.ListFiles"

	var files []string
	for obj := range m.client.ListObjects(ctx, m.bucket, minio.ListObjectsOptions{Recursive: true}) {
		if obj.Err != nil {
			return nil, fmt.Errorf("%s: %w", op, obj.Err)
		}
		files = append(files, obj.Key)
	}

	return files, nil
}
