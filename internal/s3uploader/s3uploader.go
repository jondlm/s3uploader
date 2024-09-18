package s3uploader

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:embed index.html
var indexHTML []byte

func uploadToS3(bucketName, key string, data []byte) error {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)

	var l int64 = int64(len(data))

	// Upload the byte slice
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: &l,
		ContentType:   aws.String("image/png"),
	})

	if err != nil {
		return fmt.Errorf("unable to upload image to S3, %v", err)
	}

	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(indexHTML)
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the form with a maximum memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the file's contents into a byte slice
	var uploadedFileBytes []byte
	uploadedFileBytes, err = io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	key := randString(8)

	err = uploadToS3("crossplane-bucket-7qlxw", key, uploadedFileBytes)

	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	// Respond with success message
	fmt.Fprintf(w, "File %s uploaded successfully to s3 key=%s, bytes=%d", handler.Filename, key, len(uploadedFileBytes))
}
