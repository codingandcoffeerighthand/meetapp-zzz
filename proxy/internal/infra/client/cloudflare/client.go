package infras_cloudflare

import (
	"context"
	"net/http"
	"proxy-srv/internal/configs"
	"proxy-srv/pkg/gencode/cloudflare_client"
)

func NewClient(cfg configs.CloudflareConfig) (cloudflare_client.ClientWithResponsesInterface, error) {
	return cloudflare_client.NewClientWithResponses(
		cfg.BaseURL,
		cloudflare_client.WithRequestEditorFn(
			func(ctx context.Context, req *http.Request) error {
				req.Header.Set("Authorization", "Bearer "+cfg.AppSecret)
				return nil
			},
		),
	)
}
