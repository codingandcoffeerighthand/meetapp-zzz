"use client"
import { create } from 'zustand'
import { subscribeWithSelector } from 'zustand/middleware'
import abi from '@/abi/abi.json'
import Web3 from 'web3'
import { useLocalStream } from './localStream'
import useRemoteStream from './remoteStream'
import { useCallback } from 'react'

const useWeb3V2Store = create(
    subscribeWithSelector((set, get) => ({
        isLoading: false,
        error: null,
        roomId: "",
        privateKey: "",
        account: "",
        isWeb3Connected: false,
        web3: null,
        contract: null,
        isAuthorized: false,
        localSession: "",
        remoteSession: "",

        cleanError: () => {
            set({ error: null })
        },

        addPrivateKey: (privateKey) => {
            set({ privateKey })
        },

        getPrivateKey: () => {
            return get().privateKey
        },

        connectWeb3: async () => {
            if (typeof window !== "undefined") {
                set({ isLoading: true, error: null })
                try {
                    const provider = new Web3.providers.WebsocketProvider(process.env.NEXT_PUBLIC_WS_URL);
                    const web3Instance = new Web3(provider);
                    const prK = get().privateKey
                    console.log(prK)
                    if (!prK) {
                        throw new Error("Private key not found")
                    }
                    const acc = web3Instance.eth.accounts.wallet.add(prK)
                    set({
                        account: acc[0].address,
                        isWeb3Connected: true,
                        privateKey: prK,
                        web3: web3Instance
                    })
                    const contractInstance = new web3Instance.eth.Contract(abi, process.env.NEXT_PUBLIC_CONTRACT_ADDRESS)
                    set({ contract: contractInstance })
                    get().callCheckAuthorized()
                } catch (error) {
                    set({ error: error })
                    console.error("Error initializing Web3:", error)
                } finally {
                    set({ isLoading: false })
                }
            }
        },
        _modifierContract: () => {
            const { contract, account, isWeb3Connected } = get()
            if (!contract) {
                throw new Error("Contract not initialized")
            }
            if (!isWeb3Connected) {
                throw new Error("Web3 not connected")
            }
            if (!account) {
                throw new Error("Account not connected")
            }
        },
        callCheckAuthorized: async () => {
            set({ isLoading: true, error: null })
            try {
                get()._modifierContract()
                const { contract, account } = get()
                const isAuthorized = await contract.methods.isAuthorized(account).call({ from: account })
                set({ isAuthorized })
            } catch (err) {
                console.error("Error checking auth:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        callRegister: async () => {
            set({ isLoading: true, error: null })
            try {
                get()._modifierContract()
                const { contract, account } = get()
                await contract.methods.addAuthorized().send({ from: account })
                get().callCheckAuthorized()
            } catch (err) {
                console.error("Error registering:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        callCreateRoom: async (roomId, name) => {
            set({ isLoading: true, error: null })
            try {
                if (!roomId) {
                    throw new Error("Room ID are required")
                }
                get()._modifierContract()
                const { contract, account } = get()
                await contract.methods.createRoom(roomId, name).send({ from: account })
                console.log("Room created \n", roomId, "\n", name)
            } catch (err) {
                console.error("Error creating room:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        callJoinRoom: async (roomId, localSessionId, participantName, tracks, offerStr) => {
            set({ roomId })
            set({ isLoading: true, error: null })
            try {
                if (!roomId || !localSessionId || !tracks || !offerStr) {
                    throw new Error("Room ID, localSessionId , tracks, offerStr are required")
                }
                get()._modifierContract()
                const { contract, account } = get()
                await contract.methods.joinRoom(roomId, localSessionId, participantName, tracks, offerStr).send({ from: account })

                set({ localSession: localSessionId })
                await get().subFrontendEvent()
                console.log("Room joined \n", localSessionId)
            } catch (err) {
                console.error("Error joining room:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        handleJoinRoom: async (roomId, participantName) => {
            set({ isLoading: true, error: null })
            try {
                const { sessionLocal, tracks, sdpString } = await useLocalStream.getState().setLocalStream(roomId)
                await get().callJoinRoom(roomId, sessionLocal, participantName, tracks, sdpString)
            } catch (err) {
                console.error("Error handling join room:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },

        subFrontendEvent: async () => {
            set({ isLoading: true, error: null })
            try {
                const { contract, localSession } = get()
                if (!contract) {
                    throw new Error("Contract not initialized")
                }
                if (!localSession) {
                    throw new Error("Local session not initialized")
                }
                const evtSub = contract.events.FrontendEvent({
                    filter: { seesionHash: localSession }
                }, (err, evt) => { console.info(err, evt) })

                evtSub.on("data", async evt => {
                    const data = evt?.returnValues
                    console.info("#### event", evt)
                    switch (data?.eventType) {
                        case EventType.JOINED_ROOM:
                            await get().handlerLocalConnected(data?.data)
                            await evtSub.unsubscribe()
                            evtSub.removeAllListeners()
                            await get().subFrontendEvent()
                            break
                        case EventType.RemoteConnect:
                            get().handlerRemoteConnect(data?.data)
                            break
                        case EventType.RemoteConnectSuccess:
                            await get().handlerRemoteConnectSuccess(data)
                            break
                        case EventType.NoRemoteTracks:
                            await useRemoteStream.getState().handlerNoRemoteTracks(data?.data)
                            break
                        case EventType.AddTracks:
                            await get().handlerAddTracks(data?.data)
                            break
                        default:
                            console.info("#### unknown event", evt)
                            break
                    }
                })
            } catch (err) {
                console.error("Error subscribing to frontend event:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },

        emitLocalConnected: async (roomId, data) => {
            set({ isLoading: true, error: null })
            try {
                get()._modifierContract()
                const { contract, account } = get()
                await contract.methods.emitEventToBackend(
                    roomId,
                    get().localSession,
                    EventType.LocalConnected,
                    data
                ).send({ from: account })
            } catch (err) {
                console.error("Error emitting local connected:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        handlerLocalConnected: async (data) => {
            set({ isLoading: true, error: null })
            try {
                const { roomId } = get()
                const json = JSON.parse(Web3.utils.hexToUtf8(data))
                const sdpAnswer = json?.sdp_answer
                const newSession = json?.new_session
                set({ localSession: newSession })
                console.info("new session", newSession)
                await useLocalStream.getState().localConnected(sdpAnswer)
                const rsdata = Web3.utils.toHex({
                    "room_id": roomId,
                })
                await get().emitLocalConnected(roomId, rsdata)
            } catch (err) {
                console.error("Error handling local connected:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },

        emitRemoteConnected: async (sdpAnswer) => {
            set({ isLoading: true, error: null })
            try {
                get()._modifierContract()
                const { contract, account, roomId, remoteSession } = get()
                const data = Web3.utils.toHex({
                    "room_id": roomId,
                    "remote_session_id": remoteSession,
                    "sdp_answer": sdpAnswer
                })
                await contract.methods.emitEventToBackend(
                    roomId,
                    get().localSession,
                    EventType.RemoteConnected,
                    data
                ).send({ from: account })
            } catch (err) {
                console.error("Error emitting remote connected:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },


        handlerRemoteConnect: async (data) => {
            set({ isLoading: true, error: null })
            try {
                const json = JSON.parse(Web3.utils.hexToUtf8(data))
                set({
                    remoteSession: json?.remote_session_id,
                })
                const sdpOffer = json?.sdp_offer
                const sdpAnswer = await useRemoteStream.getState().setRemoteStream(sdpOffer)
                await get().emitRemoteConnected(sdpAnswer)
            } catch (err) {
                console.error("Error handling remote connect:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        handlerRemoteConnectSuccess: async (data) => {
            set({ isLoading: true, error: null })
            try {
                console.info("#### remote connect success", data)
                const json = JSON.parse(Web3.utils.hexToUtf8(data?.data))
                console.info("json", json)
                useRemoteStream.getState().setMapRemoteStreams(json?.room?.participants, json?.tracks)
                console.info("map stream", useRemoteStream.getState().mapRemoteStreams)
                console.info("remote tracks", useRemoteStream.getState().remoteTracks)
            } catch (err) {
                console.error("Error handling remote connect success:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },

        callLeaveRoom: async (rommID, callback = () => { }) => {
            set({ isLoading: true, error: null })
            try {
                get()._modifierContract()
                const { contract, account, localSession } = get()
                await contract.methods.leaveRoom(rommID, localSession).send({ from: account })
                callback()
            } catch (err) {
                console.error("Error leaving room:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        addLocalTracks: async (stream, roomId) => {
            set({ isLoading: true })
            try {
                const { tracks, sdpString } = await useLocalStream.getState().addLocalTracks(stream, roomId, get().localSession)
                const { contract, account, localSession } = get()
                await contract.methods.addTracks(roomId, localSession, tracks, sdpString).send({ from: account })
            } catch (err) {
                console.error("Error adding local track:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        handlerAddTracks: async (data) => {
            set({ isLoading: true })
            try {
                const sdpAnswer = Web3.utils.hexToUtf8(data)
                await useLocalStream.getState().localConnected(sdpAnswer)
                const dt = Web3.utils.toHex({
                    "room_id": get().roomId,
                })
                await get().emitLocalConnected(get().roomId, dt)
            } catch (err) {
                console.error("Error handling add tracks:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        removeStream: async (streamId) => {
            set({ isLoading: true })
            try {
                console.log("remove stream", streamId)
                const { contract, account, localSession, roomId } = get()
                const mids = await useLocalStream.getState().removeStream(streamId)
                const tx = await contract.methods.removeTracks(roomId, localSession, mids).send({ from: account })
                console.info("tx", tx)
            } catch (err) {
                console.error("Error removing stream:", err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        }
    })))

const EventType = {
    JOINED_ROOM: "joined_room",
    LocalConnected: "local_connected",
    RemoteConnected: "remote_connected",
    RemoteConnect: "remote_connect",
    RemoteConnectSuccess: "remote_connect_success",
    NoRemoteTracks: "no-remote-tracks",
    AddTracks: "add_tracks",
}

export {
    EventType,
    useWeb3V2Store
};
