package infras_cloudflare

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"proxy-srv/internal/proxy/configs"
	"proxy-srv/pkg/gencode/cloudflare_client"
)

type cloudFlareInfa struct {
	cloudflare_client.ClientWithResponsesInterface
	appId string
}

func NewClient(cfg configs.Config) (*cloudFlareInfa, error) {
	cl, err := cloudflare_client.NewClientWithResponses(
		cfg.CloudflareConfig.BaseURL,
		cloudflare_client.WithRequestEditorFn(
			func(ctx context.Context, req *http.Request) error {
				req.Header.Set("Authorization", "Bearer "+cfg.CloudflareConfig.AppSecret)
				return nil
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return &cloudFlareInfa{
		ClientWithResponsesInterface: cl,
		appId:                        cfg.CloudflareConfig.AppId,
	}, nil
}

func (c *cloudFlareInfa) NewSession(ctx context.Context) (string, error) {
	resp, err := c.PostAppsAppIdSessionsNewWithResponse(ctx, c.appId)
	if err != nil {
		return "", err
	}
	if resp.JSON201 == nil {
		return "", fmt.Errorf("error new session: %s", resp.Body)
	}
	if resp.JSON201.SessionId != nil {
		return *resp.JSON201.SessionId, nil
	}
	return "", errors.New("error parsing JSON response for new session")
}

func (c *cloudFlareInfa) AddLocalTrack(ctx context.Context, sessionId string, sdpOffer string, tracks []cloudflare_client.TrackObject) (
	cloudflare_client.TracksResponse, error,
) {
	localVal := cloudflare_client.TrackObjectLocationLocal
	for i := range tracks {
		tracks[i].Location = &localVal
	}
	sdpType := cloudflare_client.SessionDescriptionType("offer")
	body := cloudflare_client.TracksRequest{
		SessionDescription: &cloudflare_client.SessionDescription{
			Type: &sdpType,
			Sdp:  &sdpOffer,
		},
		Tracks: &tracks,
	}
	resp, err := c.PostAppsAppIdSessionsSessionIdTracksNewWithResponse(
		ctx, c.appId, sessionId, body,
	)
	if resp == nil {
		return cloudflare_client.TracksResponse{}, errors.New("error passing JSON response for add local track")
	}
	return *resp.JSON200, err
}
func (c *cloudFlareInfa) AddRemoteTrack(ctx context.Context, sessionId string, tracks []cloudflare_client.TrackObject) (
	cloudflare_client.TracksResponse, error,
) {
	remoteVal := cloudflare_client.TrackObjectLocationRemote
	for i := range tracks {
		tracks[i].Location = &remoteVal
		tracks[i].Mid = nil
	}
	body := cloudflare_client.TracksRequest{
		Tracks: &tracks,
	}
	resp, err := c.PostAppsAppIdSessionsSessionIdTracksNewWithResponse(
		ctx, c.appId, sessionId, body,
	)
	if resp == nil {
		return cloudflare_client.TracksResponse{}, errors.New("error passing JSON response for add remote track")
	}
	if resp.JSON200 == nil {
		return cloudflare_client.TracksResponse{}, errors.New("error passing JSON response for add remote track")
	}
	return *resp.JSON200, err
}

func (c *cloudFlareInfa) RenegatiateSession(ctx context.Context, session string, sdpAnswer string) (cloudflare_client.SessionDescription, error) {
	answerType := cloudflare_client.SessionDescriptionTypeAnswer
	sdpInfo := cloudflare_client.SessionDescription{
		Sdp:  &sdpAnswer,
		Type: &answerType,
	}
	resp, err := c.PutAppsAppIdSessionsSessionIdRenegotiateWithResponse(
		ctx, c.appId, session,
		cloudflare_client.PutAppsAppIdSessionsSessionIdRenegotiateJSONRequestBody{
			SessionDescription: &sdpInfo,
		},
	)
	if resp == nil {
		return cloudflare_client.SessionDescription{}, errors.New("error passing JSON response for renegotiate session")
	}
	return *resp.JSON200, err
}

func (c *cloudFlareInfa) GetStatusSession(sessionId string) (
	*cloudflare_client.GetSessionStateResponse, error) {
	resp, err := c.GetAppsAppIdSessionsSessionIdWithResponse(context.Background(), c.appId, sessionId)
	if resp == nil {
		return nil, errors.New("error passing JSON response for get status session")
	}
	if resp.JSON200 == nil {
		return nil, errors.New("error passing JSON response for get status session")
	}
	return resp.JSON200, err
}
