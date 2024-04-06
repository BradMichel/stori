package tests_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/BradMichel/stori/infraestructure"
	"github.com/BradMichel/stori/internal/accounts"
)

func TestAccountsDynamoDB(t *testing.T) {
	ctx := context.Background()
	tableName := "test_account_table"
	awsConfig := GetAWSConfig()
	mySession := session.Must(session.NewSession(&awsConfig))
	dynamoDB := dynamodb.New(mySession)
	err := infraestructure.CreateAccountDynamoDB(mySession, tableName)
	if err != nil {
		t.Fatal(err)
	}

	dbConfig := accounts.RepConfig{TableName: tableName}
	dbDep := accounts.NewDynamoDep(dynamoDB)
	repo := accounts.NewDynamoRepository(dbConfig, dbDep)
	today := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)
	accountsList := []accounts.List{
		{
			ID: fmt.Sprint(1),
			Accounts: []accounts.Account{
				{
					ID:     "test_id_1",
					Email:  "test_email_1",
					Period: yesterday,
				},
				{
					ID:     "test_id_2",
					Email:  "test_email_2",
					Period: yesterday,
				},
			},
			Period: yesterday,
		},
		{
			ID: fmt.Sprint(2),
			Accounts: []accounts.Account{
				{
					ID:     "test_id_3",
					Email:  "test_email_3",
					Period: yesterday,
				},
				{
					ID:     "test_id_4",
					Email:  "test_email_4",
					Period: yesterday,
				},
				{
					ID:     "test_id_5",
					Email:  "test_email_5",
					Period: yesterday,
				},
			},
			Period: yesterday,
		},
	}

	for _, accountList := range accountsList {
		err := repo.Update(ctx, accountList.ID, accountList.Accounts)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, accountListSaved := range accountsList {
		accountsList, err := repo.GetByByListID(ctx, accountListSaved.ID)
		if err != nil {
			t.Fatal(err)
		}

		if len(accountsList) != len(accountListSaved.Accounts) {
			t.Fatalf("expected %d accounts, got %d", len(accountListSaved.Accounts), len(accountsList))
		}

		for i, account := range accountsList {
			if account.ID != accountListSaved.Accounts[i].ID {
				t.Fatalf("expected account %s, got %s", accountListSaved.Accounts[i].ID, account.ID)
			}
		}
	}
}
