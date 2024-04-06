//go:build wireinject

package v1

import (
	"github.com/BradMichel/stori/internal/notifiers"
	ses2 "github.com/BradMichel/stori/pkg/email/ses"
	"github.com/BradMichel/stori/pkg/validators"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/BradMichel/stori/pkg/aws"
)

var stdSet = wire.NewSet(
	New,
	NewEmailNotificator,
	wire.Bind(new(Notifier), new(EmailNotificator)),
	notifiers.NewEmailNotificatorConfig,
	notifiers.NewAccountSummaryEmailBuilder,
	wire.Bind(new(EmailBuilder), new(*notifiers.AccountSummaryEmailBuilder)),
	ses2.New,
	wire.Bind(new(notifiers.EmailSender), new(*ses2.Client)),
	NewValidatorConfig,
	validators.New,
	wire.Bind(new(Validator), new(*validators.Service)),
	ProvideRequestJSON,
	session.NewSession,
	aws.ConfigProvider,
	wire.Bind(new(client.ConfigProvider), new(*session.Session)),
	ses.New,
	logrus.New,
	wire.Bind(new(Logger), new(*logrus.Logger)),
)
