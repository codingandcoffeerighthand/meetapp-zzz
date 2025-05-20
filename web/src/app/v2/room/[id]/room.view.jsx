"use client"

import MediaStreamPlayer from "@/components/video/video.component"
import { useLocalStream } from "@/store/v2/localStream"
import useRemoteStream from "@/store/v2/remoteStream"

export default function RoomView({ roomId }) {
    const { localStreams } = useLocalStream()
    const { mapRemoteStreams } = useRemoteStream()
    return <div className="p-4">
        {localStreams.map((stream, index) => (
            <div key={index} className="w-[300px]">
                <MediaStreamPlayer
                    mediaStream={stream?.stream} title={`local#${index}`} isLocal={true} />
            </div>
        ))}
        {
            Object.values(mapRemoteStreams).map((value) => {
                return value?.streams?.map((stream, index) => (
                    <div key={index} className="w-[300px]">
                        <MediaStreamPlayer mediaStream={stream}
                            title={`${value?.name || value?.address}`} isLocal={false} />
                    </div>
                ))
            })
        }
    </div>
}
