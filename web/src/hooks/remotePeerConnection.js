import { createPeerConnection } from "@/utils/createPeerConnection";
import { useRef, useState } from "react";

export function useRemotePeerConnection() {
    const remotePeerConnection = useRef()
    const [remoteStreams, setRemoteStreams] = useState([])
    const [err, setErr] = useState(null)

    const initRemotePeerConnection = async (sdpOffer) => {
        if (remotePeerConnection.current) {
            setErr(new Error("Remote peer connection already exists"))
            return
        }
        try {
            remotePeerConnection.current = await createPeerConnection()
            await remotePeerConnection.current.setRemoteDescription(new RTCSessionDescription(JSON.stringify({
                type: 'offer',
                sdp: sdpOffer
            })));
            const sdpAnswer = await remotePeerConnection.current.createAnswer();

            remotePeerConnection.current.onTrack = (event) => {
                console.info("onTrack", event)
            };

            return sdpAnswer.sdp
        } catch (err) {
            console.error(err)
            setErr(err)
        }
    }
    return {
        remotePeerConnection,
        err,
        remoteStreams,
        initRemotePeerConnection,
    }
}