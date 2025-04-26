import web3 from "web3"

function listenContractEvts(contract, address, handlers = {}) {
    const evtSub = contract.events.EventForwardedToFrontend({
        filter: { participant: address },
    }, (err, evt) => { console.info(err, evt) })
    evtSub.on("data", data => {
        // const json = JSON.parse(data)
        const json = JSON.parse(web3.utils.hexToUtf8(data.returnValues.eventData))
        console.info(json)
        switch (json?.event_name) {
            case "joined_room":
                console.info("joined room")
                handlers?.joined_room?.(json)
                break
            case "pull_track":
                console.info("pull track")
                handlers?.pull_track?.(json)
                break
            case "local_peer_connection_suscess":
                console.info("local peer connection success")
                break
            case "new_participant_joined":
                console.info("new joined room")
                break
            default:
                console.error("Unknown event", json.event_name)
        }
    })
    evtSub.on("error", err => console.error(err))
}

export { listenContractEvts }