// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package smc

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

// DAppMeetingParticipant is an auto generated low-level Go binding around an user-defined struct.
type DAppMeetingParticipant struct {
	WalletAddress common.Address
	Name          string
	SessionID     string
}

// DAppMeetingTrack is an auto generated low-level Go binding around an user-defined struct.
type DAppMeetingTrack struct {
	TrackName   string
	Mid         string
	Location    string
	IsPublished bool
	SessionId   string
	RoomId      string
}

// SmcMetaData contains all meta data concerning the Smc contract.
var SmcMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addAuthorized\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addAuthorizedBackend\",\"inputs\":[{\"name\":\"_backend\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addTrack\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_newTrack\",\"type\":\"tuple\",\"internalType\":\"structDAppMeeting.Track\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedBackends\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createRoom\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forceLeaveRoom\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_participant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forwardEventToBackend\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_eventData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forwardEventToFrontend\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_participant\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_eventData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getParticipantInfo\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDAppMeeting.Participant\",\"components\":[{\"name\":\"walletAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sessionID\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParticipantTracks\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_participant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structDAppMeeting.Track[]\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParticipantTracksCount\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_participant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoomParticipantsCount\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"joinRoom\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_initialTracks\",\"type\":\"tuple[]\",\"internalType\":\"structDAppMeeting.Track[]\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"_sessionId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"leaveRoom\",\"inputs\":[{\"name\":\"_roomId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"participantIndices\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"participantTrackCount\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"participantTracks\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"participantsInRoom\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeAuthorized\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeAuthorizedBackend\",\"inputs\":[{\"name\":\"_backend\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rooms\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"creationTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EventForwardedToBackend\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"eventData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EventForwardedToFrontend\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"participant\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"eventData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ParticipantJoined\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"participant\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"initialTracks\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structDAppMeeting.Track[]\",\"components\":[{\"name\":\"trackName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"mid\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isPublished\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sessionId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"roomId\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ParticipantLeft\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"participant\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TrackAdded\",\"inputs\":[{\"name\":\"roomId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"participant\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"trackName\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false}]",
	ID:  "Smc",
}

// Smc is an auto generated Go binding around an Ethereum contract.
type Smc struct {
	abi abi.ABI
}

