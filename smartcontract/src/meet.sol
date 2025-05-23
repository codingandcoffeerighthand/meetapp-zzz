// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.29;
import {StringCompareLib} from "./StringCompareLib.sol";

contract Meet {
    using StringCompareLib for string;

    /*
     *******************************
     ***      Structures         ***
     *******************************
     */

    struct Room {
        string roomId;
        string name;
        address creator;
        Participant[] participants;
    }

    struct Participant {
        address walletAddress;
        string name;
        string sessionID;
        Track[] tracks;
    }

    struct Track {
        string trackName;
        string mid;
        uint streamNumber;
        string location;
        bool isPublished;
        string sessionId;
        string roomId;
    }

    /*
     *******************************
     ***      State Variables    ***
     *******************************
     */

    mapping(string => Room) public rooms;
    address[] public authorizedBackends;
    address public owner;
    string private iceServers;

    /*
     *******************************
     ***      Events             ***
     *******************************
     */

    event RoomCreatedEvent(string roomId, string name, address creator);
    event JoinRoomEvent(
        string roomId,
        string sessionId,
        Track[] tracks,
        string sdpOffer
    );
    event AddTracksEvent(
        string roomId,
        string sessionId,
        Track[] tracks,
        string sdpOffer
    );
    event RemoveTracksEvent(
        string roomId,
        string sessionId,
        string[] mids,
        string sdpOffer
    );
    event LeftRoomEvent(string roomId, string sessionId);
    event NewSessionEvent(
        string indexed seesionHash,
        string roomId,
        string oldSessionId,
        string newSessionId
    );
    event BackendEvent(
        string roomId,
        string sessionId,
        string eventType,
        bytes data
    );
    event FrontendEvent(
        string indexed seesionHash,
        string roomId,
        string sessionId,
        string eventType,
        bytes data
    );

    constructor() {
        owner = msg.sender;
        authorizedBackends.push(msg.sender);
    }

    /*
     *******************************
     ***      Modifiers          ***
     *******************************
     */

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    modifier onlyAuthorized() {
        bool isAuthorizedBool = false;
        for (uint256 i = 0; i < authorizedBackends.length; i++) {
            if (authorizedBackends[i] == msg.sender) {
                isAuthorizedBool = true;
                break;
            }
        }
        require(isAuthorizedBool, "Not authorized");
        _;
    }

    modifier roomExists(string memory _roomId) {
        require(bytes(rooms[_roomId].roomId).length > 0, "Room does not exist");
        _;
    }

    modifier participantExists(
        string memory _roomId,
        string memory _sessionId
    ) {
        bool isParticipantExists = false;
        for (uint256 i = 0; i < rooms[_roomId].participants.length; i++) {
            if (
                _sessionId.safeCompare(rooms[_roomId].participants[i].sessionID)
            ) {
                isParticipantExists = true;
                break;
            }
        }
        require(isParticipantExists, "Participant does not exist");
        _;
    }

    /*
     *******************************
     ***      Functions          ***
     *******************************
     */
    function be_addAuthorized(address addr) public onlyOwner {
        _addAuthorized(addr);
    }

    function addAuthorized() public {
        _addAuthorized(msg.sender);
    }

    function setIceServers(string memory _iceServers) public onlyOwner {
        iceServers = _iceServers;
    }

    function createRoom(
        string memory _roomId,
        string memory _name
    ) public onlyAuthorized {
        require(bytes(rooms[_roomId].roomId).length == 0, "room already exist");

        rooms[_roomId].roomId = _roomId;
        rooms[_roomId].name = _name;
        rooms[_roomId].creator = msg.sender;
        emit RoomCreatedEvent(_roomId, _name, msg.sender);
    }

    function joinRoom(
        string memory _roomId,
        string memory _sessionLocal,
        string memory _participantName,
        Track[] memory _initialTracks,
        string memory _sdpOffer
    ) public onlyAuthorized roomExists(_roomId) {
        Participant memory newParticipant = Participant({
            walletAddress: msg.sender,
            name: _participantName,
            sessionID: _sessionLocal,
            tracks: new Track[](0)
        });

        rooms[_roomId].participants.push(newParticipant);
        for (uint256 i = 0; i < _initialTracks.length; i++) {
            rooms[_roomId]
                .participants[rooms[_roomId].participants.length - 1]
                .tracks
                .push(_initialTracks[i]);
        }
        emit JoinRoomEvent(_roomId, _sessionLocal, _initialTracks, _sdpOffer);
    }

    function newSession(
        string memory _roomId,
        string memory _oldSessionId,
        string memory _newSessionId
    )
        public
        onlyAuthorized
        roomExists(_roomId)
        participantExists(_roomId, _oldSessionId)
    {
        uint256 participantIndex = _getIndexOfParticipantBySessionId(
            _roomId,
            _oldSessionId
        );
        rooms[_roomId].participants[participantIndex].sessionID = _newSessionId;
        for (
            uint256 i = 0;
            i < rooms[_roomId].participants[participantIndex].tracks.length;
            i++
        ) {
            rooms[_roomId]
                .participants[participantIndex]
                .tracks[i]
                .sessionId = _newSessionId;
        }
        emit NewSessionEvent(
            _oldSessionId,
            _roomId,
            _oldSessionId,
            _newSessionId
        );
    }

    function addTracks(
        string memory _roomId,
        string memory _sessionId,
        Track[] memory _tracks,
        string memory _sdpOffer
    )
        public
        onlyAuthorized
        roomExists(_roomId)
        participantExists(_roomId, _sessionId)
    {
        uint256 participantIndex = _getIndexOfParticipantBySessionId(
            _roomId,
            _sessionId
        );
        for (uint256 i = 0; i < _tracks.length; i++) {
            rooms[_roomId].participants[participantIndex].tracks.push(
                _tracks[i]
            );
            rooms[_roomId]
                .participants[participantIndex]
                .tracks[i]
                .sessionId = _sessionId;
        }
        emit AddTracksEvent(_roomId, _sessionId, _tracks, _sdpOffer);
    }

    function removeTracks(
        string memory _roomId,
        string memory _sessionId,
        string[] memory _mids
    )
        public
        onlyAuthorized
        roomExists(_roomId)
        participantExists(_roomId, _sessionId)
    {
        for (uint256 i = 0; i < _mids.length; i++) {
            uint256 index = _findIndexOfTrackByMid(
                _roomId,
                _sessionId,
                _mids[i]
            );
            _deleteTracks(_roomId, _sessionId, index);
        }
        emit RemoveTracksEvent(_roomId, _sessionId, _mids, "");
    }

    function leaveRoom(
        string memory _roomId,
        string memory _sessionId
    )
        public
        onlyAuthorized
        roomExists(_roomId)
        participantExists(_roomId, _sessionId)
    {
        uint256 participantIndex = _getIndexOfParticipantBySessionId(
            _roomId,
            _sessionId
        );
        rooms[_roomId].participants[participantIndex] = rooms[_roomId]
            .participants[rooms[_roomId].participants.length - 1];
        rooms[_roomId].participants.pop();
        emit LeftRoomEvent(_roomId, _sessionId);
    }

    /*
     *******************************
     ***      Private functions  ***
     *******************************
     */
    function _findIndexOfTrackByMid(
        string memory _roomId,
        string memory _sessionId,
        string memory _mid
    ) private view returns (uint256) {
        uint256 participantIndex = _getIndexOfParticipantBySessionId(
            _roomId,
            _sessionId
        );
        for (
            uint256 i = 0;
            i < rooms[_roomId].participants[participantIndex].tracks.length;
            i++
        ) {
            if (
                _mid.safeCompare(
                    rooms[_roomId].participants[participantIndex].tracks[i].mid
                )
            ) {
                return i;
            }
        }
        return 0;
    }

    function _deleteTracks(
        string memory _roomId,
        string memory _sessionId,
        uint256 _index
    )
        public
        onlyAuthorized
        roomExists(_roomId)
        participantExists(_roomId, _sessionId)
    {
        uint256 participantIndex = _getIndexOfParticipantBySessionId(
            _roomId,
            _sessionId
        );
        rooms[_roomId].participants[participantIndex].tracks[_index] = rooms[
            _roomId
        ].participants[participantIndex].tracks[
                rooms[_roomId].participants[participantIndex].tracks.length - 1
            ];
        rooms[_roomId].participants[participantIndex].tracks.pop();
    }

    function _addAuthorized(address addr) private {
        authorizedBackends.push(addr);
    }

    function _getIndexOfParticipantBySessionId(
        string memory _roomId,
        string memory _sessionId
    ) private view returns (uint256) {
        for (uint256 i = 0; i < rooms[_roomId].participants.length; i++) {
            if (
                _sessionId.safeCompare(rooms[_roomId].participants[i].sessionID)
            ) {
                return i;
            }
        }
        return 0;
    }

    /*
     *******************************
     ***      Emit events        ***
     *******************************
     */

    function emitEventToBackend(
        string memory _roomId,
        string memory _sessionId,
        string memory _eventType,
        bytes memory _data
    ) public onlyAuthorized {
        emit BackendEvent(_roomId, _sessionId, _eventType, _data);
    }

    function emitEventToFrontend(
        string memory _roomId,
        string memory _sessionId,
        string memory _eventType,
        bytes memory _data
    ) public onlyOwner {
        emit FrontendEvent(_sessionId, _roomId, _sessionId, _eventType, _data);
    }

    /*
     *******************************
     ***     Get info functions  ***
     *******************************
     */
    function getParticipantInfoBySessionId(
        string memory _roomId,
        string memory _sessionId
    )
        public
        view
        roomExists(_roomId)
        participantExists(_roomId, _sessionId)
        returns (Participant memory)
    {
        uint256 participantIndex = _getIndexOfParticipantBySessionId(
            _roomId,
            _sessionId
        );
        return rooms[_roomId].participants[participantIndex];
    }

    function getRoomInfo(
        string memory _roomId
    ) public view roomExists(_roomId) returns (Room memory) {
        return rooms[_roomId];
    }

    function isAuthorized(address addr) public view returns (bool) {
        for (uint256 i = 0; i < authorizedBackends.length; i++) {
            if (authorizedBackends[i] == addr) {
                return true;
            }
        }
        return false;
    }

    function getIceServers()
        public
        view
        onlyAuthorized
        returns (string memory)
    {
        return iceServers;
    }
}
