package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Client *s3.Client

func main() {

	fmt.Print("Hello, Go")
	ctx := context.Background()
	if os.Args[1] == "" {
		fmt.Println("Usage: simples3uploader <filename>")
		return
	}
	fileName := os.Args[1]
	if len(os.Args) > 2 {
		fmt.Println("Usage: simples3uploader <filename>")
		return
	}
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

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

	tempKey := ""
	Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(tempKey),
		Body:   nil,
	})
}