// NewSmc creates a new instance of Smc.
func NewSmc() *Smc {
	parsed, err := SmcMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Smc{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Smc) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddAuthorized is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x26b232b0.
//
// Solidity: function addAuthorized() returns()
func (smc *Smc) PackAddAuthorized() []byte {
	enc, err := smc.abi.Pack("addAuthorized")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAddAuthorizedBackend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2776d2e.
//
// Solidity: function addAuthorizedBackend(address _backend) returns()
func (smc *Smc) PackAddAuthorizedBackend(backend common.Address) []byte {
	enc, err := smc.abi.Pack("addAuthorizedBackend", backend)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAddTrack is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b4ddda8.
//
// Solidity: function addTrack(string _roomId, (string,string,string,bool,string,string) _newTrack) returns()
func (smc *Smc) PackAddTrack(roomId string, newTrack DAppMeetingTrack) []byte {
	enc, err := smc.abi.Pack("addTrack", roomId, newTrack)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthorizedBackends is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47100107.
//
// Solidity: function authorizedBackends(uint256 ) view returns(address)
func (smc *Smc) PackAuthorizedBackends(arg0 *big.Int) []byte {
	enc, err := smc.abi.Pack("authorizedBackends", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthorizedBackends is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x47100107.
//
// Solidity: function authorizedBackends(uint256 ) view returns(address)
func (smc *Smc) UnpackAuthorizedBackends(data []byte) (common.Address, error) {
	out, err := smc.abi.Unpack("authorizedBackends", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackCreateRoom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7306d2dd.
//
// Solidity: function createRoom(string _roomId) returns()
func (smc *Smc) PackCreateRoom(roomId string) []byte {
	enc, err := smc.abi.Pack("createRoom", roomId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackForceLeaveRoom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x209fc017.
//
// Solidity: function forceLeaveRoom(string _roomId, address _participant) returns()
func (smc *Smc) PackForceLeaveRoom(roomId string, participant common.Address) []byte {
	enc, err := smc.abi.Pack("forceLeaveRoom", roomId, participant)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackForwardEventToBackend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaaea8180.
//
// Solidity: function forwardEventToBackend(string _roomId, bytes _eventData) returns()
func (smc *Smc) PackForwardEventToBackend(roomId string, eventData []byte) []byte {
	enc, err := smc.abi.Pack("forwardEventToBackend", roomId, eventData)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackForwardEventToFrontend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3398cf6.
//
// Solidity: function forwardEventToFrontend(string _roomId, address _participant, bytes _eventData) returns()
func (smc *Smc) PackForwardEventToFrontend(roomId string, participant common.Address, eventData []byte) []byte {
	enc, err := smc.abi.Pack("forwardEventToFrontend", roomId, participant, eventData)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetParticipantInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe970cf2c.
//
// Solidity: function getParticipantInfo(string _roomId) view returns((address,string,string))
func (smc *Smc) PackGetParticipantInfo(roomId string) []byte {
	enc, err := smc.abi.Pack("getParticipantInfo", roomId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetParticipantInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe970cf2c.
//
// Solidity: function getParticipantInfo(string _roomId) view returns((address,string,string))
func (smc *Smc) UnpackGetParticipantInfo(data []byte) (DAppMeetingParticipant, error) {
	out, err := smc.abi.Unpack("getParticipantInfo", data)
	if err != nil {
		return *new(DAppMeetingParticipant), err
	}
	out0 := *abi.ConvertType(out[0], new(DAppMeetingParticipant)).(*DAppMeetingParticipant)
	return out0, err
}

// PackGetParticipantTracks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3241b95.
//
// Solidity: function getParticipantTracks(string _roomId, address _participant) view returns((string,string,string,bool,string,string)[])
func (smc *Smc) PackGetParticipantTracks(roomId string, participant common.Address) []byte {
	enc, err := smc.abi.Pack("getParticipantTracks", roomId, participant)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetParticipantTracks is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3241b95.
//
// Solidity: function getParticipantTracks(string _roomId, address _participant) view returns((string,string,string,bool,string,string)[])
func (smc *Smc) UnpackGetParticipantTracks(data []byte) ([]DAppMeetingTrack, error) {
	out, err := smc.abi.Unpack("getParticipantTracks", data)
	if err != nil {
		return *new([]DAppMeetingTrack), err
	}
	out0 := *abi.ConvertType(out[0], new([]DAppMeetingTrack)).(*[]DAppMeetingTrack)
	return out0, err
}

// PackGetParticipantTracksCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf7cc8724.
//
// Solidity: function getParticipantTracksCount(string _roomId, address _participant) view returns(uint256)
func (smc *Smc) PackGetParticipantTracksCount(roomId string, participant common.Address) []byte {
	enc, err := smc.abi.Pack("getParticipantTracksCount", roomId, participant)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetParticipantTracksCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf7cc8724.
//
// Solidity: function getParticipantTracksCount(string _roomId, address _participant) view returns(uint256)
func (smc *Smc) UnpackGetParticipantTracksCount(data []byte) (*big.Int, error) {
	out, err := smc.abi.Unpack("getParticipantTracksCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetRoomParticipantsCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f2cc59f.
//
// Solidity: function getRoomParticipantsCount(string _roomId) view returns(uint256)
func (smc *Smc) PackGetRoomParticipantsCount(roomId string) []byte {
	enc, err := smc.abi.Pack("getRoomParticipantsCount", roomId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetRoomParticipantsCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3f2cc59f.
//
// Solidity: function getRoomParticipantsCount(string _roomId) view returns(uint256)
func (smc *Smc) UnpackGetRoomParticipantsCount(data []byte) (*big.Int, error) {
	out, err := smc.abi.Unpack("getRoomParticipantsCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackJoinRoom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe82a5259.
//
// Solidity: function joinRoom(string _roomId, string _name, (string,string,string,bool,string,string)[] _initialTracks, string _sessionId) returns()
func (smc *Smc) PackJoinRoom(roomId string, name string, initialTracks []DAppMeetingTrack, sessionId string) []byte {
	enc, err := smc.abi.Pack("joinRoom", roomId, name, initialTracks, sessionId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackLeaveRoom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfab38543.
//
// Solidity: function leaveRoom(string _roomId) returns()
func (smc *Smc) PackLeaveRoom(roomId string) []byte {
	enc, err := smc.abi.Pack("leaveRoom", roomId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (smc *Smc) PackOwner() []byte {
	enc, err := smc.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (smc *Smc) UnpackOwner(data []byte) (common.Address, error) {
	out, err := smc.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackParticipantIndices is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8df58945.
//
// Solidity: function participantIndices(string , address ) view returns(uint256)
func (smc *Smc) PackParticipantIndices(arg0 string, arg1 common.Address) []byte {
	enc, err := smc.abi.Pack("participantIndices", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackParticipantIndices is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8df58945.
//
// Solidity: function participantIndices(string , address ) view returns(uint256)
func (smc *Smc) UnpackParticipantIndices(data []byte) (*big.Int, error) {
	out, err := smc.abi.Unpack("participantIndices", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackParticipantTrackCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0062b748.
//
// Solidity: function participantTrackCount(string , address ) view returns(uint256)
func (smc *Smc) PackParticipantTrackCount(arg0 string, arg1 common.Address) []byte {
	enc, err := smc.abi.Pack("participantTrackCount", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackParticipantTrackCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0062b748.
//
// Solidity: function participantTrackCount(string , address ) view returns(uint256)
func (smc *Smc) UnpackParticipantTrackCount(data []byte) (*big.Int, error) {
	out, err := smc.abi.Unpack("participantTrackCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackParticipantTracks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5b85facc.
//
// Solidity: function participantTracks(string , address , uint256 ) view returns(string trackName, string mid, string location, bool isPublished, string sessionId, string roomId)
func (smc *Smc) PackParticipantTracks(arg0 string, arg1 common.Address, arg2 *big.Int) []byte {
	enc, err := smc.abi.Pack("participantTracks", arg0, arg1, arg2)
	if err != nil {
		panic(err)
	}
	return enc
}

// ParticipantTracksOutput serves as a container for the return parameters of contract
// method ParticipantTracks.
type ParticipantTracksOutput struct {
	TrackName   string
	Mid         string
	Location    string
	IsPublished bool
	SessionId   string
	RoomId      string
}

// UnpackParticipantTracks is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5b85facc.
//
// Solidity: function participantTracks(string , address , uint256 ) view returns(string trackName, string mid, string location, bool isPublished, string sessionId, string roomId)
func (smc *Smc) UnpackParticipantTracks(data []byte) (ParticipantTracksOutput, error) {
	out, err := smc.abi.Unpack("participantTracks", data)
	outstruct := new(ParticipantTracksOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.TrackName = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Mid = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Location = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.IsPublished = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.SessionId = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.RoomId = *abi.ConvertType(out[5], new(string)).(*string)
	return *outstruct, err

}

// PackParticipantsInRoom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2aba9dda.
//
// Solidity: function participantsInRoom(string , address ) view returns(bool)
func (smc *Smc) PackParticipantsInRoom(arg0 string, arg1 common.Address) []byte {
	enc, err := smc.abi.Pack("participantsInRoom", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackParticipantsInRoom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2aba9dda.
//
// Solidity: function participantsInRoom(string , address ) view returns(bool)
func (smc *Smc) UnpackParticipantsInRoom(data []byte) (bool, error) {
	out, err := smc.abi.Unpack("participantsInRoom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackRemoveAuthorized is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x061e1fc9.
//
// Solidity: function removeAuthorized() returns()
func (smc *Smc) PackRemoveAuthorized() []byte {
	enc, err := smc.abi.Pack("removeAuthorized")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRemoveAuthorizedBackend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9ba78fd9.
//
// Solidity: function removeAuthorizedBackend(address _backend) returns()
func (smc *Smc) PackRemoveAuthorizedBackend(backend common.Address) []byte {
	enc, err := smc.abi.Pack("removeAuthorizedBackend", backend)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRooms is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbbbcc869.
//
// Solidity: function rooms(string ) view returns(string roomId, uint256 creationTime)
func (smc *Smc) PackRooms(arg0 string) []byte {
	enc, err := smc.abi.Pack("rooms", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// RoomsOutput serves as a container for the return parameters of contract
// method Rooms.
type RoomsOutput struct {
	RoomId       string
	CreationTime *big.Int
}

// UnpackRooms is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbbbcc869.
//
// Solidity: function rooms(string ) view returns(string roomId, uint256 creationTime)
func (smc *Smc) UnpackRooms(data []byte) (RoomsOutput, error) {
	out, err := smc.abi.Unpack("rooms", data)
	outstruct := new(RoomsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.RoomId = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.CreationTime = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// SmcEventForwardedToBackend represents a EventForwardedToBackend event raised by the Smc contract.
type SmcEventForwardedToBackend struct {
	RoomId    string
	Sender    common.Address
	EventData []byte
	Raw       *types.Log // Blockchain specific contextual infos
}

const SmcEventForwardedToBackendEventName = "EventForwardedToBackend"

// ContractEventName returns the user-defined event name.
func (SmcEventForwardedToBackend) ContractEventName() string {
	return SmcEventForwardedToBackendEventName
}

// UnpackEventForwardedToBackendEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EventForwardedToBackend(string roomId, address sender, bytes eventData)
func (smc *Smc) UnpackEventForwardedToBackendEvent(log *types.Log) (*SmcEventForwardedToBackend, error) {
	event := "EventForwardedToBackend"
	if log.Topics[0] != smc.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(SmcEventForwardedToBackend)
	if len(log.Data) > 0 {
		if err := smc.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range smc.abi.Events[event].Inputs {
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

// SmcEventForwardedToFrontend represents a EventForwardedToFrontend event raised by the Smc contract.
type SmcEventForwardedToFrontend struct {
	RoomId      string
	Participant common.Address
	EventData   []byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const SmcEventForwardedToFrontendEventName = "EventForwardedToFrontend"

// ContractEventName returns the user-defined event name.
func (SmcEventForwardedToFrontend) ContractEventName() string {
	return SmcEventForwardedToFrontendEventName
}

// UnpackEventForwardedToFrontendEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EventForwardedToFrontend(string roomId, address participant, bytes eventData)
func (smc *Smc) UnpackEventForwardedToFrontendEvent(log *types.Log) (*SmcEventForwardedToFrontend, error) {
	event := "EventForwardedToFrontend"
	if log.Topics[0] != smc.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(SmcEventForwardedToFrontend)
	if len(log.Data) > 0 {
		if err := smc.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range smc.abi.Events[event].Inputs {
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

// SmcParticipantJoined represents a ParticipantJoined event raised by the Smc contract.
type SmcParticipantJoined struct {
	RoomId        string
	Participant   common.Address
	InitialTracks []DAppMeetingTrack
	Raw           *types.Log // Blockchain specific contextual infos
}

const SmcParticipantJoinedEventName = "ParticipantJoined"

// ContractEventName returns the user-defined event name.
func (SmcParticipantJoined) ContractEventName() string {
	return SmcParticipantJoinedEventName
}

// UnpackParticipantJoinedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ParticipantJoined(string roomId, address participant, (string,string,string,bool,string,string)[] initialTracks)
func (smc *Smc) UnpackParticipantJoinedEvent(log *types.Log) (*SmcParticipantJoined, error) {
	event := "ParticipantJoined"
	if log.Topics[0] != smc.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(SmcParticipantJoined)
	if len(log.Data) > 0 {
		if err := smc.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range smc.abi.Events[event].Inputs {
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

// SmcParticipantLeft represents a ParticipantLeft event raised by the Smc contract.
type SmcParticipantLeft struct {
	RoomId      string
	Participant common.Address
	Raw         *types.Log // Blockchain specific contextual infos
}

const SmcParticipantLeftEventName = "ParticipantLeft"

// ContractEventName returns the user-defined event name.
func (SmcParticipantLeft) ContractEventName() string {
	return SmcParticipantLeftEventName
}

// UnpackParticipantLeftEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ParticipantLeft(string roomId, address participant)
func (smc *Smc) UnpackParticipantLeftEvent(log *types.Log) (*SmcParticipantLeft, error) {
	event := "ParticipantLeft"
	if log.Topics[0] != smc.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(SmcParticipantLeft)
	if len(log.Data) > 0 {
		if err := smc.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range smc.abi.Events[event].Inputs {
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

// SmcTrackAdded represents a TrackAdded event raised by the Smc contract.
type SmcTrackAdded struct {
	RoomId      string
	Participant common.Address
	TrackName   string
	Raw         *types.Log // Blockchain specific contextual infos
}

const SmcTrackAddedEventName = "TrackAdded"

// ContractEventName returns the user-defined event name.
func (SmcTrackAdded) ContractEventName() string {
	return SmcTrackAddedEventName
}

// UnpackTrackAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TrackAdded(string roomId, address participant, string trackName)
func (smc *Smc) UnpackTrackAddedEvent(log *types.Log) (*SmcTrackAdded, error) {
	event := "TrackAdded"
	if log.Topics[0] != smc.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(SmcTrackAdded)
	if len(log.Data) > 0 {
		if err := smc.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range smc.abi.Events[event].Inputs {
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
