package accounts

import "github.com/BradMichel/stori/pkg/time"

type Account struct {
	ID     string    `json:"id"`
	Email  string    `json:"email"`
	Period time.Time `json:"last_update"`
}

type Transaction struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Transaction float64   `json:"transaction"`
}
