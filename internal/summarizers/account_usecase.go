package summarizers

import (
	"context"

	"github.com/BradMichel/stori/internal/transactions"

	"github.com/BradMichel/stori/internal/accounts"
)

type TransactionProvider interface {
	GetByAccountID(context.Context, string) ([]transactions.Transaction, error)
}

type SummaryBuilder interface {
	BuildSummary(context.Context, string, []transactions.Transaction) (AccountSummary, error)
}

type Notifier interface {
	Notify(context.Context, interface{}) error
}

type TransactionRepository interface {
	Save(context.Context, string, []transactions.Transaction) error
}

type AccountUseCase struct {
	tranProvider   TransactionProvider
	summaryBuilder SummaryBuilder
	notifier       Notifier
	tranRepo       TransactionRepository
}

func NewAccountUseCase(
	tranProvider TransactionProvider,
	summaryBuilder SummaryBuilder,
	notifier Notifier,
	tranRepo TransactionRepository,
) *AccountUseCase {
	return &AccountUseCase{
		tranProvider:   tranProvider,
		summaryBuilder: summaryBuilder,
		notifier:       notifier,
		tranRepo:       tranRepo,
	}
}

func (u *AccountUseCase) ProduceSummary(ctx context.Context, account accounts.Account) (AccountSummary, error) {
	transactionList, err := u.tranProvider.GetByAccountID(ctx, account.ID)
	if err != nil {
		return AccountSummary{}, err
	}

	summary, err := u.summaryBuilder.BuildSummary(ctx, account.ID, transactionList)
	if err != nil {
		return AccountSummary{}, err
	}

	notification := AccountNotification{
		Account: account,
		Summary: summary,
	}
	if err := u.notifier.Notify(ctx, notification); err != nil {
		return AccountSummary{}, err
	}

	if err := u.tranRepo.Save(ctx, account.ID, transactionList); err != nil {
		return AccountSummary{}, err
	}

	return summary, nil
}
