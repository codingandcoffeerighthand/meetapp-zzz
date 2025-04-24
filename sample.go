package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	// "golang.org/x/net/websocket"
	"slices"
)

// --- Constants and Configuration ---

var (
	authRequired           bool
	port                   string
	cloudflareAppID        string
	cloudflareAppSecret    string
	secretKey              string
	cloudflareCallsBaseURL string
	cloudflareBasePath     string
	debug                  bool
)

func initConfig() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Println("Warning: Could not load .env file:", err) // Don't fatal, as we have defaults
	}

	// authRequired = getEnvBool("AUTH_REQUIRED", true)
	authRequired = false
	port = getEnv("PORT", "5000")
	cloudflareAppID = getEnv("CLOUDFLARE_APP_ID", "")
	cloudflareAppSecret = getEnv("CLOUDFLARE_APP_SECRET", "")
	secretKey = getEnv("JWT_SECRET", "thisisjustademokey")
	cloudflareCallsBaseURL = getEnv("CLOUDFLARE_APPS_URL", "https://rtc.live.cloudflare.com/v1/apps")
	cloudflareBasePath = fmt.Sprintf("%s/%s", cloudflareCallsBaseURL, cloudflareAppID)
	debug = getEnvBool("DEBUG", false)

	if cloudflareAppID == "" || cloudflareAppSecret == "" {
		log.Fatal("CLOUDFLARE_APP_ID and CLOUDFLARE_APP_SECRET must be set")
	}

	// Thêm log để kiểm tra biến môi trường
	log.Println("Biến môi trường đã tải:")
	log.Println("CLOUDFLARE_APP_ID:", cloudflareAppID)
	log.Println("CLOUDFLARE_APP_SECRET:", cloudflareAppSecret)
	log.Println("JWT_SECRET:", secretKey)
	log.Println("CLOUDFLARE_APPS_URL:", cloudflareCallsBaseURL)
	log.Println("DEBUG:", debug)
}

// Helper function to get environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func init() {
	// Initialize maps
	rooms.m = make(map[string]*Room)
	users.m = make(map[string]*User)
	wsConnections.m = make(map[string]map[string]*websocket.Conn)
}

// Helper function to get boolean environment variables
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue // Return default if parsing fails
	}
	return b
}

// --- Data Structures ---

type Room struct {
	RoomId       string         `json:"roomId"` // Thêm trường này
	Name         string         `json:"name"`
	Metadata     map[string]any `json:"metadata"`
	Participants []*Participant `json:"participants"`
	CreatedAt    int64          `json:"createdAt"`
	sync.RWMutex                // Protects concurrent access to the room
}

type Participant struct {
	UserID          string   `json:"userId"`
	SessionID       string   `json:"sessionId"`
	CreatedAt       int64    `json:"createdAt"`
	PublishedTracks []string `json:"publishedTracks"`
}

type User struct {
	UserID      string `json:"userId"`
	Username    string `json:"username"`
	IsModerator bool   `json:"isModerator"`
	Role        string `json:"role"`
}

type SessionResponse struct {
	SessionId     string        `json:"sessionId"`
	OtherSessions []SessionInfo `json:"otherSessions"`
}

type SessionInfo struct {
	UserId          string   `json:"userId"`
	SessionId       string   `json:"sessionId"`
	PublishedTracks []string `json:"publishedTracks"`
}

// Use a concurrent-safe map for rooms.
var rooms = struct {
	sync.RWMutex
	m map[string]*Room
}{m: make(map[string]*Room)}

var users = struct {
	sync.RWMutex
	m map[string]*User
}{m: make(map[string]*User)}

var wsConnections = struct {
	sync.RWMutex
	m map[string]map[string]*websocket.Conn
}{m: make(map[string]map[string]*websocket.Conn)}

// --- Cloudflare API Interaction Functions ---

func createCloudflareSession() (string, error) {
	url := fmt.Sprintf("%s/sessions/new", cloudflareBasePath)
	log.Printf("[Cloudflare API] Creating new session: %s", url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to create request: %v", err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to execute request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to read response body: %v", err)
		return "", err
	}

	log.Printf("[Cloudflare API Response] Status: %d, Body: %s", resp.StatusCode, string(body))

	var responseData map[string]any
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("[Cloudflare API Error] Failed to parse response: %v", err)
		return "", err
	}

	sessionID, ok := responseData["sessionId"].(string)
	if !ok {
		log.Printf("[Cloudflare API Error] Session ID not found in response")
		return "", fmt.Errorf("sessionId not found in response: %s", string(body))
	}

	log.Printf("[Cloudflare API Success] Created session: %s", sessionID)
	return sessionID, nil
}

func makeCloudflareRequest(method string, url string, requestBody map[string]any) (map[string]any, error, int) {
	var reqBodyBytes []byte = nil
	if requestBody != nil {
		var errMarshal error
		reqBodyBytes, errMarshal = json.Marshal(requestBody)
		if errMarshal != nil {
			log.Printf("Lỗi Marshal body yêu cầu JSON: %v", errMarshal)
			return nil, errMarshal, http.StatusInternalServerError
		}
	}

	req, errNewRequest := http.NewRequest(method, url, bytes.NewBuffer(reqBodyBytes))
	if errNewRequest != nil {
		log.Printf("Lỗi tạo yêu cầu HTTP mới: %v", errNewRequest)
		return nil, errNewRequest, http.StatusInternalServerError
	}

	// Thiết lập các header quan trọng và phổ biến
	req.Header.Set("Authorization", "Bearer "+cloudflareAppSecret) // Đảm bảo cloudflareAppSecret có giá trị
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json") // Cho các yêu cầu có body JSON
	}
	req.Header.Set("Accept", "application/json")               // Yêu cầu phản hồi JSON
	req.Header.Set("User-Agent", "CloudflareCalls-Backend-Go") // Thêm User-Agent để nhận diện backend (tùy chọn)
	// Bạn có thể thêm các header khác nếu cần, ví dụ:
	// req.Header.Set("Cache-Control", "no-cache")
	// req.Header.Set("Connection", "keep-alive")

	log.Println("Yêu cầu HTTP:")
	log.Println("  Method:", method)
	log.Println("  URL:", url)
	log.Println("  Headers:", req.Header)
	if requestBody != nil {
		log.Println("  Body yêu cầu:", string(reqBodyBytes))
	}

	client := &http.Client{}
	resp, errClientDo := client.Do(req)
	if errClientDo != nil {
		log.Printf("Lỗi khi thực hiện yêu cầu HTTP: %v", errClientDo)
		return nil, errClientDo, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	respBodyBytes, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		log.Printf("Lỗi đọc body phản hồi: %v", errReadAll)
		return nil, errReadAll, http.StatusInternalServerError
	}

	log.Println("Phản hồi HTTP:")
	log.Println("  Mã trạng thái:", resp.StatusCode)
	log.Println("  Body phản hồi:", string(respBodyBytes))

	var responseData map[string]any
	errUnmarshal := json.Unmarshal(respBodyBytes, &responseData)
	if errUnmarshal != nil {
		log.Printf("Lỗi Unmarshal body phản hồi JSON: %v", errUnmarshal)
		return nil, errUnmarshal, http.StatusInternalServerError
	}

	return responseData, nil, resp.StatusCode
}

