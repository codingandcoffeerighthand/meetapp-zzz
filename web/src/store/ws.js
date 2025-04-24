import { create } from 'zustand'
import { useWeb3Store } from './web3'
import { room } from './features/wsHandle'
const useWebSocketStore = create((set, get) => ({
    ws: { readyState: "not connected" },
    isLoading: false,
    isConnected: false,
    init: (acc) => {
        if (!acc) {
            return
        }
        set({ isLoading: true })
        const ws = new WebSocket(process.env.NEXT_PUBLIC_WS)
        ws.onopen = async () => {
            console.info("WebSocket connected")
            set({ isConnected: true })
            await get().login()
        }
        ws.onclose = () => {
            set({ isConnected: false })
        }
        set({ ws })
        set({ isLoading: false })
    },

    login: async () => {
        const { web3, account, prK } = useWeb3Store.getState()
        if (!web3 || !account) {
            get().ws?.close()
            set({ ws: null, isConnected: false })
            return
        }
        const nonce = `#nonce ${account}`;
        const { signature } = await web3.eth.accounts
            .sign(
                nonce,
                prK
            );
        const event = JSON.stringify({
            address: account,
            signature: signature,
            nonce: nonce
        })
        const ws = get().ws
        ws.send(event)
    },

    handleWs: async (roomId) => {
        const { ws } = get()
        const { joinedHandler, handleRequireRenegoiate } = room()
        ws.onmessage = async (event) => {
            const json = JSON.parse(event.data)
            console.info(json)
            switch (json.message_name) {
                case "joined_room":
                    await joinedHandler(json)
                    ws.send(JSON.stringify({
                        message_name: "pull-room", data: {
                            room_id: roomId
                        }
                    }))
                    break
                case "require_renegotiate_session":
                    console.info(json)
                    await handleRequireRenegoiate(json.session_id, json.sdp_offer)

                    break
                case "new_participant_joined_room":
                    break
                default:
                    break
            }
        }
    }
}))
export { useWebSocketStore }