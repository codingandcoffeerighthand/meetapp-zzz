import RoomView from "@/components/view/rooms.view"

export default async function Room({ params }) {
    const { id } = await params
    return <>
        <h1>Room {id}</h1>
        <p>Room</p>
        <RoomView roomId={id} />
    </>
}