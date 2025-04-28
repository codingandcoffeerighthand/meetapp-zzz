"use client";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useWeb3Store } from "@/store/web3";
import { Car } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function Home() {
  const {
    isConnected, account, addKey, getOwner, contract, initContract, isAuthorized,
    register, createRoom, isLoading, err
  } = useWeb3Store()
  const [ower, setOwer] = useState("")
  const [roomIdJoin, setRoomIdJoin] = useState("")
  const [prk, setPrk] = useState("")
  const router = useRouter()
  const [roomName, setRoomName] = useState("")
  const getOwnerHandler = async () => {
    const ownerAddress = await getOwner()
    setOwer(ownerAddress)
  }
  const addPrkHandler = async () => {
    try {
      await addKey(prk)
    } catch (err) {
      console.error(err)
    }
  }
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
        {isLoading && <p className="text-red-500">loading...</p>}
        {isConnected ? (
          <>
            <p>account: {account}</p>
          </>
        ) : (
          <>
            <Input placeholder="private key" onChange={(e) => setPrk(e.target.value)} />
            <Button onClick={addPrkHandler} className="w-[200px]">Add private key</Button>
          </>
        )}
        {contract ? (
          isAuthorized ?
            <>
              <Card className="w-full p-4">
                {/* <Input placeholder="room name" type="text" onChange={(e) => setRoomName(e.target.value)} /> */}
                <Button onClick={() => createRoom(
                  roomName, (roomId) => router.push(`/room/${roomId}`)
                )} className="">Create Room</Button>
              </Card>
            </>
            :
            <Button onClick={register} className="w-[200px]">Register</Button>
        ) : (
          <Button onClick={initContract} className="w-[200px]">Init Contract</Button>
        )}
        <div className="w-full">
          <Card className="w-full p-4">

            <h2>Join room</h2>
            <Input placeholder="room name" type="text" onChange={(e) => setRoomIdJoin(e.target.value)}
              className="border-1 border-black"
            />
            <Button asChild>
              <Link href={`/room/${roomIdJoin}`} className="flex gap-2 items-center"> <Car /> Join Room</Link>
            </Button>
          </Card>
        </div>
      </main>
    </div >
  );
}
