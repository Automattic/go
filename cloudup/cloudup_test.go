package cloudup_test

import (
	"log"
	"testing"

	"github.com/automattic/go/cloudup"
)

// generate basic auth token using Node:
// console.log(new Buffer("username:password").toString('base64'));

func TestUpload(t *testing.T) {
	filename := "--LOCAL FILE PATH--"
	client := cloudup.Client{}
	client.BasicToken = "--YOUR BASIC AUTH TOKEN--"

	stream, err := client.CreateStream("Test Stream")
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	item, err := client.CreateItem(stream.Id, filename, "Test Item")
	if err != nil {
		log.Fatalf("Error create item: %v", err)
	}

	err = client.UploadToS3(item)
	if err != nil {
		log.Fatalf("Error uploading to S3: %v", err)
	}

	err = client.CompleteItem(item)
	if err != nil {
		log.Fatalf("Error marking complete: %v", err)
	}

}
