package uploader

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	s3m "github.com/BradMichel/stori/pkg/aws/s3"
)

type Uploader struct {
	*s3m.Manager
}

func (uploader *Uploader) UploadWithContext(_ aws.Context, key string, input []byte) error {
	_, err := uploader.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(string(uploader.BucketName)),
		Key:    aws.String(fmt.Sprintf("%s/%s", uploader.FolderName, key)),
		Body:   bytes.NewReader(input),
	})
	if err != nil {
		return err
	}
	return nil
}

func New(manager *s3m.Manager) (*Uploader, error) {
	if manager == nil {
		return nil, s3m.ManagerEmptyErr
	}
	return &Uploader{
		manager,
	}, nil
}
