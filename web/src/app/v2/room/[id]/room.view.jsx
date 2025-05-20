"use client"

import MediaStreamPlayer from "@/components/video/video.component"
import { useLocalStream } from "@/store/v2/localStream"
import useRemoteStream from "@/store/v2/remoteStream"
import { useWeb3V2Store } from "@/store/v2/web3_v2"

export default function RoomView({ roomId }) {
    const { localStreams } = useLocalStream()
    const { mapRemoteStreams } = useRemoteStream()
    const { removeStream } = useWeb3V2Store()
    const handleRemoveStream = async (streamNum) => {
        await removeStream(streamNum)
    }
    return <div className="p-4 flex flex-wrap gap-4">
        {localStreams.map((stream, index) => (
            <div key={index} className="max-w-[600px]">
                <MediaStreamPlayer
                    mediaStream={stream?.stream} title={`local#${index}`} isLocal={true}
                    closeVideoCallback={() => handleRemoveStream(index)}
                />
            </div>
        ))}
        {
            Object.values(mapRemoteStreams).map((value) => {
                return value?.streams?.map((stream, index) => (
                    <div key={index} className="max-w-[600px]">
                        <MediaStreamPlayer mediaStream={stream}
                            title={`${value?.name || value?.address}`} isLocal={false} />
                    </div>
                ))
            })
        }
    </div>
}
