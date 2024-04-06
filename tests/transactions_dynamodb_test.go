package tests_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/BradMichel/stori/data"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"

	"github.com/BradMichel/stori/infraestructure"
	"github.com/BradMichel/stori/internal/transactions"
)

func TestTransactionsDynamoDB_Save(t *testing.T) {
	ctx := context.Background()
	tableName := fmt.Sprintf("test_transaction_table_%f", +rand.Float32())
	awsConfig := GetAWSConfig()
	mySession := session.Must(session.NewSession(&awsConfig))
	dynamoDB := dynamodb.New(mySession)
	err := infraestructure.CreateAccountTransactionsDynamoDB(mySession, tableName)
	if err != nil {
		t.Fatal(err)
	}

	dbConfig := transactions.DynamoRepConfig{TableName: tableName}
	dbDep := transactions.NewDynamoDep(dynamoDB)
	repo := transactions.NewDynamoRepository(dbConfig, dbDep)
	accountID := "test_account_id"
	transactionList := data.GetTransactions(accountID)
	err = repo.Save(ctx, accountID, transactionList)
	if err != nil {
		t.Fatal(err)
	}

	output, err := dynamoDB.QueryWithContext(ctx, &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"account_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(accountID),
					},
				},
			},
		},
	})

	var transactionListSaved []transactions.Transaction
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &transactionListSaved)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, transactionList, transactionListSaved)
}
