package smtp

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/BradMichel/stori/pkg/email"
	"github.com/kelseyhightower/envconfig"
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	Host     string `envconfig:"SMTP_HOST" default:"smtp.gmail.com" required:"true"`
	Port     int    `envconfig:"SMTP_PORT" default:"587" required:"true"`
	Username string `envconfig:"SMTP_USERNAME" required:"true"`
	Password string `envconfig:"SMTP_PASSWORD" required:"true"`
}

type Client struct {
	config Config
	auth   smtp.Auth
}

func NewConfig() (Config, error) {
	c := Config{}
	defaults.SetDefaults(&c)
	err := envconfig.Process("", &c)
	return c, err
}

func New(config Config) *Client {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	return &Client{config: config, auth: auth}
}

func (c *Client) Send(_ context.Context, e email.Email) error {
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", c.config.Host, c.config.Port),
		c.auth,
		e.From,
		[]string{e.To},
		[]byte(e.Content),
	)
}
