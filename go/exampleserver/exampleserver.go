package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"

	library "grpc-web-ws-example/go/_proto/example/library"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
)

func main() {
	flag.Parse()

	port := 9090

	grpcServer := grpc.NewServer()
	library.RegisterExampleServiceServer(grpcServer, &exampleService{})
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, io.Discard, io.Discard))

	wrappedServer := grpcweb.WrapServer(grpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false), grpcweb.WithOriginFunc(func(string) bool { return true }),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool { return true }))
	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handler),
	}

	grpclog.Infof("Starting server. http port: %d", port)

	if err := httpServer.ListenAndServe(); err != nil {
		grpclog.Fatalf("failed starting http server: %v", err)
	}

}

type exampleService struct {
	library.UnimplementedExampleServiceServer
}

func (s *exampleService) Ping(stream library.ExampleService_PingServer) error {
	b := library.Msg{}
	count := int64(0)
	tot := int64(0)
	for {
		if stream.RecvMsg(&b) != nil {
			break
		}
		count++
		dt := time.Now().UnixNano() - b.Val
		tot += dt
	}

	grpclog.Infof("ping finished, average round-trip time was: %d nanoseconds, over %d messages", tot/count, count)
	return stream.SendAndClose(&library.Msg{Val: tot / count})
}

func (s *exampleService) Pong(cmd *library.Msg, stream library.ExampleService_PongServer) error {
	stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Millisecond)
		b := time.Now().UnixNano()
		cmd.Val = b
		stream.SendMsg(cmd)
	}
	stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))

	return nil
}
