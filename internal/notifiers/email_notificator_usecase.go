package notifiers

import (
	"context"
	"fmt"

	"github.com/BradMichel/stori/pkg/email"
	"github.com/kelseyhightower/envconfig"
	"github.com/mcuadros/go-defaults"
)

type EmailBuilder[N interface{}] interface {
	Build(ctx context.Context, from string, notification N) (email.Email, error)
}

type EmailSender interface {
	Send(ctx context.Context, e email.Email) error
}

type EmailNotificator[N interface{}] struct {
	config  EmailNotificatorConfig
	builder EmailBuilder[N]
	client  EmailSender
}

type EmailNotificatorConfig struct {
	From string `envconfig:"EMAIL_SENDER" default:"michel.en@hotmail.com" required:"true"`
}

func NewEmailNotificator[N interface{}](
	config EmailNotificatorConfig,
	builder EmailBuilder[N],
	client EmailSender,
) *EmailNotificator[N] {
	return &EmailNotificator[N]{
		config:  config,
		builder: builder,
		client:  client,
	}
}

func NewEmailNotificatorConfig() (EmailNotificatorConfig, error) {
	c := EmailNotificatorConfig{}
	defaults.SetDefaults(&c)
	err := envconfig.Process("", &c)
	return c, err
}

func (n *EmailNotificator[N]) Notify(ctx context.Context, notification N) error {
	e, err := n.builder.Build(ctx, n.config.From, notification)
	if err != nil {
		return fmt.Errorf("error building email config %#v, notification %#v: %w", n.config, notification, err)
	}

	return n.client.Send(ctx, e)
}
