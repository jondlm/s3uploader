package main

import "fmt"
import "log"
import "net/http"
import "github.com/jondlm/s3uploader/internal/s3uploader"

func main() {
	http.HandleFunc("/", s3uploader.HandleIndex)
	http.HandleFunc("/upload", s3uploader.HandleUpload)

	fmt.Println("server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
