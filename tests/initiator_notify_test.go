package tests_test

import (
	"context"
	"fmt"
	"github.com/BradMichel/stori/infraestructure"
	sqs2 "github.com/BradMichel/stori/pkg/aws/sqs"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/BradMichel/stori/internal/accounts"
	sns2 "github.com/BradMichel/stori/pkg/aws/sns"
	"github.com/BradMichel/stori/pkg/time"
)

func TestSNStoSQS(t *testing.T) {
	ctx := context.Background()
	awsConfig := GetAWSConfig()
	sess := session.Must(session.NewSession(&awsConfig))
	svc, topicArn, svcSQS, createQueueOutput, err := infraestructure.CreateInitiatorNotificator(sess, "stori-summarize-initiator")
	if err != nil {
		t.Fatal(err)
	}

	notifierConfig := sns2.NotifierConfig{TopicArn: *topicArn}
	notifier := sns2.New(svc, notifierConfig)
	account := accounts.Account{
		ID:     "test_id_1",
		Email:  "test_email_1",
		Period: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	err = notifier.Notify(ctx, account)

	result, err := svcSQS.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            createQueueOutput.QueueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Messages) == 0 {
		t.Fatalf("expected 1 message, got %d", len(result.Messages))
	}

	var accountBody accounts.Account
	_, err = sqs2.Unmarshal(*result.Messages[0].Body, &accountBody)
	if fmt.Sprint(account) != fmt.Sprint(accountBody) {
		t.Errorf("Expected message %v, got %v \n with message %v", account, accountBody, *result.Messages[0].Body)
	}
}
