package accounts

import "github.com/BradMichel/stori/pkg/time"

type List struct {
	ID       string    `json:"id"`
	Accounts []Account `json:"accounts"`
	Period   time.Time `json:"last_update"`
}
