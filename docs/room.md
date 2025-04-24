# Design
Để xây dựng một ứng dụng họp trực tuyến (meeting app) sử dụng Cloudflare Realtime, nơi mỗi người dùng trong một phòng có một session Cloudflare Realtime riêng biệt, bạn cần một kiến trúc bao gồm client (ứng dụng web hoặc di động), một máy chủ backend (lớp ứng dụng của bạn) và dịch vụ Cloudflare Realtime đóng vai trò là SFU (Selective Forwarding Unit).

**1. Kiến trúc Tổng Quan:**

* **Client (Web/Mobile App):** Chạy trên thiết bị của người dùng. Có nhiệm vụ:
    * Thu thập luồng âm thanh/video từ micrô và camera của người dùng (sử dụng API WebRTC của trình duyệt/hệ điều hành).
    * Hiển thị các luồng âm thanh/video nhận được từ người dùng khác trong phòng.
    * Quản lý kết nối WebRTC (`RTCPeerConnection`).
    * Giao tiếp với máy chủ backend của bạn để thực hiện signaling (trao đổi thông tin kết nối).
* **Máy Chủ Backend (Lớp Ứng Dụng của bạn):** Đây là "bộ não" của ứng dụng meeting. Có nhiệm vụ:
    * Quản lý các "phòng họp" (room) và danh sách những người dùng (participants) trong mỗi phòng.
    * Xử lý logic tham gia/rời phòng, xác thực người dùng.
    * Đóng vai trò là máy chủ signaling chính. Mọi thông tin WebRTC (như SDP offer/answer) và yêu cầu thêm/bớt track đều đi qua đây.
    * Gọi Cloudflare Realtime API (HTTP API) để tạo session mới cho mỗi người dùng, thêm/bớt các track vào session đó, và kích hoạt quá trình renegotiation khi cần.
    * Theo dõi track nào của người dùng nào đang được publish và người dùng nào cần subscribe các track đó.
    * Điều phối việc trao đổi các thông tin cần thiết giữa các session Cloudflare Realtime của những người dùng trong cùng một phòng.
* **Cloudflare Realtime:** Dịch vụ SFU phân tán. Có nhiệm vụ:
    * Tiếp nhận luồng media (audio/video/data) từ session của những người dùng publish track.
    * Chuyển tiếp (forward) các luồng media này đến session của những người dùng subscribe track đó, dựa trên chỉ thị từ máy chủ backend của bạn.
    * Quản lý kết nối WebRTC ở tầng media, xử lý NAT traversal (thông qua STUN/TURN nếu cần, Cloudflare Realtime cung cấp STUN server).

**2. Quy Trình Chi Tiết:**

Để một người dùng tham gia và tương tác trong một phòng họp:

* **Bước 1: Người dùng tham gia Phòng họp (Client -> Backend):**
    * Người dùng mở ứng dụng client và chọn tham gia một phòng họp cụ thể (ví dụ: nhập mã phòng).
    * Client gửi yêu cầu đến máy chủ backend của bạn (ví dụ: API endpoint `/joinRoom`). Yêu cầu này bao gồm thông tin người dùng và phòng muốn tham gia.

* **Bước 2: Backend xử lý yêu cầu tham gia phòng và tạo Cloudflare Realtime Session (Backend -> Cloudflare Realtime API):**
    * Máy chủ backend của bạn xác thực người dùng và kiểm tra quyền tham gia phòng.
    * Backend ghi nhận người dùng này đã vào phòng.
    * Quan trọng nhất, backend gọi Cloudflare Realtime API để tạo một session mới cho người dùng này:
        * Gọi `POST /apps/{appId}/sessions/new` với `appId` của ứng dụng Cloudflare Realtime của bạn.
        * Cloudflare Realtime API trả về một phản hồi chứa `sessionId` duy nhất cho session của người dùng này. Backend cần lưu lại `sessionId` này, liên kết nó với người dùng và phòng hiện tại.

* **Bước 3: Backend thông báo cho các Client (Backend -> Clients):**
    * Backend thông báo cho client vừa tham gia về `sessionId` vừa được tạo.
    * Backend thông báo cho *tất cả các client khác* trong cùng phòng rằng có một người dùng mới đã tham gia. Thông báo này nên bao gồm thông tin định danh của người dùng mới (không nhất thiết là `sessionId` Cloudflare Realtime, có thể là ID người dùng nội bộ của bạn).

* **Bước 4: Client khởi tạo WebRTC PeerConnection và tạo SDP Offer (Client):**
    * Client vừa tham gia (và các client khác khi nhận được thông báo có người mới) khởi tạo một đối tượng `RTCPeerConnection` của WebRTC.
    * Client thêm các track media cục bộ của mình (audio, video từ mic/camera) vào `RTCPeerConnection`.
    * Client tạo một SDP Offer bằng cách gọi `peerConnection.createOffer()`.

