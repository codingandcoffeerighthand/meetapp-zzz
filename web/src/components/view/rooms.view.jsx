"use client";
import { useWeb3Store } from "@/store/web3";
import MediaStreamPlayer from "../video/video.component";
import { Button } from "../ui/button";
import { use, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Card } from "../ui/card";
import { Input } from "../ui/input";

export default function RoomView({ roomId }) {
    const {
        isConnected, account, isLoading, addLocalTrack, m,
        startStream: startLocalStream, localStreams, resetLocal,
        closeStream, leaveRoom
    } = useWeb3Store()
    useEffect(() => {
        resetLocal()
        return resetLocal
    }, [])
    const startStream = async () => {
        startLocalStream(roomId, name)
    }
    const [name, setName] = useState("")
    const stopStream = async (streamNum) => {
        await closeStream(roomId, streamNum)
    }
    const addTrackHandle = async () => {
        const t = await navigator.mediaDevices.getDisplayMedia({
            video: true,
        })
        await addLocalTrack(t, roomId)
    }
    const router = useRouter()
    const handleLeaveRoom = () => {

        leaveRoom(roomId, () => router.push("/"))
    }

    return <div>
        <div className="flex justify-between">
            <div className="p-4">
                {isLoading && <p className="text-red-500">Loading...</p>}
                <p>account: {account}</p>
            </div>
            <div className="m-auto flex gap-4 justify-center">
                {/* <Input className="inline-block" placeholder="name" type="text" value={name}
                    onChange={(e) => setName(e.target.value)}
                /> */}
                <Button className="inline-inline" onClick={startStream}>start</Button>
                <Button onClick={addTrackHandle}>share screen</Button>
                <Button onClick={handleLeaveRoom}>exit room</Button>
            </div>
        </div>
        <div className="w-full flex flex-wrap justify-center gap-4 mx-6">
            {
                Object.entries(localStreams).map(([index, stream]) => {
                    return <div key={index} className="w-[300px]">
                        <MediaStreamPlayer
                            mediaStream={stream} title={`local#${index}`}
                            isLocal={true}
                            closeVideoCallback={() => stopStream(index)}
                        />
                    </div>
                })
            }
            {/* {remoteStreams.map((stream, idx) => (<MediaStreamPlayer key={idx} isLocal={false} mediaStream={stream} title="ppp" />))} */}
            {Object.values(m).map(t => (
                <div key={t.name} className="w-[300px]">
                    <MediaStreamPlayer key={t.name} title={t.name} isLocal={false} mediaStream={t.stream} />
                </div>
            ))}
        </div>
    </div>
}