'use client'
import { useEffect, useRef } from 'react';
export default function VideoComponent({ tracks, isLocal }) {
    const videoRef = useRef();
    useEffect(() => {
        if (videoRef.current && tracks) {
            const stream = new MediaStream([...tracks]);
            videoRef.current.srcObject = stream;
        } else if (videoRef.current) {
            videoRef.current.srcObject = null;
        }
    }, [tracks])
    return (
        <video
            ref={videoRef}
            autoPlay
            muted={isLocal}
            playsInline
            style={{ width: '100%', height: '100%', objectFit: 'cover' }}
        />
    );
}