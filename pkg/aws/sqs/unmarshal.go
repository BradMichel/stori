package sqs

import (
	"encoding/json"
	"fmt"
)

var UnmarshalErr = fmt.Errorf("error unmarshalling body")

type BodyResponse struct {
	Message string `json:"Message"`
}

func Unmarshal(body string, v interface{}) (BodyResponse, error) {
	var msg BodyResponse
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return msg, fmt.Errorf("%w: %w", UnmarshalErr, err)

	}

	err = json.Unmarshal([]byte(msg.Message), v)
	if err != nil {
		return msg, fmt.Errorf("%w: %w", UnmarshalErr, err)
	}

	return msg, nil
}
