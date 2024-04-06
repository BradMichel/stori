package v1

import (
	"github.com/BradMichel/stori/pkg/aws/sns"
	"github.com/kelseyhightower/envconfig"
)

func NewNotifierConfig() (sns.NotifierConfig, error) {
	c := sns.NotifierConfig{}
	err := envconfig.Process("summarize_initiator", &c)
	return c, err
}
