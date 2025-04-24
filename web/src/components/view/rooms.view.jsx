"use client";
import { useWeb3Store } from "@/store/web3";
import { useWebSocketStore } from "@/store/ws";
import VideoComponent from "../video/video.component";
import { Button } from "../ui/button";
import { room } from "@/store/features/wsHandle";
import { v4 } from "uuid";
import { useEffect, useState } from "react";

export default function RoomView({ roomId }) {
    const { isConnected, account, contract } = useWeb3Store()
    const { ws, handleWs } = useWebSocketStore()
    const { isLoading, joinRoom, localTracks, stop, createPeerConnection, remoteViews } = room()
    handleWs(roomId)

    const getRemoteStreams = (streams) => {
        const stream = new MediaStream()
        remoteViews.forEach(view => {
            stream.addTrack(view)
        })
        streams.push(stream)
    }
    const streams = []
    getRemoteStreams(streams)




    const startStream = async () => {
        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        await joinRoom({
            account, contract
        }, { stream, roomId: roomId, participantNam: v4() })
    }
    const stopStream = () => {
        localTracks.forEach(track => track.stop())
        stop()
    }

    return <>
        {isLoading && <p className="text-red-500">Loading...</p>}
        <p>web e: {isConnected} account: {account}</p>
        <p>websocket: {ws?.readyState}</p>
        <Button onClick={startStream}>start</Button>
        <Button onClick={stopStream}>stop</Button>
        <div className="w-[60%] flex flex-col gap-4 mx-6">
            <div className="w-[400px]">
                <VideoComponent isLocal={true} tracks={localTracks} />
            </div>
            {streams.map((stream, index) => {
                return <div key={index}>
                    <VideoComponent isLocal={false} tracks={stream.getTracks()} />
                </div>
            })}
        </div>
    </>
}