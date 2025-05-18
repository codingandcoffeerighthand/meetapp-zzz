// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package meet_smc

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// MeetParticipant is an auto generated low-level Go binding around an user-defined struct.
type MeetParticipant struct {
	WalletAddress common.Address
	Name          string
	SessionID     string
	Tracks        []MeetTrack
}

// MeetRoom is an auto generated low-level Go binding around an user-defined struct.
type MeetRoom struct {
	RoomId       string
	Name         string
	Creator      common.Address
	Participants []MeetParticipant
}

// MeetTrack is an auto generated low-level Go binding around an user-defined struct.
type MeetTrack struct {
	TrackName    string
	Mid          string
	StreamNumber *big.Int
	Location     string
	IsPublished  bool
	SessionId    string
	RoomId       string
}

// MeetMetaData contains all meta data concerning the Meet contract.
var MeetMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_deleteTracks\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addAuthorized\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addTracks\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_tracks\",\"type\":\"tuple[]\",\"internalType\":\"structMeet.Track[]\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"streamNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"_sdpOffer\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedBackends\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"be_addAuthorized\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createRoom\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitEventToBackend\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_eventType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitEventToFrontend\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_eventType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getIceServers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParticipantInfoBySessionId\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sessionId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structMeet.Participant\",\"components\":[{\"name\":\"walletAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sessionID\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tracks\",\"type\":\"tuple[]\",\"internalType\":\"structMeet.Track[]\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"streamNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoomInfo\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structMeet.Room\",\"components\":[{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"participants\",\"type\":\"tuple[]\",\"internalType\":\"structMeet.Participant[]\",\"components\":[{\"name\":\"walletAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sessionID\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tracks\",\"type\":\"tuple[]\",\"internalType\":\"structMeet.Track[]\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"streamNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAuthorized\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"joinRoom\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sessionLocal\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_participantName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_initialTracks\",\"type\":\"tuple[]\",\"internalType\":\"structMeet.Track[]\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"streamNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"_sdpOffer\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"leaveRoom\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sessionId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"newSession\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_oldSessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_newSessionId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeTracks\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_mids\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rooms\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setIceServers\",\"inputs\":[{\"name\":\"_iceServers\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AddTracksEvent\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sessionId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"tracks\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structMeet.Track[]\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"streamNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"sdpOffer\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BackendEvent\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sessionId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"eventType\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FrontendEvent\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sessionId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"eventType\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"JoinRoomEvent\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sessionId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"tracks\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structMeet.Track[]\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"streamNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"sdpOffer\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LeftRoomEvent\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sessionId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewSessionEvent\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"oldSessionId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"newSessionId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoveTracksEvent\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sessionId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"mids\",\"type\":\"string[]\",\"indexed\":false,\"internalType\":\"string[]\"},{\"name\":\"sdpOffer\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoomCreatedEvent\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
	ID:  "Meet",
}

// Meet is an auto generated Go binding around an Ethereum contract.
type Meet struct {
	abi abi.ABI
}

