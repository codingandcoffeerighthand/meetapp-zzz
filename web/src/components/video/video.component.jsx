"use client"

import { useState, useRef, useEffect } from "react"
import { X, Volume2, VolumeX, Mic, MicOff } from "lucide-react"

export default function MediaStreamPlayer({
    mediaStream, title,
    closeVideoCallback = () => { },
    isLocal = false, autoPlay = true }) {
    const [isMuted, setIsMuted] = useState(isLocal)
    const [isVisible, setIsVisible] = useState(true)
    const [hasAudio, setHasAudio] = useState(false)
    const [microphoneEnabled, setMicrophoneEnabled] = useState(true)
    const videoRef = useRef(null)

    // Set the MediaStream as the video source and check for audio tracks
    useEffect(() => {
        if (videoRef.current && mediaStream) {
            videoRef.current.srcObject = mediaStream

            // Check if the stream has audio tracks
            const audioTracks = mediaStream.getAudioTracks()
            setHasAudio(audioTracks.length > 0)

            // Set initial microphone state based on the first audio track's enabled state
            if (audioTracks.length > 0) {
                setMicrophoneEnabled(audioTracks[0].enabled)
            }

            // Listen for track additions/removals to update hasAudio state
            const handleTrackChange = (evt) => {

                console.info("add track", evt, mediaStream.getAudioTracks)
                const updatedAudioTracks = mediaStream.getAudioTracks()
                setHasAudio(updatedAudioTracks.length > 0)

                // Update microphone state if tracks change
                if (updatedAudioTracks.length > 0) {
                    setMicrophoneEnabled(updatedAudioTracks[0].enabled)
                }
            }
            const handleInActive = () => {
                setIsVisible(false)
                closeVideoCallback()
            }

            mediaStream.addEventListener("addtrack", handleTrackChange)
            mediaStream.addEventListener("removetrack", handleTrackChange)
            // mediaStream.addEventListener("inactive", handleInActive)

            // Clean up when component unmounts or mediaStream changes
            return () => {
                mediaStream.removeEventListener("addtrack", handleTrackChange)
                mediaStream.removeEventListener("removetrack", handleTrackChange)
                // mediaStream.removeEventListener("inactive", handleInActive)

                if (videoRef.current && videoRef.current.srcObject) {
                    videoRef.current.srcObject = null
                }
            }
        } else {
            // Reset states if there's no mediaStream
            setHasAudio(false)
            setMicrophoneEnabled(true)
        }
    }, [mediaStream])

    // Toggle playback mute (this affects only what the user hears)
    const toggleMute = () => {
        if (videoRef.current && mediaStream.getAudioTracks().length > 0) {
            videoRef.current.muted = !videoRef.current.muted
            setIsMuted(!isMuted)
        }
    }

    // Toggle microphone (this affects the actual audio capture)
    const toggleMicrophone = () => {
        if (mediaStream && hasAudio) {
            const audioTracks = mediaStream.getAudioTracks()

            // Toggle enabled state for all audio tracks
            const newEnabledState = !microphoneEnabled
            audioTracks.forEach((track) => {
                track.enabled = newEnabledState
            })

            setMicrophoneEnabled(newEnabledState)
        }
    }

    const closeVideo = () => {
        mediaStream.getTracks().forEach((track) => track.stop())
        setIsVisible(false)
        closeVideoCallback()
    }

    // If video is closed, don't render anything
    if (!isVisible) return null

    return (
        <div className="relative w-full max-w-3xl mx-auto rounded-lg overflow-hidden shadow-lg">
            {/* Title bar */}
            {title && (
                <div className="bg-gray-800 text-white py-2 px-4 flex items-center justify-between">
                    <h2 className="font-medium text-sm md:text-base truncate">{title}</h2>
                    <div className="flex items-center gap-2">
                        {mediaStream.getAudioTracks().length > 0 ? (
                            <span className="text-xs bg-green-900/50 px-2 py-0.5 rounded">Audio</span>
                        ) : (
                            <span className="text-xs bg-gray-700/50 px-2 py-0.5 rounded">No Audio</span>
                        )}
                        {mediaStream ? (
                            <span className="flex items-center text-xs">
                                <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
                                Live
                            </span>
                        ) : (
                            <span className="flex items-center text-xs">
                                <span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
                                No stream
                            </span>
                        )}
                    </div>
                </div>
            )}

            {/* Video element for MediaStream */}
            <video
                ref={videoRef}
                className="w-full h-auto bg-black"
                autoPlay={autoPlay}
                muted={isLocal || isMuted}
                playsInline
                aria-label={title ? `Video stream: ${title}` : "Video stream"}
            />

            {/* Custom controls overlay */}
            <div className="absolute bottom-0 left-0 right-0 bg-black/50 p-3 flex items-center justify-between">
                {/* Left side - microphone control */}
                <div className="text-white">
                    {isLocal && hasAudio && (
                        <button
                            onClick={toggleMicrophone}
                            className={`hover:text-gray-300 focus:outline-none ${!microphoneEnabled ? "text-red-500" : "text-white"}`}
                            aria-label={microphoneEnabled ? "Turn off microphone" : "Turn on microphone"}
                            title={microphoneEnabled ? "Turn off microphone" : "Turn on microphone"}
                        >
                            {isLocal && microphoneEnabled ? <Mic className="w-5 h-5" /> : <MicOff className="w-5 h-5" />}
                        </button>
                    )}
                </div>

                {/* Right side controls */}
                <div className="flex items-center space-x-4">
                    {/* Mute/Unmute button - only show if stream has audio */}
                    {mediaStream.getAudioTracks().length > 0 && (
                        <button
                            onClick={toggleMute}
                            className="text-white hover:text-gray-300 focus:outline-none"
                            aria-label={isMuted ? "Unmute audio" : "Mute audio"}
                            title={isMuted ? "Unmute audio" : "Mute audio"}
                        >
                            {isMuted ? <VolumeX className="w-5 h-5" /> : <Volume2 className="w-5 h-5" />}
                        </button>
                    )}

                    {/* Close button */}
                    {isLocal &&
                        <button
                            onClick={closeVideo}
                            className="text-white hover:text-gray-300 focus:outline-none"
                            aria-label="Close video"
                            title="Close video"
                        >
                            <X className="w-5 h-5" />
                        </button>
                    }
                </div>
            </div>
        </div>
    )
}
