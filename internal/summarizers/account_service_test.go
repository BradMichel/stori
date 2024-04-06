package summarizers_test

import (
	"context"
	"testing"

	data2 "github.com/BradMichel/stori/data"

	"github.com/stretchr/testify/assert"

	"github.com/BradMichel/stori/internal/summarizers"
	"github.com/BradMichel/stori/internal/transactions"
)

func TestAccountService_BuildSummary(t *testing.T) {
	t.Parallel()

	type Args struct {
		ctx                    context.Context
		accountID              string
		transactionList        []transactions.Transaction
		accountSummaryExpected summarizers.AccountSummary
		expectedErr            error
	}

	accountID := "1"

	tests := []struct {
		name string
		args Args
	}{
		{
			name: "build summary success",
			args: Args{
				ctx:                    context.Background(),
				accountID:              accountID,
				transactionList:        data2.GetTransactions(accountID),
				accountSummaryExpected: data2.GetAccountSummary(accountID),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			args := tt.args
			service := summarizers.NewAccountService()
			summary, err := service.BuildSummary(args.ctx, args.accountID, args.transactionList)
			assert.ErrorIs(t, err, args.expectedErr)
			assert.EqualExportedValues(t, args.accountSummaryExpected, summary)
		})
	}
}
