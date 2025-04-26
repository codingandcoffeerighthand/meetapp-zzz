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
    register, createRoom, isLoading
  } = useWeb3Store()
  const [ower, setOwer] = useState("")
  const [roomIdJoin, setRoomIdJoin] = useState("")
  const [prk, setPrk] = useState("")
  const router = useRouter()
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
            <Link href="/page1" className="flex gap-2 items-center"> <Car /> Page 1</Link>
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
              <Card className="w-[200px]">
                <Button onClick={getOwnerHandler} className="w-[200px]">Get Owner</Button>
                <p>owner: {ower}</p>
              </Card>
              <Card className="w-[200px]">
                <Button onClick={() => createRoom(
                  "", (roomId) => router.push(`/room/${roomId}`)
                )} className="w-[200px]">Create Room</Button>
              </Card>
            </>
            :
            <Button onClick={register} className="w-[200px]">Register</Button>
        ) : (
          <Button onClick={initContract} className="w-[200px]">Init Contract</Button>
        )}
        <div>
          <h2>Join room</h2>
          <label> room id </label>
          <input type="text" onChange={(e) => setRoomIdJoin(e.target.value)}
            className="border-1 border-black"
          />
          <Link href={`/room/${roomIdJoin}`} className="flex gap-2 items-center"> <Car /> Join Room</Link>
        </div>
      </main>
    </div >
  );
}
