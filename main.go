package main

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

var Client *s3.Client

func getFileMetadata(file *os.File) (contentType string, extension string, err error) {
	extension = strings.ToLower(filepath.Ext(file.Name()))

	header := make([]byte, 512)
	n, readErr := file.ReadAt(header, 0)
	if readErr != nil && readErr != io.EOF {
		return "", "", fmt.Errorf("failed to read file header: %w", readErr)
	}

	contentType = http.DetectContentType(header[:n])
	if extensionType := mime.TypeByExtension(extension); extensionType != "" &&
		(contentType == "application/octet-stream" || strings.HasPrefix(contentType, "text/plain")) {
		contentType = extensionType
	}

	return contentType, extension, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	ctx := context.Background()
	if len(os.Args) != 2 || os.Args[1] == "" {
		fmt.Println("Usage: simples3uploader <filename>")
		return
	}
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	contentType, extension, err := getFileMetadata(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Content-Type: %s, extension: %s\n", contentType, extension)

	creds := credentials.NewStaticCredentialsProvider(os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_ACCESS_KEY"), "")
	sdkConfig, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(creds),
		config.WithRegion(os.Getenv("S3_REGION_NAME")),
		config.WithBaseEndpoint(os.Getenv("S3_BASE_URL")),
	)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	Client = s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	key := os.Getenv("S3_DIR") + "/" + uuid.New().String() + "/origin" + extension
	fmt.Println("key:", key)
	_, err = Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		fmt.Printf("Failed to upload %s: %v\n", fileName, err)
		return
	}
	uploadedURL := os.Getenv("S3_BASE_URL") + "/" + os.Getenv("S3_BUCKET_NAME") + "/" + key
	fmt.Printf("Uploaded %s as %s\n", fileName, uploadedURL)
}
