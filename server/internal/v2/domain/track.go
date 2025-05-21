package domain

type Track struct {
	Mid          string `json:"mid"`
	TrackName    string `json:"track_name"`
	StreamNumber uint   `json:"stream_number"`
	SessionID    string `json:"session_id"`
	Location     string `json:"location"`
}
