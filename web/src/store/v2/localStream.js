import { createPeerConnection } from '@/utils/createPeerConnection'
import { v4 } from 'uuid'
import { create } from 'zustand'
import { subscribeWithSelector } from 'zustand/middleware'
import { waitLocalConnection } from '@/utils/awaitConnected'
import Web3 from 'web3'

const useLocalStream = create(
    subscribeWithSelector((set, get) => ({
        isLoading: false,
        error: null,
        localPeerConnection: null,
        /*
            localStreams: [
                {
                    stream: MediaStream,
                    mids: [mid]
                }
            ]
        */
        localStreams: [],
        sessionLocal: "",
        roomId: "",
        createLocalPeerConnection: async () => {
            const localPeerConnection = await createPeerConnection()
            set({ localPeerConnection })
        },
        setLocalStream: async (roomId) => {
            set({ isLoading: true })
            try {
                set({ roomId })
                const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
                await get().createLocalPeerConnection()
                const localPeerConnection = get().localPeerConnection
                const transceivers = stream.getTracks().map(
                    track => localPeerConnection.addTransceiver(track, {
                        direction: "sendonly"
                    })
                )
                const offer = await localPeerConnection?.createOffer()
                await localPeerConnection.setLocalDescription(offer)
                const sessionLocal = v4()
                set({ sessionLocal })

                const localStreams = get().localStreams
                const localStreamNumber = localStreams.length
                const mids = []
                transceivers.forEach(({ mid }) => {
                    mids.push(mid)
                })
                localStreams[localStreamNumber] = {
                    stream,
                    mids
                }
                const sdpString = btoa(offer?.sdp)
                const tracks = transceivers.map(({ mid, sender }) => ([
                    sender?.track?.id, mid, localStreamNumber, "local", true, "", roomId
                ]))
                set({ localStreams, localStreamNumber, localPeerConnection })
                return {
                    sessionLocal,
                    tracks,
                    sdpString
                }
            } catch (error) {
                console.error(error)
            } finally {
                set({ isLoading: false })
            }
        },
        localConnected: async (sdpAnswer) => {
            const localPeerConnection = get().localPeerConnection
            await localPeerConnection.setRemoteDescription(
                new RTCSessionDescription({ sdp: sdpAnswer, type: "answer" }),
            );
            if (localPeerConnection.connectionState != "connected") {
                await waitLocalConnection(localPeerConnection)
            }
            set({ localPeerConnection })
        },


        // end

    }))
)
export { useLocalStream }