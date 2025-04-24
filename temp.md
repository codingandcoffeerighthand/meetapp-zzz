```js
// components/WebRTCClient.js
import { useEffect, useRef, useState } from 'react';
import io from 'socket.io-client'; // If using Socket.IO

const signalingServerUrl = 'http://localhost:8080'; // Replace with your signaling server URL

const WebRTCClient = () => {
  const localVideoRef = useRef(null);
  const remoteVideoRef = useRef(null);
  const peerConnection = useRef(null);
  const socket = useRef(null);

  useEffect(() => {
    socket.current = io(signalingServerUrl); // Connect to signaling server

    socket.current.on('connect', () => {
      console.log('Connected to signaling server');
    });

    socket.current.on('message', async (message) => {
      try {
        const parsedMessage = JSON.parse(message);
        // Handle signaling messages (offer, answer, ice candidates)
        if (parsedMessage.type === 'offer') {
          await handleOffer(parsedMessage);
        } else if (parsedMessage.type === 'answer') {
          await handleAnswer(parsedMessage);
        } else if (parsedMessage.type === 'candidate') {
          await handleCandidate(parsedMessage);
        }
      } catch (error) {
        console.error('Error processing signaling message:', error);
      }
    });

    async function startWebcam() {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        if (localVideoRef.current) {
          localVideoRef.current.srcObject = stream;
        }
        stream.getTracks().forEach(track => peerConnection.current.addTrack(track, stream));
      } catch (error) {
        console.error('Error accessing media devices:', error);
      }
    }

    async function createPeerConnection() {
      peerConnection.current = new RTCPeerConnection({
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' },
          // Add more STUN/TURN servers for better NAT traversal
        ],
      });

      peerConnection.current.onicecandidate = (event) => {
        if (event.candidate) {
          socket.current.emit('message', JSON.stringify({ type: 'candidate', candidate: event.candidate }));
        }
      };

      peerConnection.current.ontrack = (event) => {
        if (remoteVideoRef.current && event.streams && event.streams[0]) {
          remoteVideoRef.current.srcObject = event.streams[0];
        }
      };
    }

    async function handleOffer(offer) {
      if (!peerConnection.current) await createPeerConnection();
      await peerConnection.current.setRemoteDescription(new RTCSessionDescription(offer));
      const answer = await peerConnection.current.createAnswer();
      await peerConnection.current.setLocalDescription(new RTCSessionDescription(answer));
      socket.current.emit('message', JSON.stringify({ type: 'answer', sdp: answer.sdp }));
    }

    async function handleAnswer(answer) {
      if (peerConnection.current) {
        await peerConnection.current.setRemoteDescription(new RTCSessionDescription(answer));
      }
    }

    async function handleCandidate(candidate) {
      if (peerConnection.current && candidate.candidate) {
        await peerConnection.current.addIceCandidate(new RTCIceCandidate(candidate.candidate));
      }
    }

    createPeerConnection();
    startWebcam();

    return () => {
      if (socket.current) {
        socket.current.disconnect();
      }
      if (peerConnection.current) {
        peerConnection.current.close();
      }
      if (localVideoRef.current && localVideoRef.current.srcObject) {
        localVideoRef.current.srcObject.getTracks().forEach(track => track.stop());
      }
      if (remoteVideoRef.current && remoteVideoRef.current.srcObject) {
        remoteVideoRef.current.srcObject.getTracks().forEach(track => track.stop());
      }
    };
  }, []);

  return (
    <div>
      <video ref={localVideoRef} autoPlay muted playsInline />
      <video ref={remoteVideoRef} autoPlay playsInline />
    </div>
  );
};

export default WebRTCClient;
```


prompt
```
đọc thông tin trong file cloudflare_realtime.txt và connection API trong meet-app/server/pkg/gencode/cloudflare_client/api.yaml
tham khảo file sample.go
1. publish track
- client send track data and sdp offer lên smartcontract
- proxy nhận data tạo session cloudflare realtime mới và cập nhật session cho các data trên smartcontract
- proxy thực hiện add localtrack vào session
- trả về cho client sdp answer, client thêm vào localpeerconnection

2. pull track
- proxy lấy thông tin các track trong room từ smart contract
- proxy tạo seesion mới cho remote và thêm các track trong room vào remote session này
- proxy gửi sdpoffer nhận được cho client
- client add sdpoffer vào remoteperrconnection và tạo sdf answer gừi renegoiate rtc session

##
kiểm tra kĩ dựa trên các tài liệu có sẵn 2 flow trên, tối ưu thêm

triền khai code
comment rõ ràng

```