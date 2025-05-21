package domain

type Room struct {
	RoomID       string        `json:"room_id"`
	Name         string        `json:"name"`
	Creator      string        `json:"creator"`
	Participants []Participant `json:"participants"`
}

type Participant struct {
	WalletAddress string  `json:"wallet_address"`
	Name          string  `json:"name"`
	SessionID     string  `json:"session_id"`
	Tracks        []Track `json:"tracks"`
}

func (r *Room) GetAllTracks() []Track {
	tracks := make([]Track, 0)
	for _, participant := range r.Participants {
		tracks = append(tracks, participant.Tracks...)
	}
	return tracks
}
