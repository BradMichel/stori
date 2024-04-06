package downloader

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	s3m "github.com/BradMichel/stori/pkg/aws/s3"
)

type Downloader struct {
	*s3m.Manager
}

func (downloader *Downloader) DownloadWithContext(_ aws.Context, key string) ([]byte, error) {
	result, err := downloader.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(string(downloader.BucketName)),
		Key:    aws.String(fmt.Sprintf("%s/%s", downloader.FolderName, key)),
	})
	if err != nil {
		return nil, fmt.Errorf(
			"error downloading object: %w, bucket: %s, folder: %s, key: %s",
			err, downloader.BucketName, downloader.FolderName, key,
		)
	}

	defer result.Body.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, result.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func New(manager *s3m.Manager) (*Downloader, error) {
	if manager == nil {
		return nil, s3m.ManagerEmptyErr
	}

	return &Downloader{
		manager,
	}, nil
}
