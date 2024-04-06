package infraestructure

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func CreateInitiatorNotificator(sess *session.Session, tableName string) (*sns.SNS, *string, *sqs.SQS, *sqs.CreateQueueOutput, error) {
	svc := sns.New(sess)
	topic, err := svc.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(tableName + "-topic"),
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	topicArn := topic.TopicArn
	svcSQS := sqs.New(sess)

	createQueueOutput, err := svcSQS.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(tableName + "-queue"),
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	queueAttributes, err := svcSQS.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl:       createQueueOutput.QueueUrl,
		AttributeNames: []*string{aws.String("QueueArn")},
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	queueArn := queueAttributes.Attributes["QueueArn"]
	_, err = svc.Subscribe(&sns.SubscribeInput{
		Endpoint:              aws.String(*queueArn),
		Protocol:              aws.String("sqs"),
		ReturnSubscriptionArn: aws.Bool(true),
		TopicArn:              topicArn,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return svc, topicArn, svcSQS, createQueueOutput, nil

}
