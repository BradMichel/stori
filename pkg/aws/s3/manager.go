package s3

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/s3"
)

var ClientEmptyErr = errors.New("s3 client is empty")
var ManagerEmptyErr = errors.New("s3 manager is empty")

type Config struct {
	BucketName string `envconfig:"BUCKET" default:"local-topic" required:"true"`
	FolderName string `envconfig:"FOLDER" default:"folder" required:"true"`
}

type Manager struct {
	Client *s3.S3
	Config
}

func New(client *s3.S3, config Config) (*Manager, error) {
	if client == nil {
		return nil, ClientEmptyErr
	}

	return &Manager{
		Client: client,
		Config: config,
	}, nil
}
