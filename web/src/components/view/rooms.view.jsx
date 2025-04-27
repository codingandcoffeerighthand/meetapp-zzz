"use client";
import { useWeb3Store } from "@/store/web3";
import MediaStreamPlayer from "../video/video.component";
import { Button } from "../ui/button";
import { useEffect } from "react";

export default function RoomView({ roomId }) {
    const {
        isConnected, account, isLoading, addLocalTrack, m,
        startStream: startLocalStream, localStreams, resetLocal,
        closeStream
    } = useWeb3Store()
    useEffect(() => {
        resetLocal()
        return resetLocal
    }, [])
    const startStream = async () => {
        startLocalStream(roomId)
    }

    const stopStream = async () => {
        await closeStream(roomId, 1)
    }
    const addTrackHandle = async () => {
        const t = await navigator.mediaDevices.getDisplayMedia({
            video: true,
        })
        await addLocalTrack(t, roomId)
    }

    return <>
        {isLoading && <p className="text-red-500">Loading...</p>}
        <p>web e: {isConnected} account: {account}</p>
        <Button onClick={startStream}>start</Button>
        <Button onClick={stopStream}>stop</Button>
        <Button onClick={addTrackHandle}>share screen</Button>
        <div className="w-[60%] flex flex-col gap-4 mx-6">
            {
                Object.entries(localStreams).map(([index, stream]) => {
                    return <div key={index}>
                        <MediaStreamPlayer
                            mediaStream={stream} title={`local#${index}`}
                            isLocal={true} />
                    </div>
                })
            }
            {/* {remoteStreams.map((stream, idx) => (<MediaStreamPlayer key={idx} isLocal={false} mediaStream={stream} title="ppp" />))} */}
            {Object.values(m).map(t => (
                <MediaStreamPlayer key={t.name} title={t.name} isLocal={false} mediaStream={t.stream} />
            ))}
        </div>
    </>
}