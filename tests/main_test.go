package tests_test

import (
	"context"
	"testing"

	"github.com/BradMichel/stori/tests"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	tests.MainDockerCompose(ctx, m, "../infraestructure/docker-compose.yml")
}
