package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	KeyID := ""
	SecretID := ""
	token := ""
	creds := credentials.NewStaticCredentials(KeyID, SecretID, token)
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
	}
	cfg := aws.NewConfig().WithRegion("us-west-2").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	i := 0
	err = svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String("gc-demos"),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		fmt.Println("Page,", i)
		i++

		for _, obj := range p.Contents {
			fmt.Println("Object:", *obj.Key)
		}
		return true
	})
	if err != nil {
		fmt.Println("failed to list objects", err)
		return
	}

	file, err := os.Open("demo.zip")
	if err != nil {
		fmt.Printf("err opening file: %s", err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	params := &s3.PutObjectInput{
		Bucket:        aws.String("gc-demos"),
		Key:           aws.String(file.Name()),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ACL:           aws.String("public-read"),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}
	fmt.Printf("response %s", awsutil.StringValue(resp))
}
