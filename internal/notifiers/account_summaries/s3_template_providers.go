package account_summaries

import (
	"context"

	"github.com/BradMichel/stori/pkg/aws/s3/downloader"
)

func NewBlueAccountEmailTemplate(
	ctx context.Context,
	client *downloader.Downloader,
	config EmailBuilderConfig,
) (BlueAccountTemplate, error) {
	key := config.BlueAccountTemplatePath
	data, err := client.DownloadWithContext(ctx, key)
	if err != nil {
		return "", err
	}

	return BlueAccountTemplate(data), nil
}

func NewBlackCardEmailTemplate(
	ctx context.Context,
	client *downloader.Downloader,
	config EmailBuilderConfig,
) (BlackCardTemplate, error) {
	key := config.BlackCardTemplatePath
	data, err := client.DownloadWithContext(ctx, key)
	if err != nil {
		return "", err
	}

	return BlackCardTemplate(data), nil
}

func NewGreenCardEmailTemplate(
	ctx context.Context,
	client *downloader.Downloader,
	config EmailBuilderConfig,
) (GreenCardTemplate, error) {
	key := config.GreenCardTemplatePath
	data, err := client.DownloadWithContext(ctx, key)
	if err != nil {
		return "", err
	}

	return GreenCardTemplate(data), nil
}