func unpublishToCloudflare(cfUrl string, requestBody map[string]any) (map[string]any, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", cfUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var responseData map[string]any

	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func renegotiateWithCloudflare(sessionId string, body map[string]any) (map[string]any, error) {
	url := fmt.Sprintf("%s/sessions/%s/renegotiate", cloudflareBasePath, sessionId)

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData map[string]any
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func manageDataChannelsWithCloudflare(cfUrl string, dataChannels []map[string]any) (map[string]any, error) {
	jsonData, err := json.Marshal(gin.H{"dataChannels": dataChannels}) // Correct structure
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", cfUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData map[string]any
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}

	return responseData, nil
}

func getSessionStateFromCloudflare(sessionId string) (map[string]any, error) {
	url := fmt.Sprintf("%s/sessions/%s", cloudflareBasePath, sessionId)
	log.Printf("[Cloudflare API] Getting session state: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to create request: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to execute request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to read response body: %v", err)
		return nil, err
	}

	log.Printf("[Cloudflare API Response] Status: %d, Body: %s", resp.StatusCode, string(body))

	var responseData map[string]any
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("[Cloudflare API Error] Failed to parse response: %v", err)
		return nil, err
	}

	return responseData, nil
}

// --- Helper Functions ---

func serializeRoom(roomId string, room *Room) gin.H {
	return gin.H{
		"roomId":           roomId,
		"name":             room.Name,
		"metadata":         room.Metadata,
		"participantCount": len(room.Participants),
		"createdAt":        room.CreatedAt,
	}
}

// --- Route Handlers ---

func processCreateRoom(name string, metadata map[string]any) map[string]any {
	if name == "" {
		return map[string]any{"error": "Room name is required"}
	}

	roomId := uuid.NewString()
	room := &Room{
		RoomId:       roomId,
		Name:         name,
		Metadata:     metadata,
		Participants: []*Participant{},
		CreatedAt:    time.Now().Unix(),
	}

	rooms.Lock()
	rooms.m[roomId] = room
	rooms.Unlock()

	return map[string]any{
		"roomId":           room.RoomId,
		"name":             room.Name,
		"metadata":         room.Metadata,
		"participantCount": len(room.Participants),
		"createdAt":        room.CreatedAt,
	}
}

func createRoomWS(ws *websocket.Conn, payload map[string]any) {
	name, _ := payload["name"].(string)
	metadata, _ := payload["metadata"].(map[string]any)
	response := processCreateRoom(name, metadata)

	ws.WriteJSON(map[string]any{
		"type":    "create-room",
		"payload": response,
	})
}

func processInspectRooms() map[string]any {
	if os.Getenv("NODE_ENV") != "development" {
		return map[string]any{"error": "This endpoint is only available in development mode."}
	}

	rooms.RLock()
	defer rooms.RUnlock()

	users.RLock()
	defer users.RUnlock()

	wsConnections.RLock()
	defer wsConnections.RUnlock()

	return map[string]any{
		"rooms":         rooms.m,
		"roomCount":     len(rooms.m),
		"users":         users.m,
		"wsConnections": wsConnections.m,
		"raw":           rooms.m,
	}
}

func inspectRoomsWS(ws *websocket.Conn) {
	response := processInspectRooms()

	ws.WriteJSON(map[string]any{
		"type":    "inspect-rooms",
		"payload": response,
	})
}

func processJoinRoom(roomId string, currentUser *User) (SessionResponse, map[string]any) {
	if currentUser == nil {
		return SessionResponse{}, map[string]any{"error": "Forbidden: Invalid user"}
	}

	userId := currentUser.UserID

	// Tạo room nếu chưa tồn tại
	rooms.Lock()
	room, exists := rooms.m[roomId]
	if !exists {
		room = &Room{
			Name:         "New Room",
			Metadata:     make(map[string]any),
			Participants: make([]*Participant, 0),
			CreatedAt:    time.Now().Unix(),
		}
		rooms.m[roomId] = room
	}
	rooms.Unlock()

	// Tạo Calls Session từ Cloudflare
	sessionID, err := createCloudflareSession()
	if err != nil {
		return SessionResponse{}, map[string]any{"error": "Could not create Calls session"}
	}

	participant := &Participant{
		UserID:          userId,
		SessionID:       sessionID,
		CreatedAt:       time.Now().Unix(),
		PublishedTracks: make([]string, 0),
	}

	room.Lock()
	room.Participants = append(room.Participants, participant)

	// Lấy thông tin các session khác trong room
	otherSessions := make([]SessionInfo, 0)
	for _, p := range room.Participants {
		if p.UserID != userId {
			state, err := getSessionStateFromCloudflare(p.SessionID)
			if err != nil {
				log.Printf("Error getting session state for %s: %v", p.SessionID, err)
				continue
			}

			tracks := make([]string, 0)
			if tracksArray, ok := state["tracks"].([]any); ok {
				for _, track := range tracksArray {
					if trackMap, ok := track.(map[string]any); ok {
						if trackName, ok := trackMap["trackName"].(string); ok {
							tracks = append(tracks, trackName)
						}
					}
				}
			}

			otherSessions = append(otherSessions, SessionInfo{
				UserId:          p.UserID,
				SessionId:       p.SessionID,
				PublishedTracks: tracks,
			})
		}
	}
	room.Unlock()

	// Khởi tạo WebSocket connections cho phòng nếu chưa có
	wsConnections.Lock()
	if _, exists := wsConnections.m[roomId]; !exists {
		wsConnections.m[roomId] = make(map[string]*websocket.Conn)
	}
	wsConnections.Unlock()

	// Thông báo cho những người khác trong phòng
	broadcastToRoom(roomId, gin.H{
		"type": "participant-joined",
		"payload": gin.H{
			"userId":    userId,
			"username":  currentUser.Username,
			"sessionId": sessionID,
		},
	}, userId)

	response := SessionResponse{
		SessionId:     sessionID,
		OtherSessions: otherSessions,
	}

	if debug {
		log.Printf("Join room response: %+v\n", response)
	}

	return response, nil
}

