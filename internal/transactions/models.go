package transactions

type Transaction struct {
	ID        string  `json:"id"`
	AccountID string  `json:"account_id"`
	Date      string  `json:"date"`
	Amount    float64 `json:"amount"`
}
