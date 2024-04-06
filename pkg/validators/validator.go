package validators

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

var Err = fmt.Errorf("error validating")

type Config struct {
	Source string `envconfig:"VALIDATOR_SOURCE" default:"request.json" required:"true"`
}

type Service struct {
	conf   Config
	loader gojsonschema.JSONLoader
}

func New(conf Config, e embed.FS) (*Service, error) {
	schema, err := e.ReadFile(conf.Source)
	if err != nil {
		return nil, err

	}

	loader := gojsonschema.NewBytesLoader(schema)
	_, err = loader.LoadJSON()
	if err != nil {
		return nil, err
	}

	return &Service{
		conf:   conf,
		loader: loader,
	}, nil
}

func (s *Service) Validate(_ context.Context, body []byte) error {
	schemaLoader := gojsonschema.NewBytesLoader(body)
	result, err := gojsonschema.Validate(s.loader, schemaLoader)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			return s.getSyntaxError(err, jsonErr, body)
		}
		err = fmt.Errorf("%w: body %v", err, body)
		return err
	}

	if !result.Valid() {
		invalidErr := ""
		for _, resultError := range result.Errors() {
			invalidErr = fmt.Sprintf(
				"%s\n%s",
				invalidErr,
				resultError.String(),
			)
		}

		return fmt.Errorf("%w body: %s details: %s", Err, body, invalidErr)
	}

	return nil
}

func (s *Service) getSyntaxError(err error, jsonErr *json.SyntaxError, body []byte) error {
	problematicPart := ""
	if int(jsonErr.Offset) < len(body) {
		start := int(jsonErr.Offset) - 2
		if start < 0 {
			start = 0
		}
		end := int(jsonErr.Offset) + 3
		if end > len(body) {
			end = len(body)
		}
		problematicPart = string(body[start:end])
	}
	return fmt.Errorf(
		"%w: JSON syntax error at byte %v: %w. Problematic part: %s with value %s",
		Err, jsonErr.Offset, err, problematicPart, body,
	)
}