func getCurrentUser(payload map[string]any) (*User, error) {
	// Lấy userId từ payload
	userId, ok := payload["userId"].(string)
	if !ok || userId == "" {
		return nil, fmt.Errorf("invalid or missing userId")
	}

	// Lấy thông tin user từ hệ thống (map users.m hoặc database)
	users.RLock() // Đọc dữ liệu an toàn
	currentUser, exists := users.m[userId]
	users.RUnlock()

	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return currentUser, nil
}

func joinRoomWS(ws *websocket.Conn, payload map[string]any) {
	// Lấy roomId từ payload
	roomId, ok := payload["roomId"].(string)
	if !ok || roomId == "" {
		ws.WriteJSON(map[string]any{
			"type":    "join-room",
			"payload": map[string]any{"error": "Invalid or missing roomId"},
		})
		return
	}

	// Lấy User từ function `getUser`
	currentUser, err := getCurrentUser(payload)
	if err != nil {
		ws.WriteJSON(map[string]any{
			"type":    "join-room",
			"payload": map[string]any{"error": err.Error()},
		})
		return
	}

	// Gọi `processJoinRoom` để xử lý tham gia phòng
	response, errResp := processJoinRoom(roomId, currentUser)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "join-room",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "join-room",
		"payload": response,
	})
}

// Sửa lại struct request cho publishTracks để match với Express
func processPublishTracks(roomId, sessionId string, offer map[string]any, tracks []struct {
	TrackName string `json:"trackName"`
	Mid       string `json:"mid"`
	Location  string `json:"location"`
}) (map[string]any, map[string]any) {
	// Kiểm tra room có tồn tại không
	rooms.RLock()
	room, ok := rooms.m[roomId]
	rooms.RUnlock()
	if !ok {
		return nil, map[string]any{"error": "Room not found"}
	}

	// Tìm participant trong room
	var participant *Participant
	room.RLock()
	for _, p := range room.Participants {
		if p.SessionID == sessionId {
			participant = p
			break
		}
	}
	room.RUnlock()
	if participant == nil {
		return nil, map[string]any{"error": "Session not found in this room"}
	}

	if debug {
		log.Printf("Publishing tracks for session %s: %+v\n", sessionId, tracks)
	}

	// Chuẩn bị dữ liệu track cho Cloudflare API
	tracksData := make([]map[string]any, len(tracks))
	trackNames := make([]string, len(tracks))
	for i, t := range tracks {
		location := t.Location
		if location == "" {
			location = "local"
		}
		tracksData[i] = map[string]any{
			"trackName": t.TrackName,
			"location":  location,
			"mid":       t.Mid,
		}
		trackNames[i] = t.TrackName
	}

	// Gửi yêu cầu đến Cloudflare API
	requestBody := map[string]any{
		"sessionDescription": offer,
		"tracks":             tracksData,
	}

	url := fmt.Sprintf("%s/sessions/%s/tracks/new", cloudflareBasePath, sessionId)
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, map[string]any{"error": err.Error()}
	}

	cfReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, map[string]any{"error": err.Error()}
	}

	cfReq.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)
	cfReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(cfReq)
	if err != nil {
		return nil, map[string]any{"error": err.Error()}
	}
	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, map[string]any{"error": err.Error()}
	}

	// Cập nhật danh sách track đã publish của participant
	room.Lock()
	for _, t := range tracks {
		if !contains(participant.PublishedTracks, t.TrackName) {
			participant.PublishedTracks = append(participant.PublishedTracks, t.TrackName)
		}
	}
	room.Unlock()

	// Gửi sự kiện track-published đến room
	if data["sessionDescription"] != nil {
		broadcastToRoom(roomId, gin.H{
			"type": "track-published",
			"payload": gin.H{
				"userId":    participant.UserID,
				"sessionId": sessionId,
				"tracks":    trackNames,
			},
		}, participant.UserID)

		if debug {
			log.Printf("Track published event broadcasted for session %s with tracks: %v\n", sessionId, trackNames)
		}
	}

	return data, nil
}

func publishTracksWS(ws *websocket.Conn, payload map[string]any) {
	roomId, _ := payload["roomId"].(string)
	sessionId, _ := payload["sessionId"].(string)

	offer, _ := payload["offer"].(map[string]any)
	tracksData, _ := payload["tracks"].([]any)

	tracks := make([]struct {
		TrackName string `json:"trackName"`
		Mid       string `json:"mid"`
		Location  string `json:"location"`
	}, len(tracksData))

	for i, t := range tracksData {
		trackMap, ok := t.(map[string]any)
		if !ok {
			ws.WriteJSON(map[string]any{
				"type":    "publish-tracks",
				"payload": map[string]any{"error": "Invalid track format"},
			})
			return
		}
		tracks[i] = struct {
			TrackName string `json:"trackName"`
			Mid       string `json:"mid"`
			Location  string `json:"location"`
		}{
			TrackName: trackMap["trackName"].(string),
			Mid:       trackMap["mid"].(string),
			Location:  trackMap["location"].(string),
		}
	}

	response, errResp := processPublishTracks(roomId, sessionId, offer, tracks)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "publish-tracks",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "publish-tracks",
		"payload": response,
	})
}

// Helper function to check if slice contains string
func contains(slice []string, str string) bool {
	return slices.Contains(slice, str)
}

func processUnpublishTrack(roomId, sessionId string, currentUser *User, trackName, mid string, force bool, sessionDescription map[string]any) (map[string]any, map[string]any) {
	// Kiểm tra user hợp lệ
	if currentUser == nil {
		return nil, map[string]any{"error": "Forbidden: Invalid user"}
	}

	// Nếu yêu cầu force unpublish track của người khác
	if force && sessionId != currentUser.UserID {
		if !currentUser.IsModerator {
			return nil, map[string]any{
				"errorCode":        "NOT_AUTHORIZED",
				"errorDescription": "Only moderators can force unpublish other participants' tracks",
			}
		}
	}

	if debug {
		log.Println("Unpublishing track:", map[string]any{
			"roomId":    roomId,
			"sessionId": sessionId,
			"trackName": trackName,
			"mid":       mid,
			"force":     force,
		})
	}

	// Kiểm tra đầu vào hợp lệ
	if mid == "" {
		return nil, map[string]any{
			"errorCode":        "INVALID_REQUEST",
			"errorDescription": "mid is required to unpublish a track.",
		}
	}

	if sessionDescription == nil {
		return nil, map[string]any{
			"errorCode":        "INVALID_REQUEST",
			"errorDescription": "sessionDescription is required to unpublish a track.",
		}
	}

	// Gọi API Cloudflare để unpublish track
	cfUrl := fmt.Sprintf("%s/sessions/%s/tracks/close", cloudflareBasePath, sessionId)
	if debug {
		log.Println("Calling Cloudflare API:", cfUrl)
	}

	requestBody := map[string]any{
		"tracks": []map[string]string{
			{"mid": mid},
		},
		"force":              force,
		"sessionDescription": sessionDescription,
	}

	if debug {
		log.Printf("Request body: %+v\n", requestBody)
	}

	data, err := unpublishToCloudflare(cfUrl, requestBody)
	if err != nil {
		return nil, map[string]any{"error": err.Error()}
	}

	if debug {
		log.Println("Cloudflare API response:", data)
	}

	// Gửi sự kiện track-unpublished đến room
	broadcastToRoom(roomId, gin.H{
		"type":    "track-unpublished",
		"payload": gin.H{"sessionId": sessionId, "trackName": trackName},
	}, sessionId)

	return data, nil
}

