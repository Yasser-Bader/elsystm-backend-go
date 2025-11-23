package util

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var EndPoint string

type AWSMulitableFiles struct {
	URL      string
	Filetype string
}

func ConnectAws() *session.Session {
	godotenv.Load()

	AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	MyRegion = os.Getenv("AWS_REGION")
	EndPoint = os.Getenv("AWS_END_POINT")

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"",
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}

func AWSUploadSingleFile(ctx *gin.Context, requestFile string) (string, string, error) {
	// AWS Session
	sess := ConnectAws()
	uploader := s3manager.NewUploader(sess)
	MyBucket := os.Getenv("AWS_BUCKET_NAME")

	file, header, err := ctx.Request.FormFile(requestFile)

	if err != nil {
		return "", "", err
	}

	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(MyBucket),
		ACL:         aws.String(s3.BucketCannedACLPublicRead),
		ContentType: aws.String(header.Header.Get("Content-Type")),
		Key:         aws.String(AWSRenameFilename(header.Filename)),
		Body:        file,
	})
	if err != nil {
		return "", "", nil
	}
	return up.Location, header.Header.Get("Content-Type"), nil
}

func AWSUploadMultiableFiles(ctx *gin.Context, requestFile string) []AWSMulitableFiles {
	var awsMulitableFiles []AWSMulitableFiles
	// AWS Session
	sess := ConnectAws()
	uploader := s3manager.NewUploader(sess)
	MyBucket := os.Getenv("AWS_BUCKET_NAME")

	form, _ := ctx.MultipartForm()
	files := form.File[requestFile]
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			continue
		}
		up, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(MyBucket),
			ACL:         aws.String(s3.BucketCannedACLPublicRead),
			ContentType: aws.String(file.Header.Get("Content-Type")),
			Key:         aws.String(AWSRenameFilename(file.Filename)),
			Body:        f,
		})
		if err == nil {
			awsMulitableFiles = append(awsMulitableFiles, AWSMulitableFiles{
				URL:      up.Location,
				Filetype: file.Header.Get("Content-Type"),
			})
		}
	}

	return awsMulitableFiles
}

func AWSRenameFilename(filename string) string {
	fileSplit := strings.Split(filename, ".")
	return fmt.Sprintf("%d-%s.%s", time.Now().Unix(), fileSplit[0], fileSplit[len(fileSplit)-1])
}
