async function createPeerConnection(servers = []) {
    const iceServers = [
        {
            urls: "stun:stun.cloudflare.com:3478",
        },
        ...servers
    ]
    console.info(iceServers)
    const peerConnection = new RTCPeerConnection({
        iceServers,
        bundlePolicy: "max-bundle",
    });

    return peerConnection;
}

export { createPeerConnection }