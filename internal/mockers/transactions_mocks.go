package mockers

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/BradMichel/stori/internal/transactions"
)

const (
	getByAccountIDMethod = "GetByAccountID"
	saveMethod           = "Save"
)

type TransactionRepMocker struct {
	mock.Mock
}

func (m *TransactionRepMocker) MockGetByAccountID(
	ctx context.Context,
	accountID string,
	transactions []transactions.Transaction,
	err error,
	times int,
) {
	m.On(getByAccountIDMethod, ctx, accountID).Return(transactions, err).Times(times)
}

func (m *TransactionRepMocker) GetByAccountID(ctx context.Context, accountID string) ([]transactions.Transaction, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).([]transactions.Transaction), args.Error(1)
}

func (m *TransactionRepMocker) MockSave(
	ctx context.Context,
	accountID string,
	transactions []transactions.Transaction,
	err error,
	times int,
) {
	m.On(saveMethod, ctx, accountID, transactions).Return(err).Times(times)
}

func (m *TransactionRepMocker) Save(ctx context.Context, accountID string, transactions []transactions.Transaction) error {
	args := m.Called(ctx, accountID, transactions)
	return args.Error(0)
}
