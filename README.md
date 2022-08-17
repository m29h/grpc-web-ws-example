# gRPC-Web-Websocket Example: A simple Golang API server and TypeScript frontend

### Get started with HTTP 1.1 and WebSockets with both Client and Server side streaming
Makes use fo the gRPC Web library at https://github.com/improbable-eng/

* `npm install`
* `npm start` to start the Golang server and Webpack dev server
* Go to `http://localhost:8081`

The example program implements a Ping() with client side streaming and a Pong() with server side streaming
* Pong() sends every 1ms the current unix nanosecond timestamp (1000 times)
* The ts browser client echos back each Pong() stream message it receives through its Ping() stream
* once the Pong() stream terminates, the ts client terminates its Ping() stream
* finally the Ping() method echoes back the average round-trip latency in nanoseconds
* result is printed both on client and server side