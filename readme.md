# WebRTC Signaling Server
A simple signaling server for WebRTC streaming. This server is responsible for creating connections between streamers and subscribers for the WebRTC streaming testing described in my portfolio [here](https://ji-0.github.io/posts/webrtc/).

## Specifications
This server accepts connections on ports 3000 and 3001 which are upgraded to a websocket connection. After that the client can either request to be upgraded to a streamer, which means that the connection to the server is kept alive and the client can be found by other peers, or a peer can request a stream either at random or based on a name/token and is disconnected after the WebRTC handshake.

## Instructions
To run the server simply clone the repo and run the project with `go run .`.