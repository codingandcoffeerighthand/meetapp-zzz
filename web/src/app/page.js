"use client";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useContractStore } from "@/store/contract";
import useWeb3Store from "@/store/web3";
import { set } from "date-fns";
import { Car } from "lucide-react";
import { useEffect, useState } from "react";

export default function Home() {
  const { isConnected, account, connectWallet } = useWeb3Store()
  const { contract, initContract, getOwner, authorize, getAuthorizedBackend } = useContractStore()
  const [ower, setOwer] = useState("")
  const [addresAuthorized, setAddressAuthorized] = useState("")
  const [authorizeIdx, setAuthorizeIdx] = useState(0)
  const getOwnerHandler = async () => {
    const ownerAddress = await getOwner()
    setOwer(ownerAddress)
  }
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
        {isConnected ? (
          <p>account: {account}</p>
        ) : (
          <Button onClick={connectWallet} className="w-[200px]">Connect Wallet</Button>
        )}
        {
          contract ? (
            <>
              <p>Contract: {contract._address}</p>
              <Card className="w-full flex flex-col gap-4">
                <Button onClick={getOwnerHandler} className="w-full">Get Owner</Button>
                <p>Owner: {ower}</p>
              </Card>
              <Card className="w-full flex flex-col gap-4">
                <Button onClick={() => authorize(account)} className="w-full">Authorize</Button>
              </Card>
              <Card className="w-full flex flex-col gap-4">
                <div className="w-full flex gap-2">
                  <input type="number" value={authorizeIdx} onChange={(e) => setAuthorizeIdx(e.target.value)} />
                  <Button onClick={async () => {
                    setAddressAuthorized(await getAuthorizedBackend(
                      authorizeIdx
                    ))
                  }} className="">Get Authorized Backend</Button>
                </div>
                <p>Address Authorized: {addresAuthorized}</p>
              </Card>
            </>

          ) : (
            <>
              <p>Contract not initialized</p>
              <Button onClick={initContract} className="w-[200px]">Init Contract</Button>
            </>
          )
        }
      </main>
    </div>
  );
}
