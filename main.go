package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func listS3Objects(w http.ResponseWriter, r *http.Request) {
	// Create a new session with production S3 endpoint
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "prod",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})

	if err != nil {
		log.Printf("Failed to create session: %v", err)
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// List objects in the specified bucket
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("journal-react-app"), // Replace with your production bucket name
	})
	if err != nil {
		log.Printf("Failed to list objects: %v", err)
		http.Error(w, "Failed to list objects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Output the objects' details
	for _, item := range resp.Contents {
		fmt.Fprintf(w, "Name: %s, Last modified: %s\n", *item.Key, *item.LastModified)
	}
}

func main() {
	http.HandleFunc("/list-s3", listS3Objects)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
