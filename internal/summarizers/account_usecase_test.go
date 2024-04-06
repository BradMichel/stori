package summarizers_test

import (
	"context"
	"testing"

	data2 "github.com/BradMichel/stori/data"

	"github.com/stretchr/testify/assert"

	"github.com/BradMichel/stori/internal/accounts"
	"github.com/BradMichel/stori/internal/mockers"
	"github.com/BradMichel/stori/internal/summarizers"
	"github.com/BradMichel/stori/internal/transactions"
	"github.com/BradMichel/stori/pkg/time"
)

func TestAccountUseCase_ProduceSummary(t *testing.T) {
	t.Parallel()

	type Fields struct {
		transactions *mockers.TransactionRepMocker
		summary      *mockers.SummarizerMocker
		notifier     *mockers.SummarizerInitiatorNotifierMocker
	}

	type Args struct {
		ctx                    context.Context
		account                accounts.Account
		transactions           []transactions.Transaction
		getByAccountErr        error
		accountSummary         summarizers.AccountSummary
		summaryErr             error
		notifyErr              error
		saveErr                error
		accountSummaryExpected summarizers.AccountSummary
		expectedErr            error
	}

	GetFields := func() Fields {
		return Fields{
			transactions: new(mockers.TransactionRepMocker),
			summary:      new(mockers.SummarizerMocker),
			notifier:     new(mockers.SummarizerInitiatorNotifierMocker),
		}
	}

	Mock := func(f Fields, a Args) {
		f.transactions.MockGetByAccountID(a.ctx, a.account.ID, a.transactions, a.getByAccountErr, 1)
		if a.getByAccountErr != nil {
			return
		}

		f.summary.MockBuildSummary(a.ctx, a.account.ID, a.transactions, a.accountSummary, a.summaryErr, 1)
		if a.summaryErr != nil {
			return
		}

		notification := summarizers.AccountNotification{
			Account: a.account,
			Summary: a.accountSummary,
		}
		f.notifier.MockNotify(a.ctx, notification, a.notifyErr, 1)
		if a.notifyErr != nil {
			return
		}

		f.transactions.MockSave(a.ctx, a.account.ID, a.transactions, a.saveErr, 1)
	}

	accountID := "1"
	account := accounts.Account{
		ID:     accountID,
		Email:  "test_email",
		Period: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	transactionList := data2.GetTransactions(accountID)
	accountSummary := data2.GetAccountSummary(accountID)

	tests := []struct {
		name   string
		fields Fields
		args   Args
		mock   func(Fields, Args)
	}{
		{
			name:   "success: notify account summary with no errors",
			fields: GetFields(),
			args: Args{
				ctx:                    context.Background(),
				account:                account,
				transactions:           transactionList,
				accountSummary:         accountSummary,
				accountSummaryExpected: accountSummary,
			},
			mock: Mock,
		},
		{
			name:   "error: get error on transaction save",
			fields: GetFields(),
			args: Args{
				ctx:                    context.Background(),
				account:                account,
				transactions:           transactionList,
				saveErr:                assert.AnError,
				expectedErr:            assert.AnError,
				accountSummary:         accountSummary,
				accountSummaryExpected: summarizers.AccountSummary{},
			},
			mock: Mock,
		},
		{
			name:   "error: get error on notify",
			fields: GetFields(),
			args: Args{
				ctx:                    context.Background(),
				account:                account,
				transactions:           transactionList,
				notifyErr:              assert.AnError,
				expectedErr:            assert.AnError,
				accountSummary:         accountSummary,
				accountSummaryExpected: summarizers.AccountSummary{},
			},
			mock: Mock,
		},
		{
			name:   "error: get error on summary build",
			fields: GetFields(),
			args: Args{
				ctx:                    context.Background(),
				account:                account,
				transactions:           transactionList,
				summaryErr:             assert.AnError,
				expectedErr:            assert.AnError,
				accountSummary:         summarizers.AccountSummary{},
				accountSummaryExpected: summarizers.AccountSummary{},
			},
			mock: Mock,
		},
		{
			name:   "error: get error on get by account",
			fields: GetFields(),
			args: Args{
				ctx:                    context.Background(),
				account:                account,
				getByAccountErr:        assert.AnError,
				expectedErr:            assert.AnError,
				accountSummaryExpected: summarizers.AccountSummary{},
			},
			mock: Mock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := tt.fields
			a := tt.args
			tt.mock(f, a)

			uc := summarizers.NewAccountUseCase(f.transactions, f.summary, f.notifier, f.transactions)
			summary, err := uc.ProduceSummary(a.ctx, a.account)
			assert.ErrorIs(t, err, a.expectedErr)
			assert.Equal(t, a.accountSummaryExpected, summary)
			f.transactions.AssertExpectations(t)
			f.summary.AssertExpectations(t)
			f.notifier.AssertExpectations(t)
		})
	}
}
