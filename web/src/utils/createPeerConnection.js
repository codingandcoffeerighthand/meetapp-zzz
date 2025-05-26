async function createPeerConnection(servers) {
    const iceServers = [
        {
            urls: "stun:stun.cloudflare.com:3478",
        },
    ].push(...servers)
    const peerConnection = new RTCPeerConnection({
        iceServers,
        bundlePolicy: "max-bundle",
    });

    return peerConnection;
}

export { createPeerConnection }