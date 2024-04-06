package tests

import (
	"context"
	"log"
	"os"
	"testing"

	dc "github.com/testcontainers/testcontainers-go/modules/compose"
)

func MainDockerCompose(ctx context.Context, m *testing.M, path string) {
	compose, err := dc.NewDockerCompose(path)
	if err != nil {
		log.Fatalf("Could not create docker compose: %v", err)
	}

	if err := compose.Up(ctx, dc.Wait(true)); err != nil {
		log.Fatalf("Could not start docker compose: %v", err)
	}

	defer func() {
		rec := recover()
		if rec != nil {
			log.Fatalf("Panic: %v", rec)
		}
		err = compose.Down(ctx, dc.RemoveOrphans(true), dc.RemoveImagesLocal)
		if err != nil {
			log.Fatalf("Could not stop docker compose: %v", err)
		}
	}()

	code := m.Run()

	os.Exit(code)
}
