package notifiers

// AccountSummary data of an account movements summary email notification
type AccountSummary struct {
	ID                  string
	TotalBalance        float64
	AverageDebitAmount  float64
	AverageCreditAmount float64
	Months              []AccountMonthSummary
}

// AccountMonthSummary data of a month account movements summary email notification
type AccountMonthSummary struct {
	Month         string
	NTransactions int
}
