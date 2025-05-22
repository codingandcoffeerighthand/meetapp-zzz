package infra_cloudflare

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"proxy-srv/internal/v2/configs"
	domain "proxy-srv/internal/v2/domain"
	"proxy-srv/pkg/gencode/cloudflare_client"
	"sync"
	"time"
)

type cloudFlareInfa struct {
	cloudflare_client.ClientWithResponsesInterface
	appId string
	sync.Mutex
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
		Mutex:                        sync.Mutex{},
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

func (c *cloudFlareInfa) AddLocalTrack(
	ctx context.Context,
	sessionId string, sdpOffer string,
	domainTracks []domain.Track) (
	string, error,
) {
	tracks := make([]cloudflare_client.TrackObject, len(domainTracks))
	localVal := cloudflare_client.TrackObjectLocationLocal
	for i := range tracks {
		tracks[i].Location = &localVal
		tracks[i].Mid = &domainTracks[i].Mid
		tracks[i].TrackName = &domainTracks[i].TrackName
		tracks[i].SessionId = &domainTracks[i].SessionID
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
		return "", errors.New("error passing JSON response for add local track")
	}
	if resp.JSON200 == nil {
		return "", fmt.Errorf("cloudflare error %v", resp.HTTPResponse)
	}
	if resp.JSON200.SessionDescription == nil ||
		resp.JSON200.SessionDescription.Sdp == nil {
		return "", errors.New("error passing JSON response for add local track")
	}
	return *resp.JSON200.SessionDescription.Sdp, err
}

var muAddRemoteTrack sync.Mutex

func (c *cloudFlareInfa) AddRemoteTrack(
	ctx context.Context, sessionId string,
	domainTracks []domain.Track) (
	string, error,
) {
	muAddRemoteTrack.Lock()
	time.Sleep(100 * time.Millisecond)
	defer muAddRemoteTrack.Unlock()
	tracks := make([]cloudflare_client.TrackObject, len(domainTracks))
	remoteVal := cloudflare_client.TrackObjectLocationRemote
	for i := range tracks {
		tracks[i].Location = &remoteVal
		tracks[i].Mid = nil
		tracks[i].TrackName = &domainTracks[i].TrackName
		tracks[i].SessionId = &domainTracks[i].SessionID
	}
	body := cloudflare_client.TracksRequest{
		Tracks: &tracks,
	}
	resp, err := c.PostAppsAppIdSessionsSessionIdTracksNewWithResponse(
		ctx, c.appId, sessionId, body,
	)
	if resp == nil {
		return "", err
	}
	if resp.JSON200 == nil {
		return "", fmt.Errorf("error add remote track: %v", resp.HTTPResponse)
	}
	if resp.JSON200.SessionDescription == nil ||
		resp.JSON200.SessionDescription.Sdp == nil {
		return "", errors.New("error passing JSON response for add remote track")
	}
	return *resp.JSON200.SessionDescription.Sdp, err
}

func (c *cloudFlareInfa) RenegatiateSession(
	ctx context.Context,
	session string, sdpAnswer string) (string, error) {
	answerType := cloudflare_client.SessionDescriptionTypeAnswer
	sdpInfo := cloudflare_client.SessionDescription{
		Sdp:  &sdpAnswer,
		Type: &answerType,
	}
	_, err := c.PutAppsAppIdSessionsSessionIdRenegotiateWithResponse(
		ctx, c.appId, session,
		cloudflare_client.PutAppsAppIdSessionsSessionIdRenegotiateJSONRequestBody{
			SessionDescription: &sdpInfo,
		},
	)
	return "", err
}

func (c *cloudFlareInfa) GetStatusSession(sessionId string) (
	[]domain.Track, error) {
	resp, err := c.GetAppsAppIdSessionsSessionIdWithResponse(context.Background(), c.appId, sessionId)
	if resp == nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("error get status session: %v", resp.HTTPResponse)
	}
	if resp.JSON200.Tracks == nil {
		return nil, errors.New("error get status session")
	}
	tracks := make([]domain.Track, len(*resp.JSON200.Tracks))
	for i := range tracks {
		tracks[i].Mid = *(*resp.JSON200.Tracks)[i].Mid
		tracks[i].TrackName = *(*resp.JSON200.Tracks)[i].TrackName
		tracks[i].SessionID = *(*resp.JSON200.Tracks)[i].SessionId
		tracks[i].Location = string(*(*resp.JSON200.Tracks)[i].Location)
	}
	return tracks, err
}
