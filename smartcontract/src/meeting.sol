// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import {StringCompareLib} from "./StringCompareLib.sol";

/**
 * @title DAppMeeting
 * @dev A smart contract for managing virtual meeting rooms, participants, and media tracks.
 *      It includes functionalities for creating rooms, managing participants, and handling
 *      media tracks associated with participants. The contract also supports event forwarding
 *      to backend and frontend systems, and includes access control mechanisms for authorized
 *      backends and the contract owner.
 */
contract DAppMeeting {
    using StringCompareLib for string;
    /*
        *******************************
        ***      Structures         ***
        *******************************
    */

    struct Track {
        string trackName; // Name of the track
        string mid; // Media identifier in WebRTC
        uint256 streamNumber; // Media identifier in WebRTC
        string location; // Track location
        bool isPublished; // Publication status
        string sessionId; // Session ID associated with this track
        string roomId; // Room ID
    }

    struct Participant {
        address walletAddress; // Participant's wallet address
        string name; // Participant's name
        string sessionID; // Session ID from Cloudflare Calls
    }

    struct Room {
        string roomId; // Room ID
        uint256 creationTime; // Creation timestamp
        Participant[] participants; // List of participants
    }

    /*
        *******************************
        ***      State Variables    ***
        *******************************
    */
    mapping(string => Room) public rooms; // Store rooms by roomId
    mapping(string => mapping(address => uint256)) public participantIndices; // Track participant indices in arrays
    mapping(string => mapping(address => bool)) public participantsInRoom; // Check if participant is in room

    // New track mapping: roomId -> participantAddress -> trackList
    mapping(string => mapping(address => Track[])) public participantTracks;
    // Track count mapping for easy querying
    mapping(string => mapping(address => uint256)) public participantTrackCount;

    address public owner;
    address[] public authorizedBackends;

    /*
        *******************************
        ***      Events             ***
        *******************************
    */
    event ParticipantJoined(string roomId, address participant, Track[] initialTracks, string sdpOffer);
    event ParticipantLeft(string roomId, address participant);
    event TrackAdded(string roomId, address participant, Track[] tracks, string sdpOffer);
    event RemoveTracks(string roomId, string sessionId, Track[] removedTracks, string sdpOffer);
    event EventForwardedToBackend(string roomId, address sender, bytes eventData);
    event EventForwardedToFrontend(string roomId, address indexed participant, bytes eventData);
    event RoomCreated(string roomId);
    event SetParticipantSessionID(string roomId, address participant, string sessionID);
    /*
        *******************************
        ***      Constructor        ***
        *******************************
    */
    // Constructor

    constructor() {
        owner = msg.sender;
        authorizedBackends.push(msg.sender); // Owner is authorized by default
    }
    /*
        *******************************
        ***     Modifiers           ***
        *******************************
    */

    // Modifiers
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    modifier onlyAuthorized() {
        bool isAuthorized = false;
        for (uint256 i = 0; i < authorizedBackends.length; i++) {
            if (msg.sender == authorizedBackends[i]) {
                isAuthorized = true;
                break;
            }
        }
        require(isAuthorized, "Not authorized");
        _;
    }

    modifier roomExists(string memory _roomId) {
        require(bytes(rooms[_roomId].roomId).length > 0, "Room does not exist");
        _;
    }

    modifier participantExists(string memory _roomId) {
        require(participantsInRoom[_roomId][msg.sender], "You are not in this room");
        _;
    }

    /*
        *******************************
        ***      Functions          ***
        *******************************
    */

    function addAuthorized() public {
        authorizedBackends.push(msg.sender);
    }

    function _checkAuthorized(address _address) private view returns (bool) {
        for (uint256 i = 0; i < authorizedBackends.length; i++) {
            if (_address == authorizedBackends[i]) {
                return true;
            }
        }
        return false;
    }

    function checkAuthorized() public view returns (bool) {
        return _checkAuthorized(msg.sender);
    }

    // Management functions
    function addAuthorizedBackend(address _backend) public onlyOwner {
        authorizedBackends.push(_backend);
    }

    // Core functions
    function createRoom(string memory _roomId) public onlyAuthorized {
        require(bytes(rooms[_roomId].roomId).length == 0, "Room already exists");

        // Create a new room
        rooms[_roomId].roomId = _roomId;
        rooms[_roomId].creationTime = block.timestamp;
        emit RoomCreated(_roomId);
    }

    function joinRoom(
        string memory _roomId,
        string memory _name,
        Track[] memory _initialTracks,
        string memory _sdp_offer
    ) public onlyAuthorized roomExists(_roomId) {
        require(!participantsInRoom[_roomId][msg.sender], "Already in room");

        // Add participant to room
        Participant memory newParticipant = Participant({walletAddress: msg.sender, name: _name, sessionID: ""});

        rooms[_roomId].participants.push(newParticipant);
        uint256 participantIndex = rooms[_roomId].participants.length - 1;
        participantIndices[_roomId][msg.sender] = participantIndex;
        participantsInRoom[_roomId][msg.sender] = true;

        // Add initial tracks
        for (uint256 i = 0; i < _initialTracks.length; i++) {
            participantTracks[_roomId][msg.sender].push(_initialTracks[i]);
            participantTrackCount[_roomId][msg.sender]++;
        }

        // Emit the event for backend
        emit ParticipantJoined(_roomId, msg.sender, _initialTracks, _sdp_offer);
    }

    function leaveRoom(string memory _roomId) public roomExists(_roomId) participantExists(_roomId) {
        uint256 participantIndex = participantIndices[_roomId][msg.sender];

        // Delete participant from the room
        // Move the last participant to the position of the leaving participant
        if (participantIndex < rooms[_roomId].participants.length - 1) {
            Participant memory lastParticipant = rooms[_roomId].participants[rooms[_roomId].participants.length - 1];
            rooms[_roomId].participants[participantIndex] = lastParticipant;
            participantIndices[_roomId][lastParticipant.walletAddress] = participantIndex;
        }

        // Remove the last element
        rooms[_roomId].participants.pop();

        // Update mappings
        participantsInRoom[_roomId][msg.sender] = false;
        delete participantIndices[_roomId][msg.sender];
        delete participantTracks[_roomId][msg.sender];
        delete participantTrackCount[_roomId][msg.sender];

        // Emit event
        emit ParticipantLeft(_roomId, msg.sender);
    }

    function setParticipantSessionID(string memory _roomId, address _participantAddress, string memory _sessionID)
        public
        roomExists(_roomId)
        onlyOwner
    {
        require(participantsInRoom[_roomId][_participantAddress], "Participant not in room");
        uint256 participantIndex = participantIndices[_roomId][_participantAddress];
        rooms[_roomId].participants[participantIndex].sessionID = _sessionID;

        for (uint256 i = 0; i < participantTrackCount[_roomId][_participantAddress]; i++) {
            participantTracks[_roomId][_participantAddress][i].sessionId = _sessionID;
        }
        emit SetParticipantSessionID(_roomId, _participantAddress, _sessionID);
    }

    function addTrack(string memory _roomId, Track[] memory _newTracks, string memory sdpOffer)
        public
        roomExists(_roomId)
        participantExists(_roomId)
    {
        uint256 participantIndex = participantIndices[_roomId][msg.sender];
        string memory session = rooms[_roomId].participants[participantIndex].sessionID;
        for (uint256 i = 0; i < _newTracks.length; i++) {
            Track memory _newTrack = _newTracks[i];
            _newTrack.sessionId = session;
            participantTracks[_roomId][msg.sender].push(_newTrack);
            participantTrackCount[_roomId][msg.sender]++;
        }
        emit TrackAdded(_roomId, msg.sender, participantTracks[_roomId][msg.sender], sdpOffer);
    }

    function removeTrack(string memory _roomId, string[] memory _mids, string memory sdpOffer)
        public
        roomExists(_roomId)
        participantExists(_roomId)
    {
        Track[] memory removedTracks = new Track[](_mids.length);
        uint256 participantIdx = participantIndices[_roomId][msg.sender];
        string memory session = rooms[_roomId].participants[participantIdx].sessionID;
        for (uint256 i = 0; i < participantTracks[_roomId][msg.sender].length; i++) {
            for (uint256 j = 0; j < _mids.length; j++) {
                if (_mids[j].safeCompare(participantTracks[_roomId][msg.sender][i].mid)) {
                    participantTracks[_roomId][msg.sender][i].isPublished = false;
                    removedTracks[j] = participantTracks[_roomId][msg.sender][i];
                }
            }
        }
        emit RemoveTracks(_roomId, session, removedTracks, sdpOffer);
    }

    function getParticipantInfo(string memory _roomId)
        public
        view
        roomExists(_roomId)
        participantExists(_roomId)
        returns (Participant memory)
    {
        uint256 participantIndex = participantIndices[_roomId][msg.sender];
        return rooms[_roomId].participants[participantIndex];
    }

    function forwardEventToBackend(string memory _roomId, bytes memory _eventData)
        public
        roomExists(_roomId)
        participantExists(_roomId)
    {
        emit EventForwardedToBackend(_roomId, msg.sender, _eventData);
    }

    function forwardEventToFrontend(string memory _roomId, address _participant, bytes memory _eventData)
        public
        roomExists(_roomId)
        onlyOwner
    {
        require(participantsInRoom[_roomId][_participant], "Target participant not in room");
        emit EventForwardedToFrontend(_roomId, _participant, _eventData);
    }

    // New function to get a participant's tracks
    function getParticipantTracks(string memory _roomId, address _participant)
        public
        view
        roomExists(_roomId)
        returns (Track[] memory)
    {
        require(participantsInRoom[_roomId][_participant], "Participant not in room");
        return participantTracks[_roomId][_participant];
    }

    // Helper functions
    function getRoomParticipantsCount(string memory _roomId) public view roomExists(_roomId) returns (uint256) {
        return rooms[_roomId].participants.length;
    }

    function getParticipantTracksCount(string memory _roomId, address _participant)
        public
        view
        roomExists(_roomId)
        returns (uint256)
    {
        require(participantsInRoom[_roomId][_participant], "Participant not in room");
        return participantTrackCount[_roomId][_participant];
    }

    function getParticipantOfRoom(string memory _roomId)
        public
        view
        roomExists(_roomId)
        returns (Participant[] memory, Track[] memory)
    {
        uint256 participantCount = getRoomParticipantsCount(_roomId);
        uint256 trackSize = 0;
        for (uint256 i = 0; i < participantCount; i++) {
            trackSize += getParticipantTracksCount(_roomId, rooms[_roomId].participants[i].walletAddress);
        }
        Track[] memory tracks = new Track[](trackSize);
        uint256 idx = 0;
        for (uint256 i = 0; i < participantCount; i++) {
            for (
                uint256 j = 0; j < getParticipantTracksCount(_roomId, rooms[_roomId].participants[i].walletAddress); j++
            ) {
                tracks[idx++] = participantTracks[_roomId][rooms[_roomId].participants[i].walletAddress][j];
            }
        }
        return (rooms[_roomId].participants, tracks);
    }
}
