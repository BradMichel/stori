package tests_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/BradMichel/stori/data"
	"github.com/BradMichel/stori/infraestructure"
	"github.com/BradMichel/stori/internal/transactions"
	s3m "github.com/BradMichel/stori/pkg/aws/s3"
	"github.com/BradMichel/stori/pkg/aws/s3/downloader"
	uploader2 "github.com/BradMichel/stori/pkg/aws/s3/uploader"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestTransactionsS3_GetByAccountID(t *testing.T) {
	ctx := context.Background()
	awsConfig := GetAWSConfig()
	mySession := session.Must(session.NewSession(&awsConfig))
	bucketName := "transactions"
	err := infraestructure.CreateTransactionsS3Bucket(mySession, "transactions")
	if err != nil {
		t.Fatal(err)
	}

	s3Client := s3.New(mySession)
	s3Config := s3m.Config{BucketName: bucketName, FolderName: "records"}
	s3Manager, err := s3m.New(s3Client, s3Config)
	if err != nil {
		t.Fatal(err)
	}

	uploader, err := uploader2.New(s3Manager)
	if err != nil {
		t.Fatal(err)
	}

	accountID := "1"
	fileName := fmt.Sprintf(transactions.TransactionFileName, accountID)
	file, err := os.Open("./data/" + fileName)
	if err != nil {
		t.Fatal(err)
	}

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	err = uploader.UploadWithContext(ctx, fileName, buffer)
	if err != nil {
		t.Fatal(err)
	}

	s3Downloader, err := downloader.New(s3Manager)
	if err != nil {
		t.Fatal(err)
	}

	transactionsS3 := transactions.NewS3Repository(s3Downloader)
	transactionList, err := transactionsS3.GetByAccountID(ctx, accountID)
	if err != nil {
		t.Fatal(err)
	}

	expectedTransactions := data.GetTransactions(accountID)
	assert.Equal(t, expectedTransactions, transactionList)
}
