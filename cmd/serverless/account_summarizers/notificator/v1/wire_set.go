//go:build wireinject

package v1

import (
	"context"

	s3m "github.com/BradMichel/stori/pkg/aws/s3"
	"github.com/BradMichel/stori/pkg/aws/s3/downloader"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/BradMichel/stori/internal/notifiers"
	"github.com/BradMichel/stori/internal/notifiers/account_summaries"
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
	account_summaries.NewEmailBuilder,
	wire.Bind(new(EmailBuilder), new(*account_summaries.EmailBuilder)),
	ses2.New,
	wire.Bind(new(notifiers.EmailSender), new(*ses2.Client)),
	NewValidatorConfig,
	account_summaries.NewEmailBuilderConfig,
	account_summaries.NewBlueAccountEmailTemplate,
	account_summaries.NewBlackCardEmailTemplate,
	account_summaries.NewGreenCardEmailTemplate,
	NewTemplatesBucketConfig,
	validators.New,
	wire.Bind(new(Validator), new(*validators.Service)),
	ProvideRequestJSON,
	session.NewSession,
	aws.ConfigProvider,
	wire.Bind(new(client.ConfigProvider), new(*session.Session)),
	s3.New,
	s3m.New,
	downloader.New,
	ses.New,
	logrus.New,
	wire.Bind(new(Logger), new(*logrus.Logger)),
	context.Background,
)
