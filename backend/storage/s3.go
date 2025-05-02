package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"log/slog"
	"simplicity/oops"
	"strings"
)

type S3BlobStore struct {
	client *s3.Client
	bucket string
}

func NewS3BlobStore(client *s3.Client, bucket string) BlobStore {
	return &S3BlobStore{client, bucket}
}

func (s *S3BlobStore) List(ctx context.Context, prefix string, delimiter string) ([]ListResult, error) {
	var err error
	var output *s3.ListObjectsV2Output
	var result ListResult
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delimiter),
	}
	var results []ListResult
	objectPaginator := s3.NewListObjectsV2Paginator(s.client, input)
	for objectPaginator.HasMorePages() {
		output, err = objectPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, object := range output.Contents {
			result, err = toLestResult(object)
			if err != nil {
				return nil, err
			}
			results = append(results, result)

		}
	}
	return results, nil
}

func toLestResult(object types.Object) (ListResult, error) {
	if object.Key == nil {
		return ListResult{}, errors.New("object key is nil")
	}
	key := *object.Key
	if len(key) == 0 {
		return ListResult{}, errors.New("object key is empty")
	}
	return ListResult{
		IsObject: key[len(key)-1] != '/',
		Key:      key,
		Size:     int(*object.Size),
	}, nil
}

func (s *S3BlobStore) Get(ctx context.Context, key string) (io.ReadCloser, map[string]string, error) {
	slog.Info("Get", "key", key)
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") {
			return nil, nil, oops.KeyNotFound
		}
		return nil, nil, err
	}
	return output.Body, output.Metadata, nil
}

const minChunkSize = 10 * 1024 * 1024

func (s *S3BlobStore) Put(ctx context.Context, key string, reader io.Reader, metadata map[string]string) error {
	if key == "" {
		return errors.New("key is empty")
	}

	chunk, err := readChunk(reader, minChunkSize)
	if err != nil {
		return fmt.Errorf("failed to read chunk: %w", err)
	}
	if len(chunk) == 0 {
		return fmt.Errorf("chunk is empty")
	}
	if len(chunk) < minChunkSize {
		return s.PutBlob(ctx, key, chunk, metadata)
	}

	resp, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:   aws.String(s.bucket),
		Key:      aws.String(key),
		Metadata: metadata,
	})

	slog.Info("Upload multipart", "key", key, "size", len(chunk))

	var parts []types.CompletedPart
	partNumber := 1

	for {
		currentPartNumber := int32(partNumber)
		var uploadResp *s3.UploadPartOutput
		uploadResp, err = s.client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(s.bucket),
			Key:        aws.String(key),
			PartNumber: aws.Int32(currentPartNumber),
			UploadId:   resp.UploadId,
			Body:       bytes.NewReader(chunk),
		})

		if err != nil {
			return fmt.Errorf("failed to upload part %d: %w", partNumber, err)
		}

		slog.Info("Uploaded part", "key", key, "partNumber", partNumber, "size", len(chunk))

		parts = append(parts, types.CompletedPart{
			ETag:       uploadResp.ETag,
			PartNumber: aws.Int32(currentPartNumber),
		})

		partNumber++
		chunk, err = readChunk(reader, minChunkSize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	_, err = s.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(s.bucket),
		Key:      aws.String(key),
		UploadId: resp.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: parts,
		},
	})

	slog.Info("Completed multipart upload", "key", key)

	return err
}

func readChunk(reader io.Reader, chunkSize int) ([]byte, error) {
	chunk := make([]byte, chunkSize)
	n, err := io.ReadFull(reader, chunk)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return chunk[:n], nil
		}
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}
	return chunk[:n], nil
}

func (s *S3BlobStore) PutBlob(ctx context.Context, key string, data []byte, metadata map[string]string) error {
	if key == "" {
		return errors.New("key is empty")
	}
	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		Metadata:      metadata,
	}
	_, err := s.client.PutObject(ctx, input)
	slog.Info("Completed blob upload", "key", key)
	return err
}

func (s *S3BlobStore) Delete(ctx context.Context, key string) error {
	if key == "" {
		return errors.New("key is empty")
	}
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	_, err := s.client.DeleteObject(ctx, input)
	return err
}

func (s *S3BlobStore) DeleteAll(ctx context.Context, prefix string) error {
	if prefix == "" {
		return errors.New("prefix is empty")
	}
	if !strings.HasSuffix(prefix, delimiter) {
		prefix += delimiter
	}
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}
	var err error
	var output *s3.ListObjectsV2Output
	objectPaginator := s3.NewListObjectsV2Paginator(s.client, input)
	for objectPaginator.HasMorePages() {
		output, err = objectPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, object := range output.Contents {
			if object.Key == nil {
				continue
			}
			err = s.Delete(ctx, *object.Key)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
