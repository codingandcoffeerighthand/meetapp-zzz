package ws_domain

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Participant struct {
	Address   string `json:"address"`
	RoomId    string `json:"room_id"`
	SessionId string `json:"session_id"`
	Conn      *websocket.Conn
}

type Room struct {
	RoomId       string                  `json:"room_id"`
	Participants map[string]*Participant `json:"participants"`
}

type App struct {
	sync.Mutex
	Rooms map[string]*Room        `json:"rooms"`
	Users map[string]*Participant `json:"users"`
}

func (a *App) AddRoom(roomId string) {
	a.Lock()
	defer a.Unlock()
	a.Rooms[roomId] = &Room{
		RoomId:       roomId,
		Participants: make(map[string]*Participant),
	}
}

func (a *App) AddUser(user *Participant) {
	a.Lock()
	defer a.Unlock()
	a.Users[user.Address] = user
}

func (a *App) JoinRoom(walletAddress string, roomId string) {
	a.Lock()
	defer a.Unlock()
	a.Rooms[roomId].Participants[walletAddress] = a.Users[walletAddress]
}
