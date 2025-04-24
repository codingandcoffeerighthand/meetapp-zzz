import { createPeerConnection } from "@/utils/createPeerConnection"
import { create } from "zustand"
import { subscribeWithSelector } from "zustand/middleware"
import { useWebSocketStore } from "../ws"

const room = create(
    subscribeWithSelector(
        (set, get) => ({
            peerConnection: null,
            localTracks: [],
            err: null,
            isLoading: false,
            remoteConnection: null,
            remoteSession: "",
            remoteViews: [],
            onRemoteTrackCallback: (stream) => {
                console.log("default remote track callback", stream)
            },
            handleRequireRenegoiate: async (sessionId, sdp) => {
                set({ remoteSession: sessionId })
                set({ isLoading: true })
                let { remoteConnection, createRemoteConnection } = get()
                const { ws } = useWebSocketStore.getState()
                try {
                    if (!remoteConnection) {
                        await createRemoteConnection()
                    }
                    remoteConnection = get().remoteConnection
                    remoteConnection.setRemoteDescription(
                        new RTCSessionDescription({
                            sdp: sdp,
                            type: "offer"
                        })
                    )
                    const sdpAnswer = await remoteConnection.createAnswer()
                    remoteConnection.setLocalDescription(sdpAnswer)

                    ws.send(JSON.stringify({
                        message_name: "client_renegotiate_session",
                        data: {
                            session_id: sessionId,
                            sdp_answer: sdpAnswer.sdp
                        }
                    }))

                } catch (err) {
                    console.error(err)
                    set({ err })
                } finally {
                    set({ isLoading: false })
                }

            },

            createRemoteConnection: async () => {
                const { onRemoteTrackCallback } = get()
                try {
                    const p = await createPeerConnection()
                    if (p) {
                        p.ontrack = (evt) => {
                            console.info(evt)
                            set({ remoteViews: [...get().remoteViews, evt.track] })
                        }
                        set({ remoteConnection: p })
                    } else {
                        throw new Error("Peer connection not created")
                    }

                } catch (err) {
                    console.error('Error creating peer connection:', err)
                    set({ err: err })
                }
            },

            createPeerConnection: async () => {
                try {
                    const peerConnection = await createPeerConnection()
                    if (peerConnection) {
                        peerConnection.onconnectionstatechange = (evt) => {
                        }
                        set({ peerConnection })
                    } else {
                        throw new Error("Peer connection not created")
                    }

                } catch (err) {
                    console.error('Error creating peer connection:', err)
                    set({ err: err })
                }
            },
            joinRoom: async ({
                contract,
                account
            }, { stream, roomId = "a", participantNam = "aaa" }) => {
                set({ isLoading: true })
                if (!contract || !account) {
                    throw new Error("Web3 not initialized or no accounts found")
                }
                const { peerConnection, createPeerConnection } = get()
                if (!peerConnection) {
                    await createPeerConnection()
                }
                const { peerConnection: peerConn } = get()

                try {
                    set({ localTracks: stream.getTracks() })
                    const tranceivers = stream.getTracks().map(track => peerConn.addTransceiver(track, { direction: 'sendonly' }))

                    const offer = await peerConn.createOffer()
                    await peerConn.setLocalDescription(offer)
                    const offerStr = btoa(offer?.sdp)

                    const tracks = tranceivers.map(({ mid, sender }) => ([
                        sender.track?.id, mid, "local", true, "", roomId
                    ]))
                    await contract.methods.joinRoom(roomId, participantNam, tracks, offerStr).send({ from: account })
                } catch (err) {
                    console.error('Error starting stream:\n', err)
                    set({ err: err })
                } finally {
                    set({ isLoading: false })
                }
            },
            joinedHandler: async (data) => {
                const { peerConnection } = get()
                set({ isLoading: true })
                try {
                    const connected = new Promise((res, rej) => {
                        // timeout after 5s
                        setTimeout(rej, 15000);
                        const iceConnectionStateChangeHandler = () => {
                            if (peerConnection?.iceConnectionState === "connected") {
                                peerConnection?.removeEventListener(
                                    "iceconnectionstatechange",
                                    iceConnectionStateChangeHandler,
                                );
                                res(undefined);
                            }
                        };
                        peerConnection.addEventListener(
                            "iceconnectionstatechange",
                            iceConnectionStateChangeHandler,
                        );
                    });
                    peerConnection.setRemoteDescription(
                        new RTCSessionDescription({ sdp: data.sdp_answer, type: "answer" }),
                    );
                    await connected
                } catch (err) {
                    console.error('Error starting stream:\n', err)
                    set({ err: err })
                } finally {
                    set({ isLoading: false })
                }
            },
            hanlderRemoteVideo: async (eventStream, callback) => {
                callback(eventStream)
            },
            stop: () => {
                set({
                    localTracks: []
                })
            },
        })))
export { handleWs, room }