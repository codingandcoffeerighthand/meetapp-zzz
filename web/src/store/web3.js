"use client"
import { create } from 'zustand'
import Web3 from 'web3'
import abi from '@/abi/abi.json'
import { subscribeWithSelector } from 'zustand/middleware'
import { v4 } from 'uuid'
import { createPeerConnection } from '@/utils/createPeerConnection'
import { sleep, waitLocalConnection } from '@/utils/awaitConnected'
const useWeb3Store = create(
    subscribeWithSelector((set, get) => ({
        // web3
        web3: null,
        account: "",
        chainId: null,
        contract: null,
        isConnected: false,
        isLoading: false,
        error: null,
        isAuthorized: false,
        roomId: "",
        prK: "",
        evtHandlers: {},

        // Wallet and account
        initWeb3: async () => {
            if (typeof window !== "undefined") {
                set({ isLoading: true, error: null })
                try {
                    let web3Instance;
                    const provider = new Web3.providers.WebsocketProvider(process.env.NEXT_PUBLIC_INFURA_URL);
                    web3Instance = new Web3(provider);
                    const prK = localStorage?.getItem('prK')
                    if (!prK) {
                        set({ err: new Error("Private key not found") })
                        return
                    }
                    const acc = web3Instance.eth.accounts.wallet.add(prK)
                    set({ account: acc[0].address, isConnected: true, prK, web3: web3Instance })
                } catch (error) {
                    console.error("Error initializing Web3:", error)
                    set({ error: error.message })
                } finally {
                    set({ isLoading: false })
                }
            }
        },

        addKey: async (pk) => {
            set({ isLoading: true })
            try {
                localStorage.setItem('prK', pk)
                await get().initWeb3()
            } catch (err) {
                console.error(err)
                set({ err: err })
            } finally {
                set({ isLoading: false })
            }
        },

        initContract: async () => {
            set({ isLoading: true })
            try {
                const { web3, account, checkAuth, listenContractEvts } = get()
                if (!web3 || !account) {
                    set({ error: 'Web3 not initialized or no accounts found' })
                    return
                }
                const contractInstance = new web3.eth.Contract(abi, process.env.NEXT_PUBLIC_CONTRACT_ADDRESS)
                set({ contract: contractInstance })
                await checkAuth()
            } catch (err) {
                console.error('Error initializing contract:', err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },

        getOwner: async () => {
            set({ isLoading: true })
            try {
                const { contract } = get()
                const owner = await contract.methods.owner().call()
                return owner
            }
            catch (err) {
                console.error('Error getting owner:', err)
                set({ error: err })
            }
            finally {
                set({ isLoading: false })
            }
        },
        register: async () => {
            set({ isLoading: true })
            try {
                const { account, contract } = get()
                if (!contract) {
                    set({ error: 'contract not initialized' })
                    return
                }
                await contract.methods.addAuthorized().send({ from: account })
                set({ isAuthorized: true })
            }
            catch (err) {
                console.error('Error registering:', err)
                set({ error: err })
            }
            finally {
                set({ isLoading: false })
            }
        },
        checkAuth: async () => {
            set({ isLoading: true })
            try {
                const { account, contract } = get()
                if (!contract) {
                    set({ error: 'contract not initialized' })
                    return
                }
                const isAuth = await contract.methods.checkAuthorized().call({
                    from: account
                })
                set({ isAuthorized: isAuth })
            }
            catch (err) {
                console.error('Error checking auth:', err)
                set({ error: err })
            }
            finally {
                set({ isLoading: false })
            }
        },
        createRoom: async (roomId, callback) => {
            set({ isLoading: true })
            try {
                const { account, contract } = get()
                if (!contract) {
                    set({ error: 'contract not initialized' })
                    return
                }
                if (!roomId) {
                    roomId = v4()
                }
                await contract.methods.createRoom(roomId).send({ from: account })
                set({ roomId })
                callback(roomId)
            }
            catch (err) {
                console.error('Error creating room:', err)
                set({ error: err })
            }
            finally {
                set({ isLoading: false })
            }
        },
        addLocalTrack: async (stream, roomId) => {
            set({ isLoading: true })
            try {
                const { localPeerConnection, contract, account } = get()
                if (!localPeerConnection) {
                    throw new Error("Local peer connection not initialized")
                }
                const transceivers = stream.getTracks().map(
                    track => localPeerConnection?.addTransceiver(track, {
                        direction: "sendonly"
                    })
                )
                /*
                    format [trackId, mid, "local, true, "", roomid]
                */
                const tracks = transceivers.map(({ mid, sender }) => ([
                    sender?.track?.id, mid, "local", true, "", roomId
                ]))
                console.info(roomId, tracks[0])
                await contract.methods.addTrack(roomId, tracks[0]).send({ from: account })
                const data = Web3.utils.toHex({
                    event_name: "local_peer_connection_suscess"
                })
                await contract.methods.forwardEventToBackend(roomId, data).send({ from: account })
            }
            catch (err) { console.error(err) }
            finally {
                set({ isLoading: false })
            }
        },


        setEvtHandler: (eventName, handleFunction) => {
            set({
                evtHandlers: {
                    ...get().evtHandlers,
                    [eventName]: handleFunction
                }
            })
        },
        evtSub: null,
        listenContractEvts: (contract) => {
            const { setSPDAnswerLocalPeerConnection, pullTrack } = get()
            const { account: address, evtHandlers } = get()
            const evtSub = contract.events.EventForwardedToFrontend({
                filter: { participant: address },
            }, (err, evt) => { console.info(err, evt) })
            evtSub.on("data", data => {
                // const json = JSON.parse(data)
                const json = JSON.parse(Web3.utils.hexToUtf8(data.returnValues.eventData))
                console.info(data)
                console.info(json)
                switch (json?.event_name) {
                    case "joined_room":
                        setSPDAnswerLocalPeerConnection(json.sdp_answer, data.returnValues.roomId)
                        break
                    case "pull_track":
                        pullTrack(json.sdp_offer, json.remote_session, data.returnValues.roomId)
                        break
                    case "local_peer_connection_suscess":
                        break
                    case "new_participant_joined":
                        break
                    default:
                        console.error("Unknown event", json.event_name)
                }
            })
            evtSub.on("error", err => console.error(err))
            console.info(evtSub)
            set({ evtSub: evtSub })
        },
        startListen: () => {
            const { contract, listenContractEvts } = get()
            set({ localPeerConnection: null, remotePeerConnection: null, localStreams: [], remoteStreams: [] })
            if (contract) {
                listenContractEvts(contract)
            }
        },
        //  localPeerConnection
        localPeerConnection: null,
        localStreams: [],
        localStreamNumber: 0,
        startStream: async (roomId, participantName = "") => {
            set({ isLoading: true })
            const { contract, account, localStreams } = get()
            if (!participantName) {
                participantName = account
            }
            try {
                const localPeerConnection = await createPeerConnection()
                const remotePeerConnection = await createPeerConnection()
                remotePeerConnection.ontrack = (data) => {
                    if (data.track) {
                        const { remoteStreams } = get()
                        const newStream = new MediaStream()

                        newStream.addTrack(data.track)
                        var streams = []
                        if (remotePeerConnection) {
                            streams = [...remoteStreams, newStream]
                        }
                        console.info(data)
                        set({ remoteStreams: streams })
                    }
                }
                set({ remotePeerConnection })
                const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
                const transceivers = stream.getTracks().map(
                    track => localPeerConnection?.addTransceiver(track, {
                        direction: "sendonly"
                    })
                )
                const offer = await localPeerConnection?.createOffer()
                await localPeerConnection.setLocalDescription(offer)
                // base64
                const offerStr = btoa(offer?.sdp)

                /*
                    format [trackId, mid, "local, true, "", roomid]
                */
                const { localStreamNumber } = get()
                const tracks = transceivers.map(({ mid, sender }) => ([
                    sender?.track?.id, mid, localStreamNumber, "local", true, "", roomId
                ]))
                set({ localStreamNumber: localStreamNumber + 1 })
                // debugger
                await contract.methods.joinRoom(roomId, participantName, tracks, offerStr).send({ from: account })
                set({ localPeerConnection, localStreams: [...localStreams, stream] })
                return stream

            } catch (err) {
                console.error(err)
                set({ err })
            } finally {
                set({ isLoading: false })
            }
        },
        setSPDAnswerLocalPeerConnection: async (sdpAnswer, roomId = "") => {
            const { contract, localPeerConnection, account } = get()
            try {
                await localPeerConnection.setRemoteDescription(
                    new RTCSessionDescription({ sdp: sdpAnswer, type: "answer" }),
                );
                await waitLocalConnection(localPeerConnection)
                const data = Web3.utils.toHex({
                    event_name: "local_peer_connection_suscess"
                })
                await contract.methods.forwardEventToBackend(roomId, data).send({ from: account })
            } catch (err) {
                console.error(err)
            }
        },
        // remotePeerConnection
        remotePeerConnection: null,
        remoteStreams: [],
        remoteSession: "",
        pullTrack: async (sdp_offer, remoteSession, roomId) => {
            const { contract, remotePeerConnection, account } = get()
            set({ isLoading: true })
            try {
                await remotePeerConnection.setRemoteDescription(
                    new RTCSessionDescription({ sdp: sdp_offer, type: "offer" }),
                );
                const answer = await remotePeerConnection.createAnswer()
                await remotePeerConnection.setLocalDescription(answer)
                // base64
                const answerStr = btoa(answer?.sdp)
                const data = Web3.utils.toHex({
                    event_name: "renegoiate_session",
                    data: {
                        remote_session: remoteSession,
                        sdp_answer: answerStr
                    }
                })
                await contract.methods.forwardEventToBackend(roomId, data).send({ from: account })
                const { getRoomTracks } = get()
                await getRoomTracks(roomId)

            } catch (err) {
                console.error(err)
            } finally {
                set({ isLoading: false })
            }
        },

        getRoomTracks: async (roomId) => {
            set({ isLoading: true })
            try {
                const { contract, account } = get()
                const data = await contract.methods.getParticipantOfRoom(roomId).call({
                    from: account
                })
                const ps = data[0]
                const tracks = data[1]
                const mapParticipants = {}
                ps.forEach((p) => {
                    const session = p?.sessionID
                    if (session) {
                        mapParticipants[session] = {}
                        mapParticipants[session].name = p?.name
                        mapParticipants[session].address = p?.walletAddress
                    }
                })
                tracks.forEach(t => {
                    const s = t?.sessionId
                    const n = t?.trackName
                    if (s && n) {
                        if (mapParticipants[s]) {
                            mapParticipants[s][n] = t
                        }
                    }

                })
                console.info(mapParticipants)
            } catch (err) {
                console.error('Error getting room tracks:', err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        }
    })))


useWeb3Store.subscribe(
    (state) => state.account,
    (account) => {
        if (account) {
            useWeb3Store.getState().initContract()
        }
    }
)
useWeb3Store.subscribe(
    (state) => state.contract,
    (account) => {
        if (account) {
            useWeb3Store.getState().startListen()
        }
    }
)


await useWeb3Store.getState().initWeb3()

export {
    useWeb3Store
}