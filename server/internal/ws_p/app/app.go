package ws_app

import (
	"fmt"
	ws_domain "proxy-srv/internal/ws_p/domain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ProxyGrpcClient interface {
	RenegotiateSession(sessionId string, sdpAnswer string) error
	RoomPull(roomId string) error
}
type app struct {
	ws_domain.App
	gcl ProxyGrpcClient
	log *zap.Logger
}

func NewApp(gcl ProxyGrpcClient, log *zap.Logger) *app {
	return &app{
		App: ws_domain.App{
			Rooms: make(map[string]*ws_domain.Room),
			Users: make(map[string]*ws_domain.Participant),
		},
		gcl: gcl,
		log: log,
	}
}

type WsConnectReq struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
	Nonce     string `json:"nonce"`
	Conn      *websocket.Conn
}

func (a *app) WsConnection(req *WsConnectReq) error {
	if !common.IsHexAddress(req.Address) {
		err := fmt.Errorf(
			"invalid address: %s",
			req.Address,
		)
		a.log.Error(err.Error())
		return err
	}
	address := common.HexToAddress(req.Address)
	sig := common.FromHex(req.Signature)
	if len(sig) != 65 {
		e := fmt.Errorf(
			"invalid signature: %s",
			req.Signature,
		)
		a.log.Error(e.Error())
		return e
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	nonce := req.Nonce
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(nonce), nonce)
	hash := crypto.Keccak256Hash([]byte(msg))
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		a.log.Error(err.Error())
		return err
	}
	recoverAddress := crypto.PubkeyToAddress(*pubKey)
	if address != recoverAddress {
		e := fmt.Errorf(
			"invalid signature: %s",
			req.Signature,
		)
		a.log.Error(e.Error())
		return e
	}
	a.App.AddUser(&ws_domain.Participant{
		Address: req.Address,
		Conn:    req.Conn,
	})
	return nil
}

func (a *app) CreateRoom(roomId string) error {
	a.App.AddRoom(roomId)
	return nil
}
