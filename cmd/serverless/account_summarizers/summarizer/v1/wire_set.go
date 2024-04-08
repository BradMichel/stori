//go:build wireinject

package v1

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	sns2 "github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/BradMichel/stori/internal/summarizers"
	"github.com/BradMichel/stori/internal/transactions"
	"github.com/BradMichel/stori/pkg/aws"
	s3m "github.com/BradMichel/stori/pkg/aws/s3"
	"github.com/BradMichel/stori/pkg/aws/s3/downloader"
	"github.com/BradMichel/stori/pkg/aws/sns"
	"github.com/BradMichel/stori/pkg/validators"
)

var stdSet = wire.NewSet(
	summarizers.NewAccountUseCase,
	wire.Bind(new(Summarizer), new(*summarizers.AccountUseCase)),
	transactions.NewS3Repository,
	wire.Bind(new(summarizers.TransactionProvider), new(*transactions.S3Repository)),
	summarizers.NewAccountService,
	wire.Bind(new(summarizers.SummaryBuilder), new(*summarizers.AccountService)),
	transactions.NewDynamoRepository,
	wire.Bind(new(summarizers.TransactionRepository), new(*transactions.DynamoRepository)),
	transactions.NewDynamoConfig,
	transactions.NewDynamoDep,
	dynamodb.New,
	sns.New,
	wire.Bind(new(summarizers.Notifier), new(*sns.Notifier)),
	NewNotifierConfig,
	NewValidatorConfig,
	validators.New,
	wire.Bind(new(Validator), new(*validators.Service)),
	ProvideRequestJSON,
	session.NewSession,
	aws.ConfigProvider,
	wire.Bind(new(client.ConfigProvider), new(*session.Session)),
	sns2.New,
	NewTransactionsBucketConfig,
	s3.New,
	s3m.New,
	downloader.New,
	logrus.New,
	wire.Bind(new(Logger), new(*logrus.Logger)),
	New,
)
