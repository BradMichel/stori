package ses

import (
	"context"
	"fmt"

	"github.com/BradMichel/stori/pkg/email"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	charset = "UTF-8"
)

type Client struct {
	instance *ses.SES
}

func New(instance *ses.SES) *Client {
	return &Client{
		instance: instance,
	}
}

func (c *Client) Send(ctx context.Context, e email.Email) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(e.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(e.Content),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(e.Subject),
			},
		},
		Source: aws.String(e.From),
	}

	_, err := c.instance.SendEmailWithContext(ctx, input)
	if err != nil {
		err = fmt.Errorf("error sending email: %#v %w", e, err)
	}

	return err
}
