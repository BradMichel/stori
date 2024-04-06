package v1

import (
	"embed"

	"github.com/BradMichel/stori/internal/notifiers"
	"github.com/BradMichel/stori/internal/summarizers"
	"github.com/BradMichel/stori/pkg/validators"
	"github.com/kelseyhightower/envconfig"
	"github.com/mcuadros/go-defaults"
)

type EmailNotificator = *notifiers.EmailNotificator[summarizers.AccountNotification]

type EmailBuilder = notifiers.EmailBuilder[summarizers.AccountNotification]

func NewValidatorConfig() (validators.Config, error) {
	c := validators.Config{}
	defaults.SetDefaults(&c)
	err := envconfig.Process("notificator", &c)
	return c, err
}

func NewEmailNotificator(
	config notifiers.EmailNotificatorConfig,
	builder EmailBuilder,
	client notifiers.EmailSender,
) EmailNotificator {
	return notifiers.NewEmailNotificator[summarizers.AccountNotification](config, builder, client)
}

func ProvideRequestJSON() embed.FS {
	return requestJSON
}
