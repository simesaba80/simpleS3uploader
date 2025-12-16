package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	InitS3Client() // S3クライアントを初期化
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World")
	})

	http.HandleFunc("/s3", s3Handler)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

const bucketName = "assets"

func s3Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// "file"キーからファイルを取得
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// S3にアップロード
	_, err = Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(header.Filename),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})
	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	// 署名付きURLを生成（有効期限: 15分）
	presignedReq, err := PresignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(header.Filename),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
		http.Error(w, "Failed to generate presigned URL", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Uploaded: %s", presignedReq.URL)
}
