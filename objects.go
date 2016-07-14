package main

import (
	"fmt"
	"github.com/minio/minio-go"
	"io"
	"log"
	"net/http"
)

var conn *minio.Client

func connectDHO() (*minio.Client, error) {
	ssl := true

	var err error
	err = nil

	if conn == nil {
		conn, err = minio.NewV2(params.Endpoint, params.Access, params.Secret, ssl)
	}
	return conn, err
}

func objectList() (map[string]bool, error) {
	isRecursive := true

	doneCh := make(chan struct{})
	defer close(doneCh)

	minioClient, err := connectDHO()
	if err != nil {
		return nil, err
	}

	receiver := minioClient.ListObjects(params.bucket, params.prefix, isRecursive, doneCh)

	filenames := make(map[string]bool)
	for info := range receiver {
		filenames[info.Key] = true
	}
	return filenames, nil
}

//func uploadImages(original string, resized *os.File) error {
func uploadImages(r io.ReadSeeker, name string) {
	minioClient, err := connectDHO()
	if err != nil {
		log.Printf("  **there was an error connecting to DO %s: ", err)
		return
	}

	buffer := make([]byte, 512)
	_, err = r.Read(buffer)
	if err != nil {
		log.Printf("  **there was an error reading from %s: %s", name, err)
	}

	r.Seek(0, 0)
	mime := http.DetectContentType(buffer)

	objName := pathToObject(name)
	fmt.Printf("  Uploading %s:%s\n", params.bucket, objName)
	_, err = minioClient.PutObject(params.bucket, objName, r, mime)
	if err != nil {
		log.Printf("  **there was an error uploading %s: %s", objName, err)
		return
	}

	//TODO: change permissions on object to public-read
}

// Remove after setting up repo
func test_dho() error {
	ssl := false
	fmt.Println("Testing connection to DHO")
	_ = "breakpoint"
	minioClient, err := minio.NewV2(params.Endpoint, params.Access, params.Secret, ssl)
	if err != nil {
		return err
	}

	//Other libraries...
	//svc := s3.New(session.New(), &aws.Config{Endpoint: aws.String(params.Endpoint)})

	//resp, err := svc.HeadBucket(&s3.HeadBucketInput{
	//Bucket: aws.String(params.bucket),
	//})

	//auth := aws.Auth{
	//AccessKey: params.Access,
	//SecretKey: params.Secret,
	//}
	//conn := s3.New(auth, aws.Region{Name: "dho", S3Endpoint: params.Endpoint})
	//b := conn.Bucket("kmarsh")
	//res, err := b.List("", "", "", 1000)

	fmt.Println("Bucket exists?")
	err = minioClient.BucketExists(params.bucket)
	if err != nil {
		fmt.Println("error checking bucket")
		return err
	}
	fmt.Println("bucket exists")
	return nil
}
