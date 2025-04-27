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
        midsOfLocalStream: {},
        addLocalTrack: async (stream, roomId) => {
            set({ isLoading: true })
            try {
                const { localPeerConnection, contract, account, localStreams } = get()
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
                const offer = await localPeerConnection.createOffer()
                await localPeerConnection.setLocalDescription(offer)
                const offerStr = btoa(offer?.sdp)
                const localStreamNumber = Object.keys(localStreams).length
                localStreams[localStreamNumber] = stream
                set({ localStreams })
                const mids = get().midsOfLocalStream
                transceivers.forEach(({ mid }) => {
                    if (!mids[localStreamNumber]) {
                        mids[localStreamNumber] = []
                    }
                    mids[localStreamNumber].push(mid)
                })

                set({ midsOfLocalStream: mids })

                const tracks = transceivers.map(({ mid, sender }) => ([
                    sender?.track?.id, mid, localStreamNumber, "local", true, "", roomId
                ]))

                await contract.methods.addTrack(roomId, tracks, offerStr).send({ from: account })
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
            evtSub.on("data", async data => {
                // const json = JSON.parse(data)
                const json = JSON.parse(Web3.utils.hexToUtf8(data.returnValues.eventData))
                console.info(data)
                console.info(json)
                switch (json?.event_name) {
                    case "joined_room":
                        await setSPDAnswerLocalPeerConnection(json.sdp_answer, data.returnValues.roomId)
                        break
                    case "pull_track":
                        await pullTrack(json.sdp_offer, json.remote_session, data.returnValues.roomId)
                        break
                    case "local_peer_connection_suscess":
                        break
                    case "new_participant_joined":
                        break
                    case "renegoiate_success":
                        const { handleRenegoiateSuccess } = get()
                        await handleRenegoiateSuccess(data.returnValues.roomId, json)
                        break
                    default:
                        console.error("Unknown event", json.event_name)
                }
            })
            evtSub.on("error", err => console.error(err))
            console.info(evtSub)
            set({ evtSub: evtSub })
        },
        handleRenegoiateSuccess: async (roomId, data) => {
            set({ isLoading: true })
            try {
                const { getRoomTracks } = get()
                await getRoomTracks(roomId)
                const { m, remoteTracks } = get()
                const tracks = data.tracks
                const n = {}
                const tnM = {}
                tracks.forEach(t => {
                    n[t.trackName] = remoteTracks[t.mid]
                    tnM[t.trackName] = t.mid
                })
                // debugger
                Object.entries(m).forEach(([k, v]) => {
                    v.trackNames.forEach(tn => {
                        if (n[tn]) {
                            m[k].stream.addTrack(
                                n[tn],
                            )
                        }
                        if (tnM[tn]) {
                            m[k].mids.push(tnM[tn])
                        }
                    })
                })
                console.info(m)

            } catch (err) {
                console.error(err)
            } finally {
                set({ isLoading: false })
            }

        },
        startListen: () => {
            const { contract, listenContractEvts } = get()
            set({ localPeerConnection: null, remotePeerConnection: null, localStreams: {}, remoteStreams: [] })
            if (contract) {
                listenContractEvts(contract)
            }
        },
        //  localPeerConnection
        localPeerConnection: null,
        localStreams: {},
        remoteTracks: {},
        resetLocal: () => {
            set({ localStreams: {}, m: {} })
        },
        startStream: async (roomId, participantName = "") => {
            set({ isLoading: true })
            const { contract, account, localStreams } = get()
            if (!participantName) {
                participantName = account
            }
            try {
                const localPeerConnection = await createPeerConnection()

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
                const localStreamNumber = Object.keys(localStreams).length
                const mids = get().midsOfLocalStream
                transceivers.forEach(({ mid }) => {
                    if (!mids[localStreamNumber]) {
                        mids[localStreamNumber] = []
                    }
                    mids[localStreamNumber].push(mid)
                })
                const tracks = transceivers.map(({ mid, sender }) => ([
                    sender?.track?.id, mid, localStreamNumber, "local", true, "", roomId
                ]))
                // debugger
                await contract.methods.joinRoom(roomId, participantName, tracks, offerStr).send({ from: account })
                localStreams[localStreamNumber] = stream
                set({ localPeerConnection, localStreams })
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
                console.info(localPeerConnection.connectionState)
                if (localPeerConnection.connectionState != "connected") {
                    await waitLocalConnection(localPeerConnection)
                }
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
            const { contract, account } = get()
            set({ isLoading: true })
            try {
                set({ remoteTracks: {} })
                const remotePeerConnection = await createPeerConnection()
                remotePeerConnection.ontrack = (data) => {
                    if (data.track) {
                        const { remoteStreams, remoteTracks } = get()
                        const newStream = new MediaStream()

                        newStream.addTrack(data.track)
                        // const tracks = [...remoteTracks, { track: data.track, mid: data?.transceiver?.mid }]
                        const tracks = remoteTracks
                        tracks[data?.transceiver?.mid] = data.track
                        var streams = []
                        if (remotePeerConnection) {
                            streams = [...remoteStreams, newStream]
                        }
                        set({ remoteStreams: streams, remoteTracks: tracks })
                    }
                }
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
                        sdp_answer: answerStr,
                        roomId: roomId,
                        addr: account,
                    }
                })
                await contract.methods.forwardEventToBackend(roomId, data).send({ from: account })

            } catch (err) {
                console.error(err)
            } finally {
                set({ isLoading: false })
            }
        },
        m: {},
        getRoomTracks: async (roomId) => {
            set({ isLoading: true })
            try {
                const { contract, account } = get()
                const data = await contract.methods.getParticipantOfRoom(roomId).call({
                    from: account
                })
                const ps = data[0]
                const tracks = data[1].filter(t => t?.isPublished)
                const m = {}
                tracks.forEach(t => {
                    if (t.isPublished) {
                        const st = t?.streamNumber
                        const s = t?.sessionId + "#" + st
                        const n = t?.trackName
                        if (!m[s]) {
                            m[s] = {
                                trackNames: [],
                                stream: new MediaStream(),
                                mids: []
                            }
                        }
                        m[s].trackNames.push(n)
                    }
                })
                ps.forEach((p) => {
                    const session = p?.sessionID
                    if (p?.walletAddress != account) {
                        Object.keys(m).forEach((key) => {
                            if (key.startsWith(session)) {
                                m[key].addr = p?.walletAddress
                                const st = key.split("#")[1]
                                m[key].name = p?.name + "#" + st
                            }
                        })
                    } else {
                        Object.keys(m).forEach((key) => {
                            if (key.startsWith(session)) {
                                delete m[key]
                            }
                        })
                    }
                })
                set({ m })
            } catch (err) {
                console.error('Error getting room tracks:', err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },
        closeStream: async (roomId, streamNum) => {
            set({ isLoading: true })
            try {
                const { contract, account, midsOfLocalStream, localStreams } = get()
                localStreams[streamNum].getTracks().forEach(track => track.stop())
                await contract.methods.removeTrack(roomId, midsOfLocalStream[streamNum], "").send({ from: account })
            }
            catch (err) {
                console.error(err)
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