package sns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Notifier struct {
	instance *sns.SNS
	conf     NotifierConfig
}

type NotifierConfig struct {
	TopicArn string `envconfig:"TOPIC" default:"local-topic" required:"true"`
}

func New(instance *sns.SNS, config NotifierConfig) *Notifier {
	return &Notifier{
		instance: instance,
		conf:     config,
	}
}

func (n *Notifier) Notify(ctx context.Context, message interface{}) error {
	messageString, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = n.instance.PublishWithContext(ctx, &sns.PublishInput{
		Message:  aws.String(string(messageString)),
		TopicArn: aws.String(n.conf.TopicArn),
	})

	if err != nil {
		err = fmt.Errorf("error publishing message to SNS Arn %s: %w", n.conf.TopicArn, err)
	}

	return err
}
