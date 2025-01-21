package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: ", err)
		log.Println("Using environment variables instead")
	}

	s3BucketName := os.Getenv("S3_BUCKET_NAME")
	gcloudBucketName := os.Getenv("GCLOUD_BUCKET_NAME")

	s3Client, err := createS3Client()
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = ensureS3Bucket(s3Client, s3BucketName)
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = migrateFiles(gcloudBucketName, s3Client, s3BucketName)
	if err != nil {
		log.Fatalln(err)
	}
}

func createS3Client() (*minio.Client, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")

	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
}

func ensureS3Bucket(s3Client *minio.Client, s3BucketName string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	exists, errBucketExists := s3Client.BucketExists(ctx, s3BucketName)
	if errBucketExists != nil {
		return errBucketExists
	}

	if exists {
		return nil
	}

	err := s3Client.MakeBucket(ctx, s3BucketName, minio.MakeBucketOptions{})
	return err
}

func migrateFiles(gcloudBucketName string, s3Client *minio.Client, s3BucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	it := client.Bucket(gcloudBucketName).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Bucket(%q).Objects: %w", gcloudBucketName, err)
		}

		reader, _ := client.Bucket(gcloudBucketName).Object(attrs.Name).NewReader(ctx)

		info, err := s3Client.PutObject(ctx, s3BucketName, attrs.Name, reader, attrs.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Successfully uploaded %s of size %d\n", attrs.Name, info.Size)
	}

	return nil
}
