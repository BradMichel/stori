package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	v1 "github.com/BradMichel/stori/cmd/serverless/account_summarizers/initiator/v1"
)

func main() {
	h, err := v1.Initialize()
	if err != nil {
		panic(err.Error())
	}
	lambda.Start(h.Handle)
}
