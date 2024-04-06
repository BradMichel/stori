package transactions

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/BradMichel/stori/pkg/aws/s3/downloader"
)

const TransactionFileName = "%s_txns.csv"

var ErrInvalidCSV = errors.New("invalid csv")

type S3Repository struct {
	client *downloader.Downloader
}

func NewS3Repository(client *downloader.Downloader) *S3Repository {
	return &S3Repository{
		client: client,
	}
}

func (r *S3Repository) GetByAccountID(ctx context.Context, accountID string) ([]Transaction, error) {
	key := fmt.Sprintf(TransactionFileName, accountID)
	data, err := r.client.DownloadWithContext(ctx, key)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(strings.NewReader(string(data)))
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return r.readLines(accountID, lines)
}

func (r *S3Repository) readLines(accountID string, lines [][]string) ([]Transaction, error) {
	var transactions []Transaction
	for pos, line := range lines {
		if len(line) < 3 {
			return nil, fmt.Errorf("%w: there are less than 3 on: %d columns %v", ErrInvalidCSV, pos, line)
		}

		amount, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %w on: %d columns %v", ErrInvalidCSV, err, pos, line)
		}

		transaction := Transaction{
			ID:        line[0],
			AccountID: accountID,
			Date:      line[1],
			Amount:    amount,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
