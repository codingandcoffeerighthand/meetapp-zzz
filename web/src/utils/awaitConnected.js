export const waitLocalConnection = (localPeerConnection) => new Promise((res, rej) => {
    // timeout after 5s
    setTimeout(rej, 5000);
    const iceConnectionStateChangeHandler = () => {
        if (localPeerConnection.iceConnectionState === "connected") {
            localPeerConnection.removeEventListener(
                "iceconnectionstatechange",
                iceConnectionStateChangeHandler,
            );
            res(undefined);
        }
    };
    localPeerConnection.addEventListener(
        "iceconnectionstatechange",
        iceConnectionStateChangeHandler,
    );
});

export function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}