"use client";
import { useWeb3Store } from "@/store/web3";
import MediaStreamPlayer from "../video/video.component";
import { Button } from "../ui/button";
import { useEffect } from "react";

export default function RoomView({ roomId }) {
    const {
        isConnected, account, isLoading, localStreams,
        startListen, evtSub, addLocalTrack,
        startStream: startLocalStream, remoteStreams
    } = useWeb3Store()


    const startStream = async () => {
        startLocalStream(roomId)
    }

    const stopStream = () => {
    }
    const addTrackHanle = async () => {
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
        <Button onClick={addTrackHanle}>share screen</Button>
        <div className="w-[60%] flex flex-col gap-4 mx-6">
            {
                localStreams.map((stream, index) => {
                    return <div key={index}>
                        <MediaStreamPlayer
                            mediaStream={stream} title={`local#${index}`}
                            isLocal={true} />
                    </div>
                })
            }
            {remoteStreams.map((stream, idx) => (<MediaStreamPlayer key={idx} isLocal={false} mediaStream={stream} title="ppp" />))}
        </div>
    </>
}