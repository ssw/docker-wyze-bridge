# Wyze Bridge (local only)

The [IDisposable Wyze Bridge](https://github.com/IDisposable/docker-wyze-bridge)
add-on, built from source with `LOCAL_ONLY=true`: the WebUI and go2rtc (API,
RTSP, WebRTC) listen on `127.0.0.1` only, so no camera is reachable from the
network. Home Assistant Core shares the host network and pulls streams from
`rtsp://127.0.0.1:8554/<camera>`; view them in Home Assistant.

No ingress panel and no exposed ports. See upstream issue #153 for why.
