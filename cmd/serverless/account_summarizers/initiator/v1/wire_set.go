//go:build wireinject

package v1

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	sns2 "github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/BradMichel/stori/internal/accounts"
	"github.com/BradMichel/stori/pkg/aws"
	"github.com/BradMichel/stori/pkg/aws/sns"
	"github.com/BradMichel/stori/pkg/time"
)

var stdSet = wire.NewSet(
	time.New,
	wire.Bind(new(Time), new(*time.Service)),
	accounts.NewSummarizerInitiator,
	wire.Bind(new(Initiator), new(*accounts.SummarizerInitiator)),
	accounts.NewDynamoRepository,
	wire.Bind(new(accounts.Repository), new(*accounts.DynamoRepository)),
	accounts.NewDynamoConfig,
	accounts.NewDynamoDep,
	dynamodb.New,
	sns.New,
	wire.Bind(new(accounts.Notifier), new(*sns.Notifier)),
	NewNotifierConfig,
	session.NewSession,
	aws.ConfigProvider,
	wire.Bind(new(client.ConfigProvider), new(*session.Session)),
	sns2.New,
	logrus.New,
	wire.Bind(new(Logger), new(*logrus.Logger)),
	New,
)
