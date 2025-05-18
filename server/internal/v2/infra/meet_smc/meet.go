package infra_smc

import (
	"context"
	"fmt"
	"proxy-srv/internal/v2/configs"
	"proxy-srv/internal/v2/domain"
	"proxy-srv/pkg/gencode/smc_gen/meet_smc"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type smcInfra struct {
	*bind.BoundContract
	auth               *bind.TransactOpts
	contract           *meet_smc.Meet
	addr               *common.Address
	conn               bind.DeployBackend
	frontendEventMutex sync.Mutex
	errChan            chan error
}

func NewSMCInfra(cfg *configs.Config) (*smcInfra, error) {
	rpc := cfg.Web3Config.RpcUrl
	privateKeyStr := cfg.Web3Config.PrivateKey
	contractAddress := cfg.Web3Config.ContractAddress
	conn, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	chan_id, err := conn.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	contract := meet_smc.NewMeet()
	addr := common.HexToAddress(contractAddress)
	instance := contract.Instance(conn, addr)
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	auth := bind.NewKeyedTransactor(privateKey, chan_id)

	return &smcInfra{
		BoundContract:      instance,
		auth:               auth,
		contract:           contract,
		addr:               &addr,
		conn:               conn,
		frontendEventMutex: sync.Mutex{},
		errChan:            make(chan error, 20),
	}, nil
}

func (s *smcInfra) SetNewSession(
	ctx context.Context, roomID string,
	oldSessionID string, sessionID string) error {
	tx, err := bind.Transact(s.BoundContract, s.auth,
		s.contract.PackNewSession(roomID, oldSessionID, sessionID))
	if err != nil {
		return err
	}
	_, err = bind.WaitMined(ctx, s.conn, tx.Hash())
	return err
}

func (s *smcInfra) SetIceServers(ctx context.Context, iceServers string) error {
	tx, err := bind.Transact(s.BoundContract, s.auth,
		s.contract.PackSetIceServers(iceServers))
	if err != nil {
		return err
	}
	_, err = bind.WaitMined(ctx, s.conn, tx.Hash())
	return err
}

func (s *smcInfra) EmitFrontEndEvent(
	ctx context.Context,
	roomID string, sessionID string,
	eventType string, data []byte) error {
	s.frontendEventMutex.Lock()
	defer s.frontendEventMutex.Unlock()
	tx, err := bind.Transact(s.BoundContract, s.auth,
		s.contract.PackEmitEventToFrontend(roomID, sessionID, eventType, data))
	if err != nil {
		return err
	}
	_, err = bind.WaitMined(ctx, s.conn, tx.Hash())
	return err
}

func (s *smcInfra) GetRoomInfo(ctx context.Context, roomID string) (
	domain.Room, error) {
	roomInfo, err := bind.Call(s.BoundContract, &bind.CallOpts{Pending: true, From: s.auth.From},
		s.contract.PackGetRoomInfo(roomID), s.contract.UnpackGetRoomInfo)
	if err != nil {
		return domain.Room{}, err
	}
	var rs = domain.Room{
		RoomID:       roomID,
		Name:         roomInfo.Name,
		Creator:      roomInfo.Creator.Hex(),
		Participants: make([]domain.Participant, len(roomInfo.Participants)),
	}
	for i, participant := range roomInfo.Participants {
		rs.Participants[i] = domain.Participant{
			WalletAddress: participant.WalletAddress.Hex(),
			Name:          participant.Name,
			SessionID:     participant.SessionID,
			Tracks:        make([]domain.Track, len(roomInfo.Participants[i].Tracks)),
		}
		for j, track := range roomInfo.Participants[i].Tracks {
			rs.Participants[i].Tracks[j] = domain.Track{
				Mid:       track.Mid,
				TrackName: track.TrackName,
				SessionID: track.SessionId,
				Location:  track.Location,
			}
		}
	}
	return rs, nil
}

func (s *smcInfra) GetIceServers(ctx context.Context) (string, error) {
	iceServers, err := bind.Call(s.BoundContract, &bind.CallOpts{Pending: true, From: s.auth.From},
		s.contract.PackGetIceServers(), s.contract.UnpackGetIceServers)
	if err != nil {
		return "", err
	}
	return iceServers, nil
}

func (s *smcInfra) GetParticipantInfoBySessionId(ctx context.Context, roomID string,
	sessionID string) (*domain.Participant, error) {
	participant, err := bind.Call(s.BoundContract, &bind.CallOpts{Pending: true, From: s.auth.From},
		s.contract.PackGetParticipantInfoBySessionId(roomID, sessionID), s.contract.UnpackGetParticipantInfoBySessionId)
	if err != nil {
		return nil, err
	}
	rs := domain.Participant{
		WalletAddress: participant.WalletAddress.Hex(),
		Name:          participant.Name,
		SessionID:     participant.SessionID,
		Tracks:        make([]domain.Track, len(participant.Tracks)),
	}
	for i, track := range participant.Tracks {
		rs.Tracks[i] = domain.Track{
			Mid:       track.Mid,
			TrackName: track.TrackName,
			SessionID: track.SessionId,
			Location:  track.Location,
		}
	}
	return &rs, nil
}

func (s *smcInfra) ErrChan() <-chan error {
	return s.errChan
}
