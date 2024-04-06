package summarizers

import "github.com/BradMichel/stori/internal/accounts"

type AccountNotification struct {
	Account accounts.Account `json:"account"`
	Summary AccountSummary   `json:"summary"`
}
