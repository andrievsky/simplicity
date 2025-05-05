package storage

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
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

func (s *S3BlobStore) Put(ctx context.Context, key string, reader io.Reader, metadata map[string]string) error {
	if key == "" {
		return errors.New("key is empty")
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		Metadata:      metadata,
	}
	_, err = s.client.PutObject(ctx, input)
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
	if !strings.HasSuffix(prefix, Delimiter) {
		prefix += Delimiter
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