* **Bước 5: Client gửi SDP Offer đến Backend (Client -> Backend):**
    * Client gửi SDP Offer vừa tạo đến máy chủ backend của bạn. Yêu cầu này bao gồm SDP Offer và `sessionId` Cloudflare Realtime của chính client đó.

* **Bước 6: Backend xử lý SDP Offer và thêm Local Tracks vào Session Cloudflare Realtime (Backend -> Cloudflare Realtime API):**
    * Backend nhận SDP Offer từ Client A (với `sessionId A`).
    * Backend gọi Cloudflare Realtime API để thông báo về các track local mà Client A đang publish:
        * Gọi `POST /apps/{appId}/sessions/{sessionId A}/tracks/new`.
        * Trong request body (`TracksRequest`), bạn sẽ gửi `sessionDescription` chứa SDP Offer nhận được từ Client A và danh sách các `TrackObject` có `location: local`, kèm theo `mid` và `trackName` cho từng track (ví dụ: "mic", "camera"). `trackName` là định danh mà bạn sẽ dùng để các client khác subscribe track này sau này.
    * Cloudflare Realtime xử lý yêu cầu và có thể trả về một SDP Answer trong phản hồi (`TracksResponse`). Phản hồi này cũng chứa thông tin về các track local vừa được thêm vào session A, bao gồm `mid` do SFU gán. Backend cần lưu trữ liên kết giữa `sessionId` của người dùng A, `trackName` của từng track của họ, và `mid` tương ứng trong session A.

* **Bước 7: Backend gửi SDP Answer về Client (Backend -> Client):**
    * Backend gửi SDP Answer nhận được từ Cloudflare Realtime (trong Bước 6) về Client A.
    * Client A nhận SDP Answer và đặt nó làm remote description cho `RTCPeerConnection` của mình (`peerConnection.setRemoteDescription(answer)`).

* **Bước 8: Backend yêu cầu các Client khác Subscribe Remote Tracks từ Client mới (Backend -> Cloudflare Realtime API cho từng Client khác -> Backend -> Client):**
    * Đây là bước cốt lõi để các client trong phòng nhìn/nghe thấy nhau. Đối với *mỗi* client B (với `sessionId B`) trong phòng (ngoại trừ Client A), backend cần yêu cầu session B subscribe các track local của Client A.
    * Đối với Client B, backend gọi Cloudflare Realtime API:
        * Gọi `POST /apps/{appId}/sessions/{sessionId B}/tracks/new`.
        * Trong request body (`TracksRequest`), bạn sẽ gửi một danh sách các `TrackObject` có `location: remote`. Mỗi `TrackObject` remote sẽ chỉ định `sessionId` của người publish (là `sessionId A`), và `trackName` của track muốn subscribe (ví dụ: "mic" hoặc "camera" của Client A).
    * Cloudflare Realtime xử lý yêu cầu và trả về một SDP Offer trong phản hồi (`TracksResponse`) cho backend. Phản hồi này cũng có thể chứa cờ `requiresImmediateRenegotiation: true`.
    * Backend gửi SDP Offer nhận được từ Cloudflare Realtime về Client B.
    * Client B nhận SDP Offer và đặt nó làm remote description (`peerConnection.setRemoteDescription(offer)`). Sau đó, Client B tạo một SDP Answer và gửi lại cho backend.
    * Backend nhận SDP Answer từ Client B và gửi nó đến Cloudflare Realtime để hoàn tất quá trình renegotiation cho session B:
        * Gọi `PUT /apps/{appId}/sessions/{sessionId B}/renegotiate`.
        * Trong request body (`RenegotiateRequest`), gửi `sessionDescription` chứa SDP Answer nhận được từ Client B.
    * Cloudflare Realtime xử lý Answer. Từ thời điểm này, Cloudflare Realtime sẽ bắt đầu forward luồng media từ session A (các track audio/video local của Client A) đến session B (như các track remote). Client B sẽ nhận được các luồng remote này qua `RTCPeerConnection` của mình và hiển thị chúng.

* **Bước 9: Lặp lại quá trình Subscribe:**
    * Quy trình ở Bước 8 cần được lặp lại cho *tất cả* các cặp người dùng trong phòng. Ví dụ: Client A cần subscribe track của Client B, Client C cần subscribe track của Client A và Client B, v.v. Backend chịu trách nhiệm theo dõi và điều phối tất cả các yêu cầu subscribe này.

