package ws_app

import ws_domain "proxy-srv/internal/ws_p/domain"

type JoinRoomReq struct {
	ParticipantAddr string
	RoomId          string
	SessionId       string
	SdpAnswer       string
	SdpOffer        string
}
type JoinRoomRes struct {
	ParticipantAddr string `json:"participant_addr"`
	MessageName     string `json:"message_name"`
	SdpAnswer       string `json:"sdp_answer"`
	Session         string `json:"session"`
}

func (a *app) JoinRoom(req JoinRoomReq) error {
	a.App.JoinRoom(req.ParticipantAddr, req.RoomId)
	defer func() {
		if err := recover(); err != nil {
			a.log.Error(err.(error).Error())
		}
	}()
	// side effect
	go func(p *ws_domain.Participant, sdpAnswer string, sessionId string) {
		p.Conn.WriteJSON(JoinRoomRes{
			ParticipantAddr: p.Address,
			MessageName:     "joined_room",
			SdpAnswer:       sdpAnswer,
			Session:         sessionId,
		})

	}(a.App.Rooms[req.RoomId].Participants[req.ParticipantAddr], req.SdpAnswer, req.SessionId)

	go func(ps map[string]*ws_domain.Participant, pAddr string) {
		for _, p := range ps {
			if p.Address == pAddr {
				continue
			}
			p.Conn.WriteJSON(JoinRoomRes{
				ParticipantAddr: pAddr,
				MessageName:     "new_participant_joined_room",
			})
		}
	}(a.App.Rooms[req.RoomId].Participants, req.ParticipantAddr)
	return nil
}

type RequireRenegotiateSessionReq struct {
	ParticipantAddr string
	RoomId          string
	SessionId       string
	SdpOffer        string
}
type RequireRenegotiateSessionRes struct {
	MessageName string `json:"message_name"`
	SessionId   string `json:"session_id"`
	SdpOffer    string `json:"sdp_offer"`
}

func (a *app) RequireRenegotiateSession(req RequireRenegotiateSessionReq) error {
	u := a.Users[req.ParticipantAddr]
	a.Lock()
	defer a.Unlock()
	err := u.Conn.WriteJSON(RequireRenegotiateSessionRes{
		MessageName: "require_renegotiate_session",
		SessionId:   req.SessionId,
		SdpOffer:    req.SdpOffer,
	})
	if err != nil {
		a.log.Error(err.Error())
	}
	return err
}

type ClientRenegotiateSessionReq struct {
	SessionId string `json:"session_id"`
	SdpAnswer string `json:"sdp_answer"`
}

func (a *app) ClientRenegotiateSession(req ClientRenegotiateSessionReq) error {
	err := a.gcl.RenegotiateSession(req.SessionId, req.SdpAnswer)
	if err != nil {
		a.log.Error(err.Error())
	}
	return err
}

func (a *app) PullRoom(roomId string) error {
	err := a.gcl.RoomPull(roomId)
	if err != nil {
		a.log.Error(err.Error())
	}
	return err
}
