import { useState } from "react";
import { Button } from "./ui/button";
import { Card } from "./ui/card";
import { Input } from "./ui/input";
import { useRouter } from "next/navigation";

export default function CreateRoomComponent({ callCreateRoom, initialRoomId, initialRoomName }) {
    const [roomName, setRoomName] = useState(initialRoomName)
    const [roomId, setRoomId] = useState(initialRoomId)
    const router = useRouter()
    const joinRoom = () => {
        router.push(`/v2/room/${roomId}`)
    }
    return (
        <Card className="p-4 flex flex-row justify-start">
            <div className="flex flex-col gap-2 w-3/4">
                <p>You are authorized to use this contract</p>
                <Input placeholder={`Room Id... ${roomId}`} onChange={(e) => setRoomId(e.target.value)} />
                <Input placeholder={`Room Name... ${roomName}`} onChange={(e) => setRoomName(e.target.value)} />
            </div>
            <div className="flex flex-col gap-2 justify-end">
                <Button
                    onClick={joinRoom}>
                    Join Room With Id
                </Button>
                <Button
                    onClick={() => callCreateRoom(roomId, roomName)}>
                    Create Room
                </Button>
            </div>
        </Card>
    )
}