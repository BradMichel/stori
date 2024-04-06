package summarizers

// AccountMonthSummary data of a month account movements summary
type AccountMonthSummary struct {
	Date         string  `json:"date"`
	Debit        float64 `json:"debit"`
	DebitCount   int     `json:"debit_count"`
	Credit       float64 `json:"credit"`
	CreditCount  int     `json:"credit_count"`
	Transactions int     `json:"transactions"`
}

// AccountSummary data of an account movements summary
type AccountSummary struct {
	AccountID   string                `json:"account_id"`
	Months      []AccountMonthSummary `json:"months"`
	CreditAvg   float64               `json:"credit_avg"`
	DebitAvg    float64               `json:"debit_avg"`
	Total       float64               `json:"total"`
	totalDebit  float64
	countDebit  int
	totalCredit float64
	countCredit int
}
