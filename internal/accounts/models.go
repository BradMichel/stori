package accounts

import "github.com/BradMichel/stori/pkg/time"

const (
	BlueAccountKey = "blue_account"
	BlackCardKey   = "black_card"
	GreenCardKey   = "green_card"
)

type Account struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	Period time.Time
}

type Transaction struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Transaction float64   `json:"transaction"`
}
