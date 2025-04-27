import RoomView from "@/components/view/rooms.view"

export default async function Room({ params }) {
    const { id } = await params
    return <div className="p-4">
        <h1 className="text-2xl text-center item-center">Room {id}</h1>
        <RoomView roomId={id} />
    </div>
}