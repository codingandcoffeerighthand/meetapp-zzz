'use client'
import { createPeerConnection } from "@/utils/createPeerConnection";
import { use, useEffect, useRef, useState } from 'react'

export function useLocalPeerConnection() {
    const localPeerConnection = useRef()
    const [localStreams, setLocalStreams] = useState([])
    const [err, setErr] = useState(null)
    const startLocalStream = async () => {
        if (localPeerConnection.current) {
            setErr(new Error("Local peer connection already exists"))
            return
        }
        try {
            localPeerConnection.current = await createPeerConnection()
            const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
            console.info(stream)
            setLocalStreams((prev) => [...prev, stream])
            return stream
        } catch (err) {
            console.error(err)
            setErr(err)
        }
    }

    const waitLocalConnection = () => new Promise((res, rej) => {
        // timeout after 5s
        setTimeout(rej, 5000);
        const iceConnectionStateChangeHandler = () => {
            if (localPeerConnection.current.iceConnectionState === "connected") {
                localPeerConnection.current.removeEventListener(
                    "iceconnectionstatechange",
                    iceConnectionStateChangeHandler,
                );
                res(undefined);
            }
        };
        localPeerConnection.current.addEventListener(
            "iceconnectionstatechange",
            iceConnectionStateChangeHandler,
        );
    });

    const setSdpAnswer = async (sdpAnswer) => {
        try {
            debugger
            await localPeerConnection.current.setRemoteDescription(
                new RTCSessionDescription({ sdp: sdpAnswer, type: "answer" }),
            );
            await waitLocalConnection()
        } catch (err) {
            console.error(err)
            setErr(err)
        }
    }



    return {
        localPeerConnection,
        localStreams,
        err,
        startLocalStream,
        // waitLocalConnection,
        setSdpAnswer,
    }
}