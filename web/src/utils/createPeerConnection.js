async function createPeerConnection() {
    const peerConnection = new RTCPeerConnection({
        iceServers: [
            {
                urls: "stun:stun.cloudflare.com:3478",
            },
        ],
        bundlePolicy: "max-bundle",
    });

    return peerConnection;
}

export { createPeerConnection }