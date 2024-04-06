package accounts

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kelseyhightower/envconfig"
)

const (
	idKey = "id"
)

type DynamoRepository struct {
	config RepConfig
	dep    RepDep
}

type RepConfig struct {
	TableName string `envconfig:"ACCOUNT_DYNAMODB_TABLE" default:"local-stori-accounts" required:"true"`
}

type RepDep struct {
	dynamoDB *dynamodb.DynamoDB
}

func NewDynamoRepository(config RepConfig, dep RepDep) *DynamoRepository {
	return &DynamoRepository{
		config: config,
		dep:    dep,
	}
}

func NewDynamoConfig() (RepConfig, error) {
	c := RepConfig{}
	err := envconfig.Process("", &c)
	return c, err
}

func NewDynamoDep(dynamoDB *dynamodb.DynamoDB) RepDep {
	return RepDep{
		dynamoDB: dynamoDB,
	}
}

func (r *DynamoRepository) GetByByListID(ctx context.Context, listID string) ([]Account, error) {
	input := r.getByByPeriod(listID)
	result, err := r.dep.dynamoDB.QueryWithContext(ctx, &input)
	if err != nil {
		return nil, err
	}

	var accountList []List
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &accountList)
	if err != nil {
		return nil, err
	}

	if len(accountList) == 0 {
		return nil, nil

	}

	return accountList[0].Accounts, nil
}

func (r *DynamoRepository) Update(ctx context.Context, listID string, accounts []Account) error {
	if len(accounts) == 0 {
		return nil
	}

	input := r.update(listID, accounts)
	_, err := r.dep.dynamoDB.PutItemWithContext(ctx, &input)
	return err
}

func (r *DynamoRepository) getByByPeriod(listID string) dynamodb.QueryInput {
	input := dynamodb.QueryInput{
		TableName: aws.String(r.config.TableName),
		KeyConditions: map[string]*dynamodb.Condition{
			idKey: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(listID),
					},
				},
			},
		},
	}

	return input
}

func (r *DynamoRepository) update(listID string, accounts []Account) dynamodb.PutItemInput {
	list := List{
		ID:       listID,
		Accounts: accounts,
		Period:   accounts[0].Period,
	}

	item, _ := dynamodbattribute.MarshalMap(list)
	return dynamodb.PutItemInput{
		TableName: aws.String(r.config.TableName),
		Item:      item,
	}
}
