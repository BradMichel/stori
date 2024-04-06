package infraestructure

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func CreateTransactionsS3Bucket(sess *session.Session, bucketName string) error {
	s3Client := s3.New(sess)

	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: &bucketName,
	})

	return err
}
