"use client"
import { useWeb3V2Store } from "@/store/v2/web3_v2"
import LoadingScreen from "@/components/loading.component"
import { Card } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import CreateRoomComponent from "@/components/create-room.component"
import { v4 } from "uuid"
import ErrorComponent from "@/components/error.component"

export default function V2() {
    const { isAuthorized, isLoading, error, isWeb3Connected,
        account, contract,
        addPrivateKey, connectWeb3,
        callRegister, callCreateRoom
    } = useWeb3V2Store()

    return (
        <div className="flex flex-col gap-4 w-[80%] p-20 mx-auto">
            {error && <ErrorComponent error={error} />}
            <Card className="p-4 gap-2 flex flex-row justify-start gap-20">
                <div>
                    <p>Account: {account ? account : "<Empty>"}</p>
                    <p>Contract: {contract ? contract.options.address : "<Empty>"}</p>
                </div>
                <div>
                    <p>Connected: {isWeb3Connected ? "true" : "false"}</p>
                    <p>Authorized: {isAuthorized ? "true" : "false"}</p>
                </div>
            </Card>
            {isWeb3Connected || (
                <Card className="p-4 gap-2 flex flex-row justify-start gap-20">
                    <Input placeholder="Private Key" onChange={(e) => addPrivateKey(e.target.value)} />
                    <Button onClick={connectWeb3}>Connect</Button>
                </Card>
            )}
            {isWeb3Connected && !isAuthorized && (
                <Card className="p-4 gap-2 flex flex-row justify-start gap-20">
                    <p>You are not authorized to use this contract</p>
                    <Button onClick={callRegister}>
                        Register
                    </Button>
                </Card>
            )}
            {isAuthorized && isAuthorized && (
                <CreateRoomComponent callCreateRoom={callCreateRoom} initialRoomId={v4()} initialRoomName={``} />
            )}
            {isLoading && <LoadingScreen />}
        </div>
    )
}