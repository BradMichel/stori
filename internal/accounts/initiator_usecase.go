package accounts

import (
	"context"
	"fmt"
	"github.com/BradMichel/stori/pkg/time"
)

type Repository interface {
	GetByByListID(ctx context.Context, listID string) ([]Account, error)
	Update(ctx context.Context, listID string, accounts []Account) error
}

type Notifier interface {
	Notify(ctx context.Context, account interface{}) error
}

type SummarizerInitiator struct {
	repository Repository
	notifier   Notifier
}

func NewSummarizerInitiator(repository Repository, notifier Notifier) *SummarizerInitiator {
	return &SummarizerInitiator{
		repository: repository,
		notifier:   notifier,
	}
}

func (s *SummarizerInitiator) Launch(ctx context.Context, date time.Time) error {
	var err error
	listID := 1
	for {
		var accounts []Account
		var getErr error
		accounts, getErr = s.repository.GetByByListID(ctx, fmt.Sprint(listID))
		if getErr != nil {
			return getErr
		}

		if len(accounts) == 0 {
			break
		}

		for pos, account := range accounts {
			if notifyErr := s.notifier.Notify(ctx, account); notifyErr != nil {
				err = fmt.Errorf("failed to notify account %s: %w; %w", account.ID, notifyErr, err)
				continue
			}

			accounts[pos].Period = date
		}

		updateErr := s.repository.Update(ctx, fmt.Sprint(listID), accounts)
		listID++

		if updateErr != nil {
			err = fmt.Errorf("failed to update account %v: %w; %w", listID, updateErr, err)
			continue
		}
	}

	return err
}
