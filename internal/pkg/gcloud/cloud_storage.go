package gcloud

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"os"
	"time"
)

type FileUploaderClient struct {
	cl         *storage.Client
	projectID  string
	bucketName string
}

func NewFileUploaderClient() FileUploaderClient {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON"))))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return FileUploaderClient{
		cl:         client,
		bucketName: os.Getenv("GCLOUD_BUCKET_NAME"),
		projectID:  os.Getenv("GCLOUD_PROJECT_ID"),
	}
}

func (c *FileUploaderClient) UploadFile(file multipart.File, path string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(path).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func (c *FileUploaderClient) GetURL(path string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, path)
}
