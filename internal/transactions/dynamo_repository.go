package transactions

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kelseyhightower/envconfig"
)

type DynamoRepository struct {
	config DynamoRepConfig
	dep    DynamoRepDep
}

type DynamoRepConfig struct {
	TableName string `envconfig:"ACCOUNT_TRANSACTIONS_DYNAMODB_TABLE" default:"local-stori-transactions" required:"true"`
}

type DynamoRepDep struct {
	dynamoDB *dynamodb.DynamoDB
}

func NewDynamoRepository(config DynamoRepConfig, dep DynamoRepDep) *DynamoRepository {
	return &DynamoRepository{
		config: config,
		dep:    dep,
	}
}

func NewDynamoConfig() (DynamoRepConfig, error) {
	c := DynamoRepConfig{}
	err := envconfig.Process("", &c)
	return c, err
}

func NewDynamoDep(dynamoDB *dynamodb.DynamoDB) DynamoRepDep {
	return DynamoRepDep{
		dynamoDB: dynamoDB,
	}
}

func (r *DynamoRepository) Save(ctx context.Context, accountID string, transactions []Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	writeRequests := make([]*dynamodb.WriteRequest, len(transactions))
	for pos, transaction := range transactions {
		transaction.AccountID = accountID
		item, err := dynamodbattribute.MarshalMap(transaction)
		if err != nil {
			return err
		}

		writeRequests[pos] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: item,
			},
		}
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			r.config.TableName: writeRequests,
		},
	}

	_, err := r.dep.dynamoDB.BatchWriteItemWithContext(ctx, input)
	return err
}
