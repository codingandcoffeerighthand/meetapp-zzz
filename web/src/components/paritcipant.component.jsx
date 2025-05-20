"use client"

import { useWeb3V2Store } from "@/store/v2/web3_v2"
import { Button } from "./ui/button"
import { Card } from "./ui/card"
import { Input } from "./ui/input"
import { Label } from "./ui/label"
import { useState } from "react"
import ErrorComponent from "./error.component"
import LoadingScreen from "./loading.component"

export default function ParticipantComponent({ roomId }) {
    const { account, isLoading: web3Loading, error: web3Error,
        handleJoinRoom, callLeaveRoom, addLocalTracks
    } = useWeb3V2Store()
    if (!account) {
        if (window) {
            window.location.href = "/v2"
        }
    }
    const [participantName, setParticipantName] = useState("")
    const [isStream, setIsStream] = useState(false)
    const handlerAddLocalTracks = async () => {
        const stream = await navigator.mediaDevices.getDisplayMedia({
            video: true,
        })
        await addLocalTracks(stream, roomId)
    }
    return (
        <>
            {web3Error && <ErrorComponent error={web3Error} />}
            <div className="flex flex-row gap-4 justify-between">

                <Card className="p-4 flex flex-row gap-4 w-full">
                    <div className="flex flex-row gap-2 w-full">
                        <Label htmlFor="name">Name</Label>
                        <Input id="name" placeholder={`Participant Name ... ${account} `} onChange={(e) => setParticipantName(e.target.value)} />
                    </div>
                </Card>
                <Card className="p-4 flex flex-row gap-4 w-full">
                    <Button
                        disabled={isStream}
                        onClick={() => {
                            handleJoinRoom(roomId, participantName)
                            setIsStream(true)
                        }}
                    >
                        Start stream
                    </Button>
                    <Button
                        onClick={() => {
                            handlerAddLocalTracks()
                        }}
                    >
                        Share screen
                    </Button>
                    <Button
                        onClick={() => {
                            callLeaveRoom(roomId, () => {
                                if (window) {
                                    window.location.href = "/v2"
                                }
                            })
                        }}
                    >
                        Leave room
                    </Button>
                </Card>
            </div>
            {web3Loading && <LoadingScreen />}
        </>
    )
}
