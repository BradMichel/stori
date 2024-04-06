package sendgrid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/sendgrid/sendgrid-go"
	smail "github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/BradMichel/stori/pkg/email"
)

var SendEmailErr = fmt.Errorf("error when the message was send")

type Config struct {
	APIKey string `envconfig:"SENDGRID_API_KEY" default:"" required:"true"`
}

type Client struct {
	client *sendgrid.Client
	config Config
}

func NewClient(config Config) *Client {
	client := sendgrid.NewSendClient(config.APIKey)
	return &Client{client: client, config: config}
}

func NewConfig() (Config, error) {
	c := Config{}
	err := envconfig.Process("", &c)
	return c, err
}

func (c *Client) Send(ctx context.Context, e email.Email) error {
	from := smail.NewEmail("sender", e.From)
	to := smail.NewEmail("receiver", e.To)
	message := smail.NewSingleEmail(from, e.Subject, to, "", e.Content)
	response, err := c.client.SendWithContext(ctx, message)
	if err != nil {
		return fmt.Errorf("%w with: %w", SendEmailErr, err)
	}

	if response.StatusCode > http.StatusMultipleChoices {
		return fmt.Errorf("%w with: %s", SendEmailErr, response.Body)
	}

	return nil
}
