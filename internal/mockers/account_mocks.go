package mockers

import (
	"context"
	"github.com/BradMichel/stori/internal/accounts"
	"time"

	"github.com/stretchr/testify/mock"
)

const (
	launchMethod        = "Launch"
	getByByPeriodMethod = "GetByByListID"
	updateMethod        = "Update"
	notifierMethod      = "Notify"
)

type SummarizerInitiatorMocker struct {
	mock.Mock
}

type RepositoryMocker struct {
	mock.Mock
}

type SummarizerInitiatorNotifierMocker struct {
	mock.Mock
}

func (m *SummarizerInitiatorMocker) MockLaunch(ctx context.Context, date time.Time, err error, times int) {
	m.On(launchMethod, ctx, date).Return(err).Times(times)
}

func (m *SummarizerInitiatorMocker) Launch(ctx context.Context, date time.Time) error {
	args := m.Called(ctx, date)
	return args.Error(0)
}

func (m *RepositoryMocker) MockGetByByPeriod(
	ctx context.Context,
	listID string,
	accounts []accounts.Account,
	err error,
	times int,
) {
	m.On(getByByPeriodMethod, ctx, listID).Return(accounts, err).Times(times)
}

func (m *RepositoryMocker) GetByByListID(
	ctx context.Context,
	listID string,
) ([]accounts.Account, error) {
	args := m.Called(ctx, listID)
	return args.Get(0).([]accounts.Account), args.Error(1)
}

func (m *RepositoryMocker) MockUpdate(ctx context.Context, listID string, accountList []accounts.Account, err error, times int) {
	m.On(updateMethod, ctx, listID, accountList).Return(err).Times(times)
}

func (m *RepositoryMocker) Update(ctx context.Context, listID string, accounts []accounts.Account) error {
	args := m.Called(ctx, listID, accounts)
	return args.Error(0)
}

func (m *SummarizerInitiatorNotifierMocker) MockNotify(ctx context.Context, notification interface{}, err error, times int) {
	m.On(notifierMethod, ctx, notification).Return(err).Times(times)
}

func (m *SummarizerInitiatorNotifierMocker) Notify(ctx context.Context, notification interface{}) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}
