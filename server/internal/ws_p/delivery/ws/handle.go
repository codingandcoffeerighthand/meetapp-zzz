package ws_server

type Message struct {
	MessageName string `json:"message_name"`
	Data        any    `json:"data"`
}
type ClientRenegotiateSession struct {
	SdpAnswer string `json:"sdp_answer"`
	SessionId string `json:"session_id"`
}
type PullRoom struct {
	RoomId string `json:"room_id"`
}
