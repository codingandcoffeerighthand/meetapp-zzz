package domain

type JoinRoomEvent struct {
	RoomID    string  `json:"room_id"`
	SessionID string  `json:"session_id"`
	Tracks    []Track `json:"tracks"`
	SdpOffer  string  `json:"sdp_offer"`
}

var EventJoinedRoomName = "joined_room"

type AddTracksEvent struct {
	RoomID    string  `json:"room_id"`
	SessionID string  `json:"session_id"`
	Tracks    []Track `json:"tracks"`
	SdpOffer  string  `json:"sdp_offer"`
}

var EventAddTracksName = "add_tracks"

type RemoveTracksEvent struct {
	RoomID    string  `json:"room_id"`
	SessionID string  `json:"session_id"`
	Tracks    []Track `json:"tracks"`
	SdpOffer  string  `json:"sdp_offer"`
}

var EventRemoveTracksName = "remove_tracks"

type LeaveRoomEvent struct {
	RoomID    string `json:"room_id"`
	SessionID string `json:"session_id"`
}

var EventLeaveRoomName = "leave_room"

type BackendEvent struct {
	RoomID    string `json:"room_id"`
	SessionID string `json:"session_id"`
	EventType string `json:"event_type"`
	Data      []byte `json:"data"`
}

var EventBackendName = "backend_event"

type FrontendEvent struct {
	RoomID    string `json:"room_id"`
	SessionID string `json:"session_id"`
	EventType string `json:"event_type"`
	Data      []byte `json:"data"`
}

var EventFrontendName = "frontend_event"

type NewSessionEvent struct {
	RoomID       string `json:"room_id"`
	OldSessionID string `json:"old_session_id"`
	NewSessionID string `json:"new_session_id"`
}

var EventNewSessionName = "new_session"

type LocalConnectedEvent struct {
	RoomID string `json:"room_id"`
}

var EventLocalConnectedName = "local_connected"

type RemoteInfoEvent struct {
	RoomID          string `json:"room_id"`
	RemoteSessionID string `json:"remote_session_id"`
	SdpOffer        string `json:"sdp_offer"`
}

var RemoteInfoEventName = "remote_connect"

type RemoteConnectedEvent struct {
	RoomID          string `json:"room_id"`
	RemoteSessionID string `json:"remote_session_id"`
	SdpAnswer       string `json:"sdp_answer"`
}

var EventRemoteConnectedName = "remote_connected"

type RemoteConectSuccessEvent struct {
	RoomID Room    `json:"room"`
	Tracks []Track `json:"tracks"`
}

var EventRemoteConectSuccessName = "remote_connect_success"
