package domain

/*
	JOINED EVENT
*/

type EventJoinedRoom struct {
	EventName string `json:"event_name"`
	SdpAnswer string `json:"sdp_answer"`
}

const EventJoinedRoomName = "joined_room"

func (e *EventJoinedRoom) ContractEventName() string {
	return e.EventName
}

func (e *EventJoinedRoom) SetName() {
	e.EventName = EventJoinedRoomName
}

/*
	NEW PARTICIPANT JOINED EVENT
*/

type EventNewParticipantJoined struct {
	EventName       string `json:"event_name"`
	Participant     string `json:"participant"`
	ParticipantName string `json:"participant_name"`
}

const EventNewParticipantJoinedName = "new_participant_joined"

func (e *EventNewParticipantJoined) ContractEventName() string {
	return e.EventName
}

func (e *EventNewParticipantJoined) SetName() {
	e.EventName = EventNewParticipantJoinedName
}

type EventForwardedToBackend struct {
	EventName string `json:"event_name"`
	Data      *any   `json:"data,omitempty"`
}

const EventLocalPeerConnectionSuccess = "local_peer_connection_suscess"
const EventRenegoiateSession = "renegoiate_session"

func (e *EventForwardedToBackend) ContractEventName() string {
	return e.EventName
}

type EventPullTrack struct {
	EventName     string `json:"event_name"`
	SdpOffer      string `json:"sdp_offer"`
	RemoteSession string `json:"remote_session"`
}

const EventPullTrackName = "pull_track"

func (e *EventPullTrack) ContractEventName() string {
	return e.EventName
}

func (e *EventPullTrack) SetName() {
	e.EventName = EventPullTrackName
}
