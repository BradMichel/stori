package accounts_test

import (
	"context"
	"errors"
	"github.com/BradMichel/stori/internal/mockers"
	"testing"

	"github.com/BradMichel/stori/internal/accounts"
	"github.com/BradMichel/stori/pkg/time"
)

func TestSummarizerInitiator_Launch(t *testing.T) {
	t.Parallel()

	type Fields struct {
		repository *mockers.RepositoryMocker
		notifier   *mockers.SummarizerInitiatorNotifierMocker
	}

	type Args struct {
		ctx              context.Context
		date             time.Time
		accounts         [][]accounts.Account
		listID           []string
		getByByPeriodErr error
		notifyErr        [][]error
		updateErr        []error
		expectedErr      []error
	}

	getFields := func() Fields {
		return Fields{
			repository: new(mockers.RepositoryMocker),
			notifier:   new(mockers.SummarizerInitiatorNotifierMocker),
		}
	}

	notifyErr := errors.New("notify error")
	updateErr := errors.New("update error")
	today := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)
	Mock := func(f Fields, a Args) {
		for i, accountList := range a.accounts {
			f.repository.MockGetByByPeriod(
				a.ctx,
				a.listID[i],
				append([]accounts.Account{}, accountList...),
				a.getByByPeriodErr,
				1,
			)
			for j, account := range accountList {
				f.notifier.MockNotify(a.ctx, account, a.notifyErr[i][j], 1)
				if a.notifyErr[i][j] != nil {
					continue
				}
				accountList[j].Period = a.date
			}

			if len(accountList) == 0 {
				continue
			}

			f.repository.MockUpdate(a.ctx, a.listID[i], accountList, a.updateErr[i], 1)
		}
	}

	tests := []struct {
		name   string
		fields Fields
		args   Args
		mock   func(Fields, Args)
	}{
		{
			name:   "error: notify with some errors",
			fields: getFields(),
			args: Args{
				ctx:  context.Background(),
				date: today,
				accounts: [][]accounts.Account{
					{
						{
							ID:     "id-1",
							Email:  "email-1",
							Period: yesterday,
						},
						{
							ID:     "id-2",
							Email:  "email-2",
							Period: yesterday,
						},
					},
					{
						{
							ID:     "id-3",
							Email:  "email-3",
							Period: yesterday,
						},
						{
							ID:     "id-4",
							Email:  "email-4",
							Period: yesterday,
						},
					},
					{},
				},
				listID:           []string{"1", "2", "3"},
				getByByPeriodErr: nil,
				notifyErr:        [][]error{{nil, nil}, {notifyErr, nil}},
				updateErr:        []error{nil, updateErr, nil},
				expectedErr:      []error{notifyErr, updateErr},
			},
			mock: Mock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := tt.fields
			args := tt.args
			mock := tt.mock
			mock(fields, args)
			s := accounts.NewSummarizerInitiator(fields.repository, fields.notifier)
			err := s.Launch(args.ctx, args.date)
			for _, expectedErr := range args.expectedErr {
				if !errors.Is(err, expectedErr) {
					t.Errorf("expected error was: %v but received: %v", expectedErr, err)
				}
			}
			fields.repository.AssertExpectations(t)
			fields.notifier.AssertExpectations(t)
		})
	}
}
