package initial

import (
	"log"
	"proxy-srv/internal/utils"
	"proxy-srv/internal/v2/app"
	"proxy-srv/internal/v2/configs"
	infra_cloudflare "proxy-srv/internal/v2/infra/cloudflare"
	infra_smc "proxy-srv/internal/v2/infra/meet_smc"
)

func NewApp(cfg *configs.Config) (*app.App, func(), error) {
	clf, err := infra_cloudflare.NewClient(*cfg)
	if err != nil {
		log.Fatalf("error new cloudflare client: %v", err)
		return nil, nil, err
	}
	meet, err := infra_smc.NewSMCInfra(cfg)
	if err != nil {
		log.Fatalf("error new meet client: %v", err)
		return nil, nil, err
	}
	logger, cl, err := utils.NewLogger()
	if err != nil {
		log.Fatalf("error new logger: %v", err)
		return nil, nil, err
	}
	cleanUp := func() {
		cl()
	}
	return app.NewApp(clf, meet, logger), cleanUp, nil
}
