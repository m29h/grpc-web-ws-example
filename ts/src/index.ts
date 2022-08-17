import {grpc} from "@improbable-eng/grpc-web";
import {ExampleService} from "../_proto/example/library/example_service_pb_service";
import {msg} from "../_proto/example/library/example_service_pb";

declare const USE_TLS: boolean;
const host = USE_TLS ? "https://localhost:9091" : "http://localhost:9090";


function PingPong() {
  const sendMsg = new msg();
  sendMsg.setVal(0);

  const ping = grpc.client(ExampleService.Ping, {
    host: host,
    transport: grpc.WebsocketTransport(),
  });
  ping.onHeaders((headers: grpc.Metadata) => {
    console.log("ping.onHeaders", headers);
  });
  ping.onMessage((message: msg) => {
    console.log("ping finished, average round-trip time was ",message.getVal()," ns" );
  });
  ping.onEnd((code: grpc.Code, msgi: string, trailers: grpc.Metadata) => {
    console.log("ping.onEnd error code", code, msgi, trailers);
  });
  ping.start();
    
  const pong = grpc.client(ExampleService.Pong, {
    host: host,
    transport: grpc.WebsocketTransport(),
  });
  pong.onHeaders((headers: grpc.Metadata) => {
    console.log("pong.onHeaders", headers);
  });
  pong.onMessage((message: msg) => {
    ping.send(message) // ping back the message we just received from pong
  });
  pong.onEnd((code: grpc.Code, msgi: string, trailers: grpc.Metadata) => {
    console.log("pong.onEnd error code", code, msgi, trailers);
    ping.finishSend()
  });
  pong.start();
  pong.send(sendMsg);
}
PingPong();

