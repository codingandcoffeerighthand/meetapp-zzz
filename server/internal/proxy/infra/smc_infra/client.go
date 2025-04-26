package smc_infra

import (
	"context"
	"encoding/json"
	"fmt"
	"proxy-srv/internal/proxy/configs"
	"proxy-srv/pkg/gencode/smc_gen"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type smcInfra struct {
	*bind.BoundContract
	auth     *bind.TransactOpts
	contract *smc_gen.Meeeting
	addr     *common.Address
	conn     bind.DeployBackend
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
	contract := smc_gen.NewMeeeting()
	addr := common.HexToAddress(contractAddress)
	instance := contract.Instance(conn, addr)
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	auth := bind.NewKeyedTransactor(privateKey, chan_id)

	return &smcInfra{
		BoundContract: instance,
		auth:          auth,
		contract:      contract,
		addr:          &addr,
		conn:          conn,
	}, nil
}

func (s *smcInfra) CheckAuthorized(addressStr string) (bool, error) {
	address := common.HexToAddress(addressStr)
	callOpts := &bind.CallOpts{
		Pending: true,
		From:    s.auth.From,
	}
	return bind.Call(s.BoundContract, callOpts, s.contract.PackCheckAuthorizedBackend(address), s.contract.UnpackCheckAuthorizedBackend)
}

var block *uint64 = nil

func (s *smcInfra) SetParticipantSessionID(
	ctx context.Context,
	room_id string, addr string, session_id string,
) error {
	address := common.HexToAddress(addr)
	tx, err := bind.Transact(s.BoundContract, s.auth, s.contract.PackSetParticipantSessionID(
		room_id, address, session_id,
	))
	if err != nil {
		return err
	}
	_, err = bind.WaitMined(ctx, s.conn, tx.Hash())
	return err
}

func (s *smcInfra) GetParticipantsAndTracksOfRoom(roomId string) (smc_gen.GetParticipantOfRoomOutput, error) {
	return bind.Call(s.BoundContract,
		&bind.CallOpts{Pending: true, From: s.auth.From},
		s.contract.PackGetParticipantOfRoom(roomId), s.contract.UnpackGetParticipantOfRoom)
}

var mutexSendEventToFrontend = sync.Mutex{}

func (s *smcInfra) EmitEventToFrontend(rooId string, addrStr string, data any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	addr := common.HexToAddress(addrStr)
	// _, err = bind.Call(s.BoundContract, &bind.CallOpts{Pending: true, From: s.auth.From},
	// 	s.contract.PackForwardEventToFrontend(rooId, addr, dataBytes),
	// 	utils.UnpackEmpty)
	mutexSendEventToFrontend.Lock()
	defer mutexSendEventToFrontend.Unlock()
	tx, err := bind.Transact(s.BoundContract, s.auth, s.contract.PackForwardEventToFrontend(rooId, addr, dataBytes))
	if err != nil {
		return err
	}
	_, err = bind.WaitMined(context.Background(), s.conn, tx.Hash())
	return err
}
