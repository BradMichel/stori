package mockers

import (
	"context"

	"github.com/BradMichel/stori/internal/transactions"

	"github.com/stretchr/testify/mock"

	"github.com/BradMichel/stori/internal/accounts"
	"github.com/BradMichel/stori/internal/summarizers"
)

const (
	produceSummaryMethod = "ProduceSummary"
	buildSummaryMethod   = "BuildSummary"
	notifyMethod         = "Notify"
)

type SummarizerMocker struct {
	mock.Mock
}

func (m *SummarizerMocker) MockProduceSummary(
	ctx context.Context,
	account accounts.Account,
	summary summarizers.AccountSummary,
	err error,
	times int,
) {
	m.On(produceSummaryMethod, ctx, account).Return(summary, err).Times(times)
}

func (m *SummarizerMocker) ProduceSummary(ctx context.Context, account accounts.Account) (summarizers.AccountSummary, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(summarizers.AccountSummary), args.Error(1)
}

func (m *SummarizerMocker) MockBuildSummary(
	ctx context.Context,
	accountID string,
	transactions []transactions.Transaction,
	summary summarizers.AccountSummary,
	err error,
	times int,
) {
	m.On(buildSummaryMethod, ctx, accountID, transactions).Return(summary, err).Times(times)
}

func (m *SummarizerMocker) BuildSummary(
	ctx context.Context,
	accountID string,
	transactions []transactions.Transaction,
) (summarizers.AccountSummary, error) {
	args := m.Called(ctx, accountID, transactions)
	return args.Get(0).(summarizers.AccountSummary), args.Error(1)
}

func (m *SummarizerMocker) MockNotify(
	ctx context.Context,
	notification summarizers.AccountNotification,
	err error,
	times int,
) {
	m.On(notifyMethod, ctx, notification).Return(err).Times(times)
}

func (m *SummarizerMocker) Notify(ctx context.Context, notification summarizers.AccountNotification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}