func unpublishTrackWS(ws *websocket.Conn, payload map[string]any) {
	roomId, _ := payload["roomId"].(string)
	sessionId, _ := payload["sessionId"].(string)

	// Lấy User từ function `getUser`
	currentUser, err := getCurrentUser(payload)
	if err != nil {
		ws.WriteJSON(map[string]any{
			"type":    "unpublish-track",
			"payload": map[string]any{"error": err.Error()},
		})
		return
	}

	trackName, _ := payload["trackName"].(string)
	mid, _ := payload["mid"].(string)
	force, _ := payload["force"].(bool)
	sessionDescription, _ := payload["sessionDescription"].(map[string]any)

	response, errResp := processUnpublishTrack(roomId, sessionId, currentUser, trackName, mid, force, sessionDescription)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "unpublish-track",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "unpublish-track",
		"payload": response,
	})
}

func processPullTracks(roomId, sessionId, remoteSessionId, trackName string) (map[string]any, map[string]any, int) {
	// Lấy thông tin room từ bộ nhớ
	rooms.RLock()
	room, ok := rooms.m[roomId]
	rooms.RUnlock()
	if !ok {
		return nil, map[string]any{"error": "Room not found"}, http.StatusNotFound
	}

	// Tìm participant trong room
	participant := findParticipantBySessionId(room, sessionId)
	if participant == nil {
		return nil, map[string]any{"error": "Session not found in this room"}, http.StatusNotFound
	}

	// Tạo cấu trúc tracksToPull cho Cloudflare API
	tracksToPull := []map[string]any{
		{
			"location":  "remote",
			"sessionId": remoteSessionId,
			"trackName": trackName,
		},
	}

	// Gọi Cloudflare API để pull track
	url := fmt.Sprintf("%s/sessions/%s/tracks/new", cloudflareBasePath, sessionId)
	requestBody := map[string]any{
		"tracks": tracksToPull,
	}

	data, err, statusCode := makeCloudflareRequest("POST", url, requestBody)
	if err != nil {
		return nil, map[string]any{"error": "Failed to pull track", "detail": err.Error()}, http.StatusInternalServerError
	}

	// Kiểm tra mã trạng thái lỗi từ Cloudflare API
	if statusCode >= 400 {
		return nil, map[string]any{"error": "Cloudflare API error during pull track", "detail": data}, statusCode
	}

	return data, nil, http.StatusOK
}

