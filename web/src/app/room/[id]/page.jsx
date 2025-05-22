import ParticipantComponent from "@/components/paritcipant.component"
import RoomView from "./room.view"
export default async function Room({ params }) {
    const { id } = await params
    return <div className="flex flex-col gap-4 p-4 mx-auto">
        <h1 className="text-2xl text-center item-center">Room {id}</h1>
        <ParticipantComponent roomId={id} />
        <RoomView roomId={id} />
    </div>
}