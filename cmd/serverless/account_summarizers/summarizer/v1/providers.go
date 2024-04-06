package v1

import (
	"embed"

	s3m "github.com/BradMichel/stori/pkg/aws/s3"
	"github.com/BradMichel/stori/pkg/aws/sns"
	"github.com/BradMichel/stori/pkg/validators"
	"github.com/kelseyhightower/envconfig"
	"github.com/mcuadros/go-defaults"
)

func NewNotifierConfig() (sns.NotifierConfig, error) {
	c := sns.NotifierConfig{}
	err := envconfig.Process("account_summarize_finished", &c)
	return c, err
}

func NewTransactionsBucketConfig() (s3m.Config, error) {
	c := s3m.Config{}
	err := envconfig.Process("account_transactions", &c)
	return c, err
}

func NewValidatorConfig() (validators.Config, error) {
	c := validators.Config{}
	defaults.SetDefaults(&c)
	err := envconfig.Process("summarizer", &c)
	return c, err
}

func ProvideRequestJSON() embed.FS {
	return requestJSON
}