// NewMeet creates a new instance of Meet.
func NewMeet() *Meet {
	parsed, err := MeetMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Meet{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Meet) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackDeleteTracks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b2a87be.
//
// Solidity: function _deleteTracks(string _roomId, string _sessionId, uint256 _index) returns()
func (meet *Meet) PackDeleteTracks(roomId string, sessionId string, index *big.Int) []byte {
	enc, err := meet.abi.Pack("_deleteTracks", roomId, sessionId, index)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAddAuthorized is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x26b232b0.
//
// Solidity: function addAuthorized() returns()
func (meet *Meet) PackAddAuthorized() []byte {
	enc, err := meet.abi.Pack("addAuthorized")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAddTracks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x96d7d8db.
//
// Solidity: function addTracks(string _roomId, string _sessionId, (string,string,uint256,string,bool,string,string)[] _tracks, string _sdpOffer) returns()
func (meet *Meet) PackAddTracks(roomId string, sessionId string, tracks []MeetTrack, sdpOffer string) []byte {
	enc, err := meet.abi.Pack("addTracks", roomId, sessionId, tracks, sdpOffer)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthorizedBackends is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47100107.
//
// Solidity: function authorizedBackends(uint256 ) view returns(address)
func (meet *Meet) PackAuthorizedBackends(arg0 *big.Int) []byte {
	enc, err := meet.abi.Pack("authorizedBackends", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthorizedBackends is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x47100107.
//
// Solidity: function authorizedBackends(uint256 ) view returns(address)
func (meet *Meet) UnpackAuthorizedBackends(data []byte) (common.Address, error) {
	out, err := meet.abi.Unpack("authorizedBackends", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBeAddAuthorized is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe37777cd.
//
// Solidity: function be_addAuthorized(address addr) returns()
func (meet *Meet) PackBeAddAuthorized(addr common.Address) []byte {
	enc, err := meet.abi.Pack("be_addAuthorized", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCreateRoom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfa68daad.
//
// Solidity: function createRoom(string _roomId, string _name) returns()
func (meet *Meet) PackCreateRoom(roomId string, name string) []byte {
	enc, err := meet.abi.Pack("createRoom", roomId, name)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEmitEventToBackend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xeb40db26.
//
// Solidity: function emitEventToBackend(string _roomId, string _sessionId, string _eventType, bytes _data) returns()
func (meet *Meet) PackEmitEventToBackend(roomId string, sessionId string, eventType string, data []byte) []byte {
	enc, err := meet.abi.Pack("emitEventToBackend", roomId, sessionId, eventType, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEmitEventToFrontend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x92e35653.
//
// Solidity: function emitEventToFrontend(string _roomId, string _sessionId, string _eventType, bytes _data) returns()
func (meet *Meet) PackEmitEventToFrontend(roomId string, sessionId string, eventType string, data []byte) []byte {
	enc, err := meet.abi.Pack("emitEventToFrontend", roomId, sessionId, eventType, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetIceServers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb8c3eea1.
//
// Solidity: function getIceServers() view returns(string)
func (meet *Meet) PackGetIceServers() []byte {
	enc, err := meet.abi.Pack("getIceServers")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetIceServers is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb8c3eea1.
//
// Solidity: function getIceServers() view returns(string)
func (meet *Meet) UnpackGetIceServers(data []byte) (string, error) {
	out, err := meet.abi.Unpack("getIceServers", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackGetParticipantInfoBySessionId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xae3f25e3.
//
// Solidity: function getParticipantInfoBySessionId(string _roomId, string _sessionId) view returns((address,string,string,(string,string,uint256,string,bool,string,string)[]))
func (meet *Meet) PackGetParticipantInfoBySessionId(roomId string, sessionId string) []byte {
	enc, err := meet.abi.Pack("getParticipantInfoBySessionId", roomId, sessionId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetParticipantInfoBySessionId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xae3f25e3.
//
// Solidity: function getParticipantInfoBySessionId(string _roomId, string _sessionId) view returns((address,string,string,(string,string,uint256,string,bool,string,string)[]))
func (meet *Meet) UnpackGetParticipantInfoBySessionId(data []byte) (MeetParticipant, error) {
	out, err := meet.abi.Unpack("getParticipantInfoBySessionId", data)
	if err != nil {
		return *new(MeetParticipant), err
	}
	out0 := *abi.ConvertType(out[0], new(MeetParticipant)).(*MeetParticipant)
	return out0, err
}

// PackGetRoomInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3d600e7c.
//
// Solidity: function getRoomInfo(string _roomId) view returns((string,string,address,(address,string,string,(string,string,uint256,string,bool,string,string)[])[]))
func (meet *Meet) PackGetRoomInfo(roomId string) []byte {
	enc, err := meet.abi.Pack("getRoomInfo", roomId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetRoomInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3d600e7c.
//
// Solidity: function getRoomInfo(string _roomId) view returns((string,string,address,(address,string,string,(string,string,uint256,string,bool,string,string)[])[]))
func (meet *Meet) UnpackGetRoomInfo(data []byte) (MeetRoom, error) {
	out, err := meet.abi.Unpack("getRoomInfo", data)
	if err != nil {
		return *new(MeetRoom), err
	}
	out0 := *abi.ConvertType(out[0], new(MeetRoom)).(*MeetRoom)
	return out0, err
}

// PackIsAuthorized is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe9fbb80.
//
// Solidity: function isAuthorized(address addr) view returns(bool)
func (meet *Meet) PackIsAuthorized(addr common.Address) []byte {
	enc, err := meet.abi.Pack("isAuthorized", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsAuthorized is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfe9fbb80.
//
// Solidity: function isAuthorized(address addr) view returns(bool)
func (meet *Meet) UnpackIsAuthorized(data []byte) (bool, error) {
	out, err := meet.abi.Unpack("isAuthorized", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackJoinRoom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x763b80a1.
//
// Solidity: function joinRoom(string _roomId, string _sessionLocal, string _participantName, (string,string,uint256,string,bool,string,string)[] _initialTracks, string _sdpOffer) returns()
func (meet *Meet) PackJoinRoom(roomId string, sessionLocal string, participantName string, initialTracks []MeetTrack, sdpOffer string) []byte {
	enc, err := meet.abi.Pack("joinRoom", roomId, sessionLocal, participantName, initialTracks, sdpOffer)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackLeaveRoom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x420dcd45.
//
// Solidity: function leaveRoom(string _roomId, string _sessionId) returns()
func (meet *Meet) PackLeaveRoom(roomId string, sessionId string) []byte {
	enc, err := meet.abi.Pack("leaveRoom", roomId, sessionId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackNewSession is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97ca04b0.
//
// Solidity: function newSession(string _roomId, string _oldSessionId, string _newSessionId) returns()
func (meet *Meet) PackNewSession(roomId string, oldSessionId string, newSessionId string) []byte {
	enc, err := meet.abi.Pack("newSession", roomId, oldSessionId, newSessionId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (meet *Meet) PackOwner() []byte {
	enc, err := meet.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (meet *Meet) UnpackOwner(data []byte) (common.Address, error) {
	out, err := meet.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackRemoveTracks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6baa1c74.
//
// Solidity: function removeTracks(string _roomId, string _sessionId, string[] _mids) returns()
func (meet *Meet) PackRemoveTracks(roomId string, sessionId string, mids []string) []byte {
	enc, err := meet.abi.Pack("removeTracks", roomId, sessionId, mids)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRooms is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbbbcc869.
//
// Solidity: function rooms(string ) view returns(string roomId, string name, address creator)
func (meet *Meet) PackRooms(arg0 string) []byte {
	enc, err := meet.abi.Pack("rooms", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// RoomsOutput serves as a container for the return parameters of contract
// method Rooms.
type RoomsOutput struct {
	RoomId  string
	Name    string
	Creator common.Address
}

// UnpackRooms is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbbbcc869.
//
// Solidity: function rooms(string ) view returns(string roomId, string name, address creator)
func (meet *Meet) UnpackRooms(data []byte) (RoomsOutput, error) {
	out, err := meet.abi.Unpack("rooms", data)
	outstruct := new(RoomsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.RoomId = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Creator = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	return *outstruct, err

}

// PackSetIceServers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe850e914.
//
// Solidity: function setIceServers(string _iceServers) returns()
func (meet *Meet) PackSetIceServers(iceServers string) []byte {
	enc, err := meet.abi.Pack("setIceServers", iceServers)
	if err != nil {
		panic(err)
	}
	return enc
}

// MeetAddTracksEvent represents a AddTracksEvent event raised by the Meet contract.
type MeetAddTracksEvent struct {
	RoomId    common.Hash
	SessionId common.Hash
	Tracks    []MeetTrack
	SdpOffer  string
	Raw       *types.Log // Blockchain specific contextual infos
}

const MeetAddTracksEventEventName = "AddTracksEvent"

// ContractEventName returns the user-defined event name.
func (MeetAddTracksEvent) ContractEventName() string {
	return MeetAddTracksEventEventName
}

// UnpackAddTracksEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AddTracksEvent(string indexed roomId, string indexed sessionId, (string,string,uint256,string,bool,string,string)[] tracks, string sdpOffer)
func (meet *Meet) UnpackAddTracksEventEvent(log *types.Log) (*MeetAddTracksEvent, error) {
	event := "AddTracksEvent"
	if log.Topics[0] != meet.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MeetAddTracksEvent)
	if len(log.Data) > 0 {
		if err := meet.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range meet.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MeetBackendEvent represents a BackendEvent event raised by the Meet contract.
type MeetBackendEvent struct {
	RoomId    common.Hash
	SessionId common.Hash
	EventType string
	Data      []byte
	Raw       *types.Log // Blockchain specific contextual infos
}

const MeetBackendEventEventName = "BackendEvent"

// ContractEventName returns the user-defined event name.
func (MeetBackendEvent) ContractEventName() string {
	return MeetBackendEventEventName
}

// UnpackBackendEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BackendEvent(string indexed roomId, string indexed sessionId, string eventType, bytes data)
func (meet *Meet) UnpackBackendEventEvent(log *types.Log) (*MeetBackendEvent, error) {
	event := "BackendEvent"
	if log.Topics[0] != meet.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MeetBackendEvent)
	if len(log.Data) > 0 {
		if err := meet.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range meet.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MeetFrontendEvent represents a FrontendEvent event raised by the Meet contract.
type MeetFrontendEvent struct {
	RoomId    common.Hash
	SessionId common.Hash
	EventType string
	Data      []byte
	Raw       *types.Log // Blockchain specific contextual infos
}

const MeetFrontendEventEventName = "FrontendEvent"

// ContractEventName returns the user-defined event name.
func (MeetFrontendEvent) ContractEventName() string {
	return MeetFrontendEventEventName
}

// UnpackFrontendEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FrontendEvent(string indexed roomId, string indexed sessionId, string eventType, bytes data)
func (meet *Meet) UnpackFrontendEventEvent(log *types.Log) (*MeetFrontendEvent, error) {
	event := "FrontendEvent"
	if log.Topics[0] != meet.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MeetFrontendEvent)
	if len(log.Data) > 0 {
		if err := meet.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range meet.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MeetJoinRoomEvent represents a JoinRoomEvent event raised by the Meet contract.
type MeetJoinRoomEvent struct {
	RoomId    common.Hash
	SessionId common.Hash
	Tracks    []MeetTrack
	SdpOffer  string
	Raw       *types.Log // Blockchain specific contextual infos
}

const MeetJoinRoomEventEventName = "JoinRoomEvent"

// ContractEventName returns the user-defined event name.
func (MeetJoinRoomEvent) ContractEventName() string {
	return MeetJoinRoomEventEventName
}

// UnpackJoinRoomEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event JoinRoomEvent(string indexed roomId, string indexed sessionId, (string,string,uint256,string,bool,string,string)[] tracks, string sdpOffer)
func (meet *Meet) UnpackJoinRoomEventEvent(log *types.Log) (*MeetJoinRoomEvent, error) {
	event := "JoinRoomEvent"
	if log.Topics[0] != meet.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MeetJoinRoomEvent)
	if len(log.Data) > 0 {
		if err := meet.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range meet.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MeetLeftRoomEvent represents a LeftRoomEvent event raised by the Meet contract.
type MeetLeftRoomEvent struct {
	RoomId    common.Hash
	SessionId common.Hash
	Raw       *types.Log // Blockchain specific contextual infos
}

const MeetLeftRoomEventEventName = "LeftRoomEvent"

// ContractEventName returns the user-defined event name.
func (MeetLeftRoomEvent) ContractEventName() string {
	return MeetLeftRoomEventEventName
}

// UnpackLeftRoomEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LeftRoomEvent(string indexed roomId, string indexed sessionId)
func (meet *Meet) UnpackLeftRoomEventEvent(log *types.Log) (*MeetLeftRoomEvent, error) {
	event := "LeftRoomEvent"
	if log.Topics[0] != meet.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MeetLeftRoomEvent)
	if len(log.Data) > 0 {
		if err := meet.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range meet.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MeetNewSessionEvent represents a NewSessionEvent event raised by the Meet contract.
type MeetNewSessionEvent struct {
	RoomId       common.Hash
	OldSessionId common.Hash
	NewSessionId string
	Raw          *types.Log // Blockchain specific contextual infos
}

const MeetNewSessionEventEventName = "NewSessionEvent"

// ContractEventName returns the user-defined event name.
func (MeetNewSessionEvent) ContractEventName() string {
	return MeetNewSessionEventEventName
}

// UnpackNewSessionEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event NewSessionEvent(string indexed roomId, string indexed oldSessionId, string newSessionId)
func (meet *Meet) UnpackNewSessionEventEvent(log *types.Log) (*MeetNewSessionEvent, error) {
	event := "NewSessionEvent"
	if log.Topics[0] != meet.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MeetNewSessionEvent)
	if len(log.Data) > 0 {
		if err := meet.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range meet.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MeetRemoveTracksEvent represents a RemoveTracksEvent event raised by the Meet contract.
type MeetRemoveTracksEvent struct {
	RoomId    common.Hash
	SessionId common.Hash
	Mids      []string
	SdpOffer  string
	Raw       *types.Log // Blockchain specific contextual infos
}

const MeetRemoveTracksEventEventName = "RemoveTracksEvent"

// ContractEventName returns the user-defined event name.
func (MeetRemoveTracksEvent) ContractEventName() string {
	return MeetRemoveTracksEventEventName
}

// UnpackRemoveTracksEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RemoveTracksEvent(string indexed roomId, string indexed sessionId, string[] mids, string sdpOffer)
func (meet *Meet) UnpackRemoveTracksEventEvent(log *types.Log) (*MeetRemoveTracksEvent, error) {
	event := "RemoveTracksEvent"
	if log.Topics[0] != meet.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MeetRemoveTracksEvent)
	if len(log.Data) > 0 {
		if err := meet.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range meet.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MeetRoomCreatedEvent represents a RoomCreatedEvent event raised by the Meet contract.
type MeetRoomCreatedEvent struct {
	RoomId  common.Hash
	Name    string
	Creator common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const MeetRoomCreatedEventEventName = "RoomCreatedEvent"

// ContractEventName returns the user-defined event name.
func (MeetRoomCreatedEvent) ContractEventName() string {
	return MeetRoomCreatedEventEventName
}

// UnpackRoomCreatedEventEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoomCreatedEvent(string indexed roomId, string name, address indexed creator)
func (meet *Meet) UnpackRoomCreatedEventEvent(log *types.Log) (*MeetRoomCreatedEvent, error) {
	event := "RoomCreatedEvent"
	if log.Topics[0] != meet.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(MeetRoomCreatedEvent)
	if len(log.Data) > 0 {
		if err := meet.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range meet.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}
