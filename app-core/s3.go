package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Obj struct {
	FileName string
	//LastMod     string
	StorageType string
	ETag        string
	Url         string
}

func home(w http.ResponseWriter, r *http.Request) {
	AWS_ID := ""
	AWS_SECRET := ""
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(AWS_ID, AWS_SECRET, ""),
	})

	bucket := ""
	folderLoc := ""

	resp, err := GetObjects(sess, &bucket, &folderLoc)
	if err != nil {
		fmt.Println("Got an error retrieving buckets:")
		fmt.Println(err)
		return
	}

	var s3ObjList []s3Obj
	for _, item := range resp.Contents {
		url, _ := GetObjUrl(sess, item.Key)
		s3ObjList = append(s3ObjList, s3Obj{*item.Key, *item.StorageClass, *item.ETag, url})
	}

	b, err := json.Marshal(s3ObjList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}

func GetObjects(sess *session.Session, bucket *string, prefix *string) (*s3.ListObjectsV2Output, error) {
	// snippet-start:[s3.go.list_objects.call]
	svc := s3.New(sess)

	// Get the list of items
	resp, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: bucket,
			Prefix: prefix,
		})
	// snippet-end:[s3.go.list_objects.call]
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetObjUrl(sess *session.Session, key *string) (string, error) {
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("sima.ai"),
		Key:    key,
	})
	str, err := req.Presign(15 * time.Minute)

	return str, err
}

func main() {
	fmt.Println("FOO: ", os.Getenv("FOO"))
	//http.HandleFunc("/getObjects", home)
	//http.ListenAndServe(":8080", nil)
	/*
		downloader := s3manager.NewDownloader(sess)
		file, err := os.Create("test.txt")
		if err != nil {
			fmt.Println("Failed to create file %d", *lastObj.Key)
		}

		// Write the contents of S3 Object to the file
		n, err := downloader.Download(file, &s3.GetObjectInput{
			Bucket: aws.String("sima.ai"),
			Key:    aws.String("foobar/foobar.txt"),
		})
		if err != nil {
			fmt.Println("Failed to write to the file")
		}
		fmt.Println("File downloaded, %d bytes\n", n)
	*/
}