* **Bước 10: Xử lý các sự kiện khác:**
    * **Tắt/Bật mic/camera:** Client phát hiện người dùng tắt/bật thiết bị media. Client gửi thông báo đến backend. Backend gọi `PUT /apps/{appId}/sessions/{sessionId}/tracks/close` để tạm dừng hoặc xóa track tương ứng khỏi session Cloudflare Realtime, hoặc sử dụng `PUT /apps/{appId}/sessions/{sessionId}/tracks/update` để thay đổi trạng thái của track (ví dụ: tắt luồng media nhưng giữ transceiver). Backend cũng cần thông báo cho các client khác trong phòng về sự thay đổi này để họ có thể cập nhật giao diện (ví dụ: hiển thị biểu tượng micrô bị tắt).
    * **Người dùng rời phòng:** Client gửi yêu cầu rời phòng đến backend. Backend ghi nhận người dùng đã rời đi. Backend gọi API để đóng session Cloudflare Realtime của người dùng đó (API `PUT /apps/{appId}/sessions/{sessionId}/tracks/close` có thể được dùng để đóng tất cả các track hoặc có thể có API riêng để đóng toàn bộ session). Backend thông báo cho các client khác trong phòng rằng người dùng này đã rời đi.
    * **ICE Candidates:** WebRTC client sẽ tạo ICE candidates trong quá trình thiết lập kết nối. Mặc dù Cloudflare Realtime API chủ yếu tập trung vào signaling SDP, việc trao đổi ICE candidates thường diễn ra trực tiếp giữa client và SFU sau khi SDP được trao đổi, sử dụng thông tin trong SDP. Cloudflare Realtime hoạt động như một điểm cuối cho ICE. Backend của bạn không cần xử lý ICE candidates trực tiếp giữa các client, chỉ cần đảm bảo SDP exchange diễn ra đúng.
    * **Data Channels:** Nếu bạn cần trao đổi dữ liệu phi-media (ví dụ: tin nhắn chat, tín hiệu giơ tay, dữ liệu màn hình chia sẻ), bạn có thể sử dụng Data Channels. Lưu ý rằng Cloudflare Realtime SFU Data Channels mặc định là một chiều (publisher -> subscribers). Đối với giao tiếp hai chiều hoặc broadcast dữ liệu từ một người đến nhiều người, backend của bạn sẽ cần nhận dữ liệu từ Data Channel của người gửi và sau đó phân phối lại (có thể thông qua các Data Channel khác hoặc kênh giao tiếp riêng) đến những người nhận trong phòng.

**3. Sử dụng OpenAPI Spec:**

OpenAPI spec bạn cung cấp là tài liệu tham khảo chi tiết cho các cuộc gọi API mà máy chủ backend của bạn cần thực hiện. Bạn sẽ dựa vào nó để biết:

* Đường dẫn (path) và phương thức (method) cho mỗi thao tác (tạo session, thêm track, renegotiate, đóng track, lấy trạng thái session).
* Các tham số cần thiết cho mỗi yêu cầu (trong path, query).
* Cấu trúc của request body (payload gửi đi), bao gồm các schema như `TracksRequest`, `RenegotiateRequest`, `CloseTracksRequest`, `UpdateTracksRequest`.
* Cấu trúc của response body (dữ liệu nhận về), bao gồm các schema như `NewSessionResponse`, `TracksResponse`, `RenegotiateResponse`, `CloseTracksResponse`, `GetSessionStateResponse`, `UpdateTracksResponse`.
* Định nghĩa chi tiết của các đối tượng schema được tham chiếu (như `SessionDescription`, `TrackObject`, `CloseTrackObject`), mô tả các trường dữ liệu, kiểu dữ liệu, và mục đích của chúng.
* Yêu cầu bảo mật (`security: - secret: []`) cho thấy bạn cần sử dụng Bearer Token (App Secret của Cloudflare Realtime App) để xác thực các cuộc gọi API từ backend.

**Tóm lại:**

Việc xây dựng ứng dụng meeting với Cloudflare Realtime và session riêng cho mỗi người dùng đòi hỏi một backend server đủ mạnh để quản lý trạng thái phòng, danh sách người dùng và thực hiện *tất cả* các cuộc gọi signaling cần thiết đến Cloudflare Realtime API cho từng session của người dùng. Cloudflare Realtime đảm nhận phần khó khăn của việc xử lý và chuyển tiếp luồng media phân tán, giúp bạn không phải xây dựng và quản lý một hệ thống SFU phức tạp. Bạn sử dụng OpenAPI spec đã cung cấp làm bản đồ chi tiết để lập trình phần giao tiếp giữa backend của bạn và Cloudflare Realtime API.