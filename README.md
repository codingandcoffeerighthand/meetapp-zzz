# Meet App

A simple decentralized meeting application using smart contracts on EVM.

---

## Features

- Create Room
- Join Room
- Add/Remove Tracks (screen, mute mic, ...)
- Exit Room

---

## Architecture

- **Smart Contract**: Manages rooms, participants, and tracks on-chain.
- **Proxy/Backend**: Handles Cloudflare Realtime session management and signaling, interacts with the smart contract.
- **Client**: Sends track data and SDP offers to the smart contract, receives SDP answers and updates.

---

## API

> See `/docs/room.md` and backend OpenAPI spec for detailed API design and flows.

---

## Data Model

### Room

- `roomId`: string
- `name`: string

### User

- `walletAddress`: address
- `name`: string

### Participant

- `walletAddress`: address
- `name`: string
- `sessionId`: string

### Track

- `sessionId`: string
- `mid`: string
- `trackName`: string

---

## Flows

### 1. Publish Track

1. Client sends track data and SDP offer to the smart contract.
2. Proxy receives data, creates a new Cloudflare Realtime session, and updates the session info on the smart contract.
3. Proxy adds local track to the session via Cloudflare API.
4. Proxy returns SDP answer to the client.

### 2. Pull Track

1. Proxy fetches track info in the room from the smart contract.
2. Proxy creates a new session for the remote participant and adds tracks to this session.
3. Proxy sends SDP offer to the client.
4. Client sets remote description and returns SDP answer for renegotiation.

---

## Development

- Smart contract: Solidity, Foundry
- Backend: Go (see `/server`)
- Frontend: Next.js (see `/web`)

---

## Getting Started

See each subdirectory (`/smartcontract`, `/server`, `/web`) for build and run instructions.

---