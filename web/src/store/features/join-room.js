import { useWeb3Store } from "@/store/web3"

const joinRoom = async ({ roomId, participantNam, tracks, sdpOffer }) => {
    const { account, contract } = useWeb3Store.getState()
    if (!contract || !account) {
        throw new Error("Web3 not initialized or no accounts found")
    }
    await contract.methods.joinRoom(roomId, participantNam, tracks, sdpOffer).send({ from: account })
}

export { joinRoom }