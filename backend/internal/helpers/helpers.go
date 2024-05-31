package helpers

import (
	"DiplomaV2/backend/internal/validator"
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"io"
	"net/url"
	"strconv"
	"strings"
)

func ReadString(qs url.Values, key string, defaultValue string) string {

	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

func ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func UploadFileToGCS(ctx context.Context, client *storage.Client, bucketName, objectName string, src io.Reader) error {
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(wc, src); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func NewStorageClient(ctx context.Context, keyFilePath string) (*storage.Client, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(keyFilePath))
	if err != nil {
		return nil, err
	}
	return client, nil
}
