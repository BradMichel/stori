package account_summaries

type Email struct {
	Headers   EmailHeaders
	Data      AccountSummary
	Resources AccountSummaryEmailResources
}

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

type EmailHeaders struct {
	Title           string `json:"title" default:"Account Summary"`
	CardHeader      string `json:"card_header" default:"Stori Cuenta: Account Summary for"`
	TotalBalance    string `json:"total_balance" default:"Total Balance:"`
	AvgDebitAmount  string `json:"avg_debit" default:"Average Debit Amount:"`
	AvgCreditAmount string `json:"avg_credit" default:"Average Credit Amount:"`
	MonthlySummary  string `json:"monthly_summary" default:"Monthly Summary"`
	Month           string `json:"month" default:"Month"`
	NTransactions   string `json:"n_transactions" default:"Number of Transactions"`
}

type AccountSummaryEmailResources struct {
	FootImage string
}

type FootImageByAccountType struct {
	BlueAccount string `envconfig:"BLUE_ACCOUNT_IMAGE_URL" default:"https://www.storicard.com/_next/static/media/deposits_products_carousel.39672734.webp"`
	BlackCard   string `envconfig:"BLACK_CARD_IMAGE_URL" default:"https://www.storicard.com/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Fstori_black_beta_2.3fa8492e.webp&w=1920&q=100"`
	GreenCard   string `envconfig:"GREEN_CARD_IMAGE_URL" default:"https://www.storicard.com/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Fstoricard_products_carousel.bf76e094.webp&w=1920&q=75"`
}