func pullTracksAPI(c *gin.Context) {
	roomId := c.Param("roomId")
	sessionId := c.Param("sessionId")

	var req struct {
		RemoteSessionId string `json:"remoteSessionId"`
		TrackName       string `json:"trackName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, errResp, statusCode := processPullTracks(roomId, sessionId, req.RemoteSessionId, req.TrackName)
	if errResp != nil {
		c.JSON(statusCode, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func pullTracksWS(ws *websocket.Conn, payload map[string]any) {
	roomId, _ := payload["roomId"].(string)
	sessionId, _ := payload["sessionId"].(string)
	remoteSessionId, _ := payload["remoteSessionId"].(string)
	trackName, _ := payload["trackName"].(string)

	response, errResp, statusCode := processPullTracks(roomId, sessionId, remoteSessionId, trackName)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "pull-tracks",
			"payload": errResp,
			"status":  statusCode,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "pull-tracks",
		"payload": response,
		"status":  http.StatusOK,
	})
}

// Helper function để tìm participant theo sessionId (tái sử dụng từ code của bạn)
func findParticipantBySessionId(room *Room, sessionId string) *Participant {
	room.RLock()
	defer room.RUnlock()
	for _, p := range room.Participants {
		if p.SessionID == sessionId {
			return p
		}
	}
	return nil
}

func processRenegotiateSession(sessionId, sdp, sdpType string) (map[string]any, map[string]any) {
	// Chuẩn bị body gửi đến Cloudflare
	body := map[string]any{
		"sessionDescription": map[string]string{
			"sdp":  sdp,
			"type": sdpType,
		},
	}

	// Gửi yêu cầu renegotiate session đến Cloudflare
	data, err := renegotiateWithCloudflare(sessionId, body)
	if err != nil {
		if data != nil && data["errorCode"] != nil { // Xử lý lỗi từ Cloudflare
			return nil, data
		}
		return nil, map[string]any{"error": err.Error()}
	}

	return data, nil
}
func renegotiateSessionAPI(c *gin.Context) {
	sessionId := c.Param("sessionId")

	var req struct {
		SDP     string `json:"sdp"`
		SDPType string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, errResp := processRenegotiateSession(sessionId, req.SDP, req.SDPType)
	if errResp != nil {
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func renegotiateSessionWS(ws *websocket.Conn, payload map[string]any) {
	sessionId, _ := payload["sessionId"].(string)
	sdp, _ := payload["sdp"].(string)
	sdpType, _ := payload["type"].(string)

	response, errResp := processRenegotiateSession(sessionId, sdp, sdpType)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "renegotiate-session",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "renegotiate-session",
		"payload": response,
	})
}

func processManageDataChannels(roomId, sessionId string, dataChannels []struct {
	Location        string `json:"location"`
	DataChannelName string `json:"dataChannelName"`
}) (map[string]any, map[string]any) {
	// Kiểm tra room có tồn tại không
	rooms.RLock()
	room, ok := rooms.m[roomId]
	print(room)
	rooms.RUnlock()
	if !ok {
		return nil, map[string]any{"error": "Room not found"}
	}

	// Chuẩn bị dữ liệu gửi đến Cloudflare
	cfUrl := fmt.Sprintf("%s/sessions/%s/datachannels/new", cloudflareBasePath, sessionId)
	dataChannelsRequest := make([]map[string]any, len(dataChannels))
	for i, dc := range dataChannels {
		dataChannelsRequest[i] = map[string]any{
			"location":        dc.Location,
			"dataChannelName": dc.DataChannelName,
		}
	}

	// Gọi Cloudflare API
	data, err := manageDataChannelsWithCloudflare(cfUrl, dataChannelsRequest)
	if err != nil {
		if data != nil && data["errorCode"] != nil { // Xử lý lỗi từ Cloudflare
			return nil, data
		}
		return nil, map[string]any{"error": err.Error()}
	}

	// Nếu user đang publish channel, có thể lưu vào `participant.publishedDataChannels`
	for _, dc := range dataChannels {
		if dc.Location == "local" {
			// TODO: Lưu trữ danh sách data channels đã publish
			// E.g. participant.publishedDataChannels = append(participant.publishedDataChannels, dc.DataChannelName)
		}
	}

	return data, nil
}
func manageDataChannelsAPI(c *gin.Context) {
	roomId := c.Param("roomId")
	sessionId := c.Param("sessionId")

	var req struct {
		DataChannels []struct {
			Location        string `json:"location"`
			DataChannelName string `json:"dataChannelName"`
		} `json:"dataChannels"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, errResp := processManageDataChannels(roomId, sessionId, req.DataChannels)
	if errResp != nil {
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func manageDataChannelsWS(ws *websocket.Conn, payload map[string]any) {
	roomId, _ := payload["roomId"].(string)
	sessionId, _ := payload["sessionId"].(string)

	channelsData, _ := payload["dataChannels"].([]any)

	dataChannels := make([]struct {
		Location        string `json:"location"`
		DataChannelName string `json:"dataChannelName"`
	}, len(channelsData))

	for i, ch := range channelsData {
		chMap, ok := ch.(map[string]any)
		if !ok {
			ws.WriteJSON(map[string]any{
				"type":    "manage-datachannels",
				"payload": map[string]any{"error": "Invalid dataChannel format"},
			})
			return
		}
		dataChannels[i] = struct {
			Location        string `json:"location"`
			DataChannelName string `json:"dataChannelName"`
		}{
			Location:        chMap["location"].(string),
			DataChannelName: chMap["dataChannelName"].(string),
		}
	}

	response, errResp := processManageDataChannels(roomId, sessionId, dataChannels)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "manage-datachannels",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "manage-datachannels",
		"payload": response,
	})
}

func processGetParticipants(roomId string) (map[string]any, map[string]any) {
	// Lấy danh sách participants của room
	rooms.RLock()
	room, ok := rooms.m[roomId]
	rooms.RUnlock()

	if !ok {
		return nil, map[string]any{"error": "Room not found"}
	}

	return map[string]any{"participants": room.Participants}, nil
}
func getParticipantsAPI(c *gin.Context) {
	roomId := c.Param("roomId")

	response, errResp := processGetParticipants(roomId)
	if errResp != nil {
		c.JSON(http.StatusNotFound, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func getParticipantsWS(ws *websocket.Conn, payload map[string]any) {
	roomId, _ := payload["roomId"].(string)

	response, errResp := processGetParticipants(roomId)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "get-participants",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "get-participants",
		"payload": response,
	})
}

func processGetParticipantTracks(roomId, sessionId string) (map[string]any, map[string]any) {
	// Kiểm tra room có tồn tại không
	rooms.RLock()
	room, ok := rooms.m[roomId]
	rooms.RUnlock()
	if !ok {
		return nil, map[string]any{"error": "Room not found"}
	}

	// Tìm participant trong room
	var participant *Participant
	room.RLock()
	for _, p := range room.Participants {
		if p.SessionID == sessionId {
			participant = p
			break
		}
	}
	room.RUnlock()

	if participant == nil {
		return nil, map[string]any{"error": "Participant not found"}
	}

	return map[string]any{"publishedTracks": participant.PublishedTracks}, nil
}
func getParticipantTracksAPI(c *gin.Context) {
	roomId := c.Param("roomId")
	sessionId := c.Param("sessionId")

	response, errResp := processGetParticipantTracks(roomId, sessionId)
	if errResp != nil {
		c.JSON(http.StatusNotFound, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func getParticipantTracksWS(ws *websocket.Conn, payload map[string]any) {
	roomId, _ := payload["roomId"].(string)
	sessionId, _ := payload["sessionId"].(string)

	response, errResp := processGetParticipantTracks(roomId, sessionId)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "get-participant-tracks",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "get-participant-tracks",
		"payload": response,
	})
}

func processGetICEServers() map[string]any {
	cloudflareTurnID := os.Getenv("CLOUDFLARE_TURN_ID")
	cloudflareTurnToken := os.Getenv("CLOUDFLARE_TURN_TOKEN")

	if cloudflareTurnID == "" || cloudflareTurnToken == "" {
		return map[string]any{
			"iceServers": []map[string]any{
				{"urls": "stun:stun.cloudflare.com:3478"},
			},
		}
	}

	lifetime := 600 // Credentials valid for 10 minutes (600 seconds)
	timestamp := time.Now().Unix() + int64(lifetime)
	username := fmt.Sprintf("%d:%s", timestamp, cloudflareTurnID)

	// Create HMAC-SHA256 hash using CLOUDFLARE_TURN_TOKEN as the key
	h := hmac.New(sha256.New, []byte(cloudflareTurnToken))
	h.Write([]byte(username))
	credential := base64.StdEncoding.EncodeToString(h.Sum(nil))

	iceServers := map[string]any{
		"iceServers": []map[string]any{
			{"urls": "stun:stun.cloudflare.com:3478"},
			{
				"urls":       "turn:turn.cloudflare.com:3478?transport=udp",
				"username":   username,
				"credential": credential,
			},
			{
				"urls":       "turn:turn.cloudflare.com:3478?transport=tcp",
				"username":   username,
				"credential": credential,
			},
			{
				"urls":       "turns:turn.cloudflare.com:5349?transport=tcp",
				"username":   username,
				"credential": credential,
			},
		},
	}

	return iceServers
}
func getICEServersAPI(c *gin.Context) {
	response := processGetICEServers()
	c.JSON(http.StatusOK, response)
}
func getICEServersWS(ws *websocket.Conn) {
	response := processGetICEServers()

	ws.WriteJSON(map[string]any{
		"type":    "get-ice-servers",
		"payload": response,
	})
}

func processGetToken(username string) map[string]any {
	if username == "" {
		return map[string]any{"error": "Username is required"}
	}

	userId := uuid.NewString()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":      userId,
		"username":    username,
		"role":        "demo",
		"isModerator": true, // Should come from DB
		"exp":         time.Now().Add(time.Hour * 8).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return map[string]any{"error": "Could not generate token"}
	}

	// Store user info
	users.Lock()
	users.m[userId] = &User{
		UserID:      userId,
		Username:    username,
		IsModerator: true,
		Role:        "demo",
	}
	users.Unlock()

	return map[string]any{"userId": userId, "token": tokenString}
}

func getTokenAPI(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := processGetToken(req.Username)
	c.JSON(http.StatusOK, response)
}
func getTokenWS(ws *websocket.Conn, payload map[string]any) {
	username, _ := payload["username"].(string)
	response := processGetToken(username)
	ws.WriteJSON(map[string]any{
		"type":    "auth-token",
		"payload": response,
	})
}

func processGetSessionState(sessionId string) (map[string]any, map[string]any) {
	// Lấy trạng thái session từ Cloudflare
	data, err := getSessionStateFromCloudflare(sessionId)
	if err != nil {
		return nil, map[string]any{
			"errorCode":        "SESSION_STATE_ERROR",
			"errorDescription": err.Error(),
		}
	}

	return data, nil
}
func getSessionStateAPI(c *gin.Context) {
	sessionId := c.Param("sessionId")

	response, errResp := processGetSessionState(sessionId)
	if errResp != nil {
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func getSessionStateWS(ws *websocket.Conn, payload map[string]any) {
	sessionId, _ := payload["sessionId"].(string)

	response, errResp := processGetSessionState(sessionId)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "get-session-state",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "get-session-state",
		"payload": response,
	})
}

func processGetUserInfo(userIdParam string, currentUser *User) (map[string]any, map[string]any) {
	// Kiểm tra user hợp lệ
	if currentUser == nil {
		return nil, map[string]any{"error": "Forbidden: Invalid user"}
	}

	// Nếu yêu cầu thông tin của chính user ("me")
	if userIdParam == "me" {
		users.RLock()
		userInfo, ok := users.m[currentUser.UserID]
		users.RUnlock()
		if !ok {
			return nil, map[string]any{
				"errorCode":        "USER_NOT_FOUND",
				"errorDescription": "Current user not found",
			}
		}
		return map[string]any{"userInfo": userInfo}, nil
	}

	// Lấy thông tin user khác
	users.RLock()
	requestedUser, ok := users.m[userIdParam]
	users.RUnlock()
	if !ok {
		return nil, map[string]any{
			"errorCode":        "USER_NOT_FOUND",
			"errorDescription": "User not found",
		}
	}

	// Trả về thông tin giới hạn của user khác
	return map[string]any{
		"userId":   requestedUser.UserID,
		"username": requestedUser.Username,
	}, nil
}
func getUserInfoAPI(c *gin.Context) {
	userIdParam := c.Param("userId")
	user, _ := c.Get("user")
	currentUser, ok := user.(*User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: Invalid user"})
		return
	}

	response, errResp := processGetUserInfo(userIdParam, currentUser)
	if errResp != nil {
		c.JSON(http.StatusNotFound, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func getUserInfoWS(ws *websocket.Conn, payload map[string]any) {
	userIdParam, _ := payload["userId"].(string)

	// Lấy User từ function `getUser`
	currentUser, err := getCurrentUser(payload)
	if err != nil {
		ws.WriteJSON(map[string]any{
			"type":    "join-room",
			"payload": map[string]any{"error": err.Error()},
		})
		return
	}

	response, errResp := processGetUserInfo(userIdParam, currentUser)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "get-user-info",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "get-user-info",
		"payload": response,
	})
}

func processLeaveRoom(roomId, sessionId string) (map[string]any, map[string]any) {
	// Lấy room từ danh sách
	rooms.RLock()
	room, ok := rooms.m[roomId]
	rooms.RUnlock()
	if !ok {
		return nil, map[string]any{"error": "Room not found"}
	}

	// Tìm participant trong danh sách
	participantIndex := -1
	var participant *Participant
	room.Lock()
	for i, p := range room.Participants {
		if p.SessionID == sessionId {
			participantIndex = i
			participant = p
			break
		}
	}

	// Nếu tìm thấy participant
	if participantIndex != -1 {
		room.Participants = append(room.Participants[:participantIndex], room.Participants[participantIndex+1:]...)

		// Thông báo rời phòng đến những người khác
		broadcastToRoom(roomId, gin.H{
			"type": "participant-left",
			"payload": gin.H{
				"sessionId": sessionId,
				"userId":    participant.UserID,
			},
		}, sessionId)

		// Nếu phòng trống, xóa phòng
		if len(room.Participants) == 0 {
			rooms.Lock()
			delete(rooms.m, roomId)
			rooms.Unlock()
		}
	}
	room.Unlock()

	return map[string]any{"success": true}, nil
}
func leaveRoomAPI(c *gin.Context) {
	roomId := c.Param("roomId")

	var req struct {
		SessionId string `json:"sessionId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, errResp := processLeaveRoom(roomId, req.SessionId)
	if errResp != nil {
		c.JSON(http.StatusNotFound, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func leaveRoomWS(ws *websocket.Conn, payload map[string]any) {
	roomId, _ := payload["roomId"].(string)
	sessionId, _ := payload["sessionId"].(string)

	response, errResp := processLeaveRoom(roomId, sessionId)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "leave-room",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "leave-room",
		"payload": response,
	})
}

func processUpdateTrackStatus(roomId, sessionId string, currentUser *User, req struct {
	TrackId string `json:"trackId"`
	Kind    string `json:"kind"`
	Enabled bool   `json:"enabled"`
	Force   bool   `json:"force"`
}) (map[string]any, map[string]any) {
	// Kiểm tra user hợp lệ
	if currentUser == nil {
		return nil, map[string]any{"error": "Forbidden: Invalid user"}
	}

	// Kiểm tra nếu user muốn ép buộc thay đổi track của người khác
	if req.Force && sessionId != currentUser.UserID {
		if !currentUser.IsModerator {
			return nil, map[string]any{
				"errorCode":        "NOT_AUTHORIZED",
				"errorDescription": "Only moderators can force change other participants' tracks",
			}
		}
	}

	// Gửi thông báo đến các participant khác trong room
	broadcastToRoom(roomId, gin.H{
		"type": "track-status-changed",
		"payload": gin.H{
			"sessionId": sessionId,
			"trackId":   req.TrackId,
			"kind":      req.Kind,
			"enabled":   req.Enabled,
		},
	}, sessionId)

	return map[string]any{"success": true}, nil
}
func updateTrackStatusAPI(c *gin.Context) {
	roomId := c.Param("roomId")
	sessionId := c.Param("sessionId")

	user, _ := c.Get("user")
	currentUser, ok := user.(*User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: Invalid user"})
		return
	}

	var req struct {
		TrackId string `json:"trackId"`
		Kind    string `json:"kind"`
		Enabled bool   `json:"enabled"`
		Force   bool   `json:"force"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, errResp := processUpdateTrackStatus(roomId, sessionId, currentUser, req)
	if errResp != nil {
		c.JSON(http.StatusForbidden, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func updateTrackStatusWS(ws *websocket.Conn, payload map[string]any) {
	roomId, _ := payload["roomId"].(string)
	sessionId, _ := payload["sessionId"].(string)

	// Lấy User từ function `getUser`
	currentUser, err := getCurrentUser(payload)
	if err != nil {
		ws.WriteJSON(map[string]any{
			"type":    "join-room",
			"payload": map[string]any{"error": err.Error()},
		})
		return
	}

	req := struct {
		TrackId string `json:"trackId"`
		Kind    string `json:"kind"`
		Enabled bool   `json:"enabled"`
		Force   bool   `json:"force"`
	}{
		TrackId: payload["trackId"].(string),
		Kind:    payload["kind"].(string),
		Enabled: payload["enabled"].(bool),
		Force:   payload["force"].(bool),
	}

	response, errResp := processUpdateTrackStatus(roomId, sessionId, currentUser, req)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "update-track-status",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "update-track-status",
		"payload": response,
	})
}

func processUpdateRoomMetadata(roomId string, req struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata"`
}) (map[string]any, map[string]any) {
	// Kiểm tra room có tồn tại không
	rooms.RLock()
	room, ok := rooms.m[roomId]
	rooms.RUnlock()
	if !ok {
		return nil, map[string]any{"error": "Room not found"}
	}

	// Cập nhật thông tin phòng
	room.Lock()
	if req.Name != "" {
		room.Name = req.Name
	}

	if req.Metadata != nil {
		if room.Metadata == nil {
			room.Metadata = make(map[string]any)
		}
		for k, v := range req.Metadata {
			room.Metadata[k] = v
		}
	}
	room.Unlock()

	// Gửi thông báo đến các participant trong room
	broadcastToRoom(roomId, gin.H{
		"type": "room-metadata-updated",
		"payload": gin.H{
			"roomId":   roomId,
			"name":     room.Name,
			"metadata": room.Metadata,
		},
		"from": "server",
	}, "") // Không loại trừ bất kỳ user nào

	return serializeRoom(roomId, room), nil
}
func updateRoomMetadataAPI(c *gin.Context) {
	roomId := c.Param("roomId")

	var req struct {
		Name     string         `json:"name"`
		Metadata map[string]any `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, errResp := processUpdateRoomMetadata(roomId, req)
	if errResp != nil {
		c.JSON(http.StatusNotFound, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
func updateRoomMetadataWS(ws *websocket.Conn, payload map[string]any) {
	roomId, _ := payload["roomId"].(string)

	req := struct {
		Name     string         `json:"name"`
		Metadata map[string]any `json:"metadata"`
	}{
		Name:     payload["name"].(string),
		Metadata: payload["metadata"].(map[string]any),
	}

	response, errResp := processUpdateRoomMetadata(roomId, req)

	if errResp != nil {
		ws.WriteJSON(map[string]any{
			"type":    "update-room-metadata",
			"payload": errResp,
		})
		return
	}

	ws.WriteJSON(map[string]any{
		"type":    "update-room-metadata",
		"payload": response,
	})
}

func processGetRooms() map[string]any {
	rooms.RLock()
	defer rooms.RUnlock()

	roomList := make([]gin.H, 0, len(rooms.m))
	for roomId, room := range rooms.m {
		roomList = append(roomList, serializeRoom(roomId, room))
	}

	return map[string]any{"rooms": roomList}
}
func getRoomsAPI(c *gin.Context) {
	response := processGetRooms()
	c.JSON(http.StatusOK, response)
}
func getRoomsWS(ws *websocket.Conn) {
	response := processGetRooms()

	ws.WriteJSON(map[string]any{
		"type":    "get-rooms",
		"payload": response,
	})
}

// --- WebSocket Handling ---

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Cẩn thận với setting này trong production
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(c *gin.Context) {
	if debug {
		log.Printf("Incoming WebSocket request from: %s\n", c.Request.RemoteAddr)
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}
	defer ws.Close()

	handleWebSocket(ws)
}

func handleWebSocket(ws *websocket.Conn) {
	if debug {
		log.Printf("New WebSocket connection established from: %s\n", ws.RemoteAddr().String())
	}
	defer ws.Close()

	var userId string
	var roomId string

	for {
		// Đọc message
		messageType, message, err := ws.ReadMessage()
		log.Println("messageType", messageType)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			handleWSDisconnect(ws, roomId, userId)
			break
		}

		// Parse message
		var data map[string]any
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		messageTypeStr, ok := data["type"].(string)
		if !ok {
			log.Println("Invalid message format: missing or invalid 'type'")
			continue
		}

		payload, _ := data["payload"].(map[string]any)

		switch messageTypeStr {
		case "join-websocket":
			// Format payload giống Express
			payload := map[string]any{
				"roomId": data["payload"].(map[string]any)["roomId"],
				"userId": data["payload"].(map[string]any)["userId"],
				"token":  data["payload"].(map[string]any)["token"],
			}
			handleWSJoin(ws, payload)

		case "data-message":
			// Format payload giống Express
			payload := map[string]any{
				"from": data["payload"].(map[string]any)["from"],
				"data": data["payload"].(map[string]any)["message"],
			}
			handleDataMessage(ws, payload)
		case "auth-token":
			getTokenWS(ws, payload)
		case "create-room":
			createRoomWS(ws, payload)
		case "get-rooms":
			getRoomsWS(ws)
		case "inspect-rooms":
			inspectRoomsWS(ws)
		case "join-room":
			joinRoomWS(ws, payload)
		case "publish-tracks":
			publishTracksWS(ws, payload)
		case "unpublish-track":
			unpublishTrackWS(ws, payload)
		case "pull-tracks":
			pullTracksWS(ws, payload)
		case "renegotiate-session":
			renegotiateSessionWS(ws, payload)
		case "manage-datachannels":
			manageDataChannelsWS(ws, payload)
		case "get-participants":
			getParticipantsWS(ws, payload)
		case "get-participant-tracks":
			getParticipantTracksWS(ws, payload)
		case "get-ice-servers":
			getICEServersWS(ws)
		case "get-session-state":
			getSessionStateWS(ws, payload)
		case "get-user-info":
			getUserInfoWS(ws, payload)
		case "leave-room":
			leaveRoomWS(ws, payload)
		case "update-track-status":
			updateTrackStatusWS(ws, payload)
		case "update-room-metadata":
			updateRoomMetadataWS(ws, payload)
		default:
			log.Println("Unknown message type:", messageTypeStr)
			ws.WriteJSON(map[string]any{
				"type":    "error",
				"payload": map[string]any{"error": "Unknown message type"},
			})
		}
	}
}

// --- WebSocket Helper Functions ---

func getRoomIdByUserId(userId string) string {
	rooms.RLock()
	defer rooms.RUnlock()
	for roomId, room := range rooms.m {
		for _, p := range room.Participants {
			if p.UserID == userId {
				return roomId
			}
		}
	}
	return ""
}

func getRoomIdBySessionId(sessionId string) string {
	rooms.RLock()
	defer rooms.RUnlock()
	for roomId, room := range rooms.m {
		for _, p := range room.Participants {
			if p.SessionID == sessionId {
				return roomId
			}
		}
	}
	return ""
}

func getWebSocketByUserId(userId string) *websocket.Conn {
	wsConnections.RLock()
	defer wsConnections.RUnlock()
	for _, userMap := range wsConnections.m {
		if conn, ok := userMap[userId]; ok {
			return conn
		}
	}
	return nil
}

func handleWSJoin(ws *websocket.Conn, payload map[string]any) {
	roomId, ok1 := payload["roomId"].(string)
	userId, ok2 := payload["userId"].(string)
	token, ok3 := payload["token"].(string)
	//check roomid, userid and token
	if !ok1 || !ok2 || (authRequired && !ok3) {
		log.Println("Missing roomId, userId, or token in WS join")
		_ = ws.WriteJSON(map[string]string{"error": "Missing roomId, userId, or token"})
		return
	}
	if authRequired {
		_, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			log.Println("Invalid token in WS join:", err)
			_ = ws.WriteJSON(map[string]string{"error": "Invalid or expired token"})
			return
		}
	}

	// Add user to the room's WebSocket connections
	wsConnections.Lock()
	if _, ok := wsConnections.m[roomId]; !ok {
		wsConnections.m[roomId] = make(map[string]*websocket.Conn)
	}
	wsConnections.m[roomId][userId] = ws
	wsConnections.Unlock()

	if debug {
		log.Printf("User %s joined room %s via WS\n", userId, roomId)
	}
	response := map[string]string{"message": "Joined room successfully"}
	if err := ws.WriteJSON(response); err != nil {
		log.Println("Error sending join response:", err)
	}
}

func handleDataMessage(ws *websocket.Conn, payload map[string]any) {
	from, ok1 := payload["from"].(string)
	data, ok2 := payload["data"] //  "data", not "message"

	if !ok1 || !ok2 {
		log.Println("Invalid data-message payload:", payload)
		return
	}

	// Get room ID from the session ID.  Crucially, this now uses *session* ID.
	roomId := getRoomIdBySessionId(from)
	if roomId == "" {
		log.Printf("Room not found for session: %s\n", from)
		return
	}

	// Broadcast to all participants in the room except the sender
	broadcastToRoom(roomId, gin.H{
		"type": "data-message",
		"payload": gin.H{
			"from": from,
			"data": data, //  "data", not "message"
		},
	}, from) // Exclude the sender
}

func handleWSDisconnect(ws *websocket.Conn, roomId string, userId string) {
	wsConnections.Lock()
	defer wsConnections.Unlock()
	//remove user from room
	if _, ok := wsConnections.m[roomId]; ok {
		delete(wsConnections.m[roomId], userId)
		if debug {
			log.Printf("User %s disconnected from room %s\n", userId, roomId)
		}
	}
}

func broadcastToRoom(roomId string, message gin.H, excludeUserId string) {
	if debug {
		log.Printf("Broadcasting to room %s: %+v (excluding: %s)\n", roomId, message, excludeUserId)
	}
	rooms.RLock()
	_, ok := rooms.m[roomId]
	rooms.RUnlock()
	if !ok {
		log.Printf("Room %s not found for broadcast\n", roomId)
		return
	}

	wsConnections.RLock()
	connections, ok := wsConnections.m[roomId]
	wsConnections.RUnlock()
	if !ok {
		log.Printf("No WebSocket connections found for room %s\n", roomId)
		return
	}

	// Serialize message once
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error serializing broadcast message: %v\n", err)
		return
	}

	if debug {
		log.Printf("Serialized message for broadcast: %s\n", string(msgBytes))
	}

	wsConnections.RLock()
	defer wsConnections.RUnlock()

	for userId, conn := range connections {
		if userId == excludeUserId {
			continue
		}
		if conn != nil {
			err := conn.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				log.Printf("Error broadcasting to user %s: %v\n", userId, err)
				// Clean up failed connection
				wsConnections.Lock()
				delete(wsConnections.m[roomId], userId)
				wsConnections.Unlock()
			} else if debug {
				log.Printf("Successfully sent broadcast message to user: %s\n", userId)
			}
		}
	}
}

func main() {
	initConfig() // Initialize configuration

	r := gin.Default()

	// CORS middleware (configure as needed)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // For development only!
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Upgrade, Connection")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")

		// Handle WebSocket pre-flight
		if c.Request.Method == "OPTIONS" {
			if c.Request.Header.Get("Upgrade") == "websocket" {
				c.Writer.Header().Set("Connection", "Upgrade")
				c.Writer.Header().Set("Upgrade", "websocket")
			}
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/auth/token", getTokenAPI)
	// Protected routes (require JWT)
	// api := r.Group("/api", verifyToken)
	// {
	// 	// Match Express routes
	// 	api.POST("/rooms", createRoomAPI)
	// 	api.GET("/rooms", getRoomsAPI)
	// 	api.GET("/inspect-rooms", inspectRoomsAPI)
	// 	api.POST("/rooms/:roomId/join", joinRoomAPI)
	// 	api.POST("/rooms/:roomId/sessions/:sessionId/publish", publishTracksAPI)
	// 	api.POST("/rooms/:roomId/sessions/:sessionId/unpublish", unpublishTrackAPI)
	// 	api.POST("/rooms/:roomId/sessions/:sessionId/pull", pullTracksAPI)
	// 	api.PUT("/rooms/:roomId/sessions/:sessionId/renegotiate", renegotiateSessionAPI)
	// 	api.POST("/rooms/:roomId/sessions/:sessionId/datachannels/new", manageDataChannelsAPI)
	// 	api.GET("/rooms/:roomId/participants", getParticipantsAPI)
	// 	api.GET("/rooms/:roomId/participant/:sessionId/tracks", getParticipantTracksAPI)
	// 	api.GET("/ice-servers", getICEServersAPI)
	// 	api.GET("/rooms/:roomId/sessions/:sessionId/state", getSessionStateAPI)
	// 	api.GET("/users/:userId", getUserInfoAPI)
	// 	api.POST("/rooms/:roomId/leave", leaveRoomAPI)
	// 	api.POST("/rooms/:roomId/sessions/:sessionId/track-status", updateTrackStatusAPI)
	// 	api.PUT("/rooms/:roomId/metadata", updateRoomMetadataAPI)
	// }

	r.GET("/ws", websocketHandler)

	// Start the server
	log.Printf("Server listening on http://localhost:%s\n", port)

	//listen and serve for http
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
