package ws_server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"proxy-srv/internal/utils"
	ws_app "proxy-srv/internal/ws_p/app"
	"proxy-srv/internal/ws_p/config"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type App interface {
	WsConnection(req *ws_app.WsConnectReq) error
	ClientRenegotiateSession(req ws_app.ClientRenegotiateSessionReq) error
	PullRoom(roomId string) error
}

type server struct {
	port int
	app  App
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Cẩn thận với setting này trong production
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewServer(cfg config.ServerConfig, app App) *server {
	return &server{
		port: cfg.Port,
		app:  app,
	}
}

func (s *server) Listen() error {
	r := mux.NewRouter()
	// CORS middleware (configure as needed)
	r.HandleFunc("/ws", s.HandleConnections).Methods("GET")

	// Route cho trang chủ (tùy chọn, để kiểm tra server có chạy không)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "HTTP server is running. Connect to /ws for WebSocket.")
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), r)
}

func (s *server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}
	defer ws.Close()
	// login
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err = utils.WithContext(ctx, func() error {
		loginReq := &ws_app.WsConnectReq{}
		loginReq.Conn = ws
		err = ws.ReadJSON(loginReq)
		if err != nil {
			log.Println(err)
			return err
		}
		err = s.app.WsConnection(loginReq)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		ws.WriteJSON(Message{
			MessageName: "error",
			Data:        err.Error(),
		})
		ws.Close()
		log.Println(err)
		return
	}
	ws.WriteJSON(Message{
		MessageName: "success",
	})
	// handle messages
	var msg = Message{}
	for {
		err = ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
		switch msg.MessageName {
		case "client_renegotiate_session":
			req := ws_app.ClientRenegotiateSessionReq{}
			data, ok := msg.Data.(map[string]any)
			if !ok {
				log.Println("invalid message")
				continue
			}
			req.SessionId = data["session_id"].(string)
			req.SdpAnswer = data["sdp_answer"].(string)
			err = s.app.ClientRenegotiateSession(req)
			if err != nil {
				log.Println(err)
			}
		case "pull-room":
			var roomId, ok = msg.Data.(map[string]any)
			if !ok {
				log.Println("invalid message")
				continue
			}
			req := PullRoom{
				RoomId: roomId["room_id"].(string),
			}
			err = s.app.PullRoom(req.RoomId)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Println("unknown message:", msg.MessageName)
		}
	}
}
