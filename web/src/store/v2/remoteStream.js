import { waitLocalConnection } from '@/utils/awaitConnected'
import { createPeerConnection } from '@/utils/createPeerConnection'
import { create } from 'zustand'
import { subscribeWithSelector } from 'zustand/middleware'

const useRemoteStream = create(
    subscribeWithSelector((set, get) => ({
        isLoading: false,
        error: null,
        remoteStreams: [],
        remoteTracks: {},
        roomId: "",
        remotePeerConnection: null,
        createRemotePeerConnection: async () => {
            const remotePeerConnection = await createPeerConnection()
            remotePeerConnection.ontrack = (event) => {
                if (event.track) {
                    const newStream = new MediaStream()
                    newStream.addTrack(event.track)
                    const { remoteTracks } = get()
                    remoteTracks[event?.transceiver?.mid] = event.track
                    set({
                        remoteStreams: [...get().remoteStreams, newStream],
                        remoteTracks: remoteTracks
                    })
                }
            }
            set({ remotePeerConnection })
        },
        setRemoteStream: async (sdpOffer) => {
            set({ isLoading: true })
            try {
                set({ remoteStreams: [], remoteTracks: [] })
                await get().createRemotePeerConnection()
                const remotePeerConnection = get().remotePeerConnection
                await remotePeerConnection.setRemoteDescription(
                    new RTCSessionDescription({ sdp: sdpOffer, type: "offer" }),
                );
                const answer = await remotePeerConnection.createAnswer()
                await remotePeerConnection.setLocalDescription(answer)
                return answer.sdp
            } catch (error) {
                set({ error })
            } finally {
                set({ isLoading: false })
            }
        },
        remoteConnected: async (sdpAnswer) => {
            const remotePeerConnection = get().remotePeerConnection
            await remotePeerConnection.setRemoteDescription(
                new RTCSessionDescription({ sdp: sdpAnswer, type: "answer" }),
            );
        },
        mapRemoteStreams: {},
        setMapRemoteStreams: (participants, tracks, localSession) => {
            const { remoteTracks } = get()
            const mapRemoteStreams = {}
            participants.filter(
                (p) => p?.session_id !== localSession
            ).forEach((participant) => {
                mapRemoteStreams[participant?.session_id] = {
                    name: participant?.name,
                    address: participant?.wallet_address,
                    midsOfStream: [],
                    streams: []
                }
                let mids = {}
                participant?.tracks?.forEach((pt) => {
                    tracks.filter(
                        (t) => t?.track_name === pt?.track_name
                    ).forEach((t) => {
                        if (!mids[pt?.stream_number]) {
                            mids[pt?.stream_number] = []
                        }
                        mids[pt?.stream_number].push(t?.mid)
                    })
                })
                console.info('mids', mids)
                mapRemoteStreams[participant?.session_id].midsOfStream = mids
                Object.values(mids).forEach((mid) => {
                    const newStream = new MediaStream()
                    mid.forEach((m) => {
                        newStream.addTrack(remoteTracks[m])
                    })
                    mapRemoteStreams[participant?.session_id].streams.push(newStream)
                })
            })
            set({ mapRemoteStreams })
        },
        handlerNoRemoteTracks: () => {
            set({ mapRemoteStreams: {}, remoteTracks: {}, remoteStreams: [] })
        }
    }))
)

export default useRemoteStream
