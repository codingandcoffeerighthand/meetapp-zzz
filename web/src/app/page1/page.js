"use client";
import useWeb3Store from "@/store/web3";
import { useWebSocketStore } from "@/store/ws";
import { useEffect, useState } from "react";
export default function Page() {
    const { isConnected, account, connectWallet } = useWeb3Store()
    const { ws } = useWebSocketStore()
    const [messages, setMessages] = useState([])
    useEffect(() => {
        if (ws) {

            ws.onmessage = (event) => {
                console.log("Message from server ", event.data);
                setMessages((prevMessages) => [...prevMessages, event.data])
            }
        }
    }, [ws])
    console.info(ws)
    return (
        <div className="mx-auto my-20 flex justify-center items-center flex-col gap-4">
            {isConnected ? (
                <p>account: {account}</p>
            ) : (
                <p>Not Connected</p>
            )}
            {isConnected && <p>Connected</p>}
            <br />
            <p>websocket</p>
            <p>readyState: {ws?.readyState || "not connected"}</p>
            <button
                onClick={() => {
                    ws.send("Hello from client")
                }}
                className="bg-blue-500 text-white px-4 py-2 rounded"
            >
                Send Message
            </button>
            {
                messages.map((message, index) => (
                    <div key={index} className="bg-gray-200 p-2 rounded">
                        {message}
                    </div>
                ))
            }
        </div>
    );
}