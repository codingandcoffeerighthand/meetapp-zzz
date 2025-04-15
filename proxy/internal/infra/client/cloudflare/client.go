package infras_cloudflare

import (
	"context"
	"net/http"
	"proxy-srv/gencode/cloudflare_client"
	"proxy-srv/internal/configs"
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
