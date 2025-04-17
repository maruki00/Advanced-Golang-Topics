package util

import (
	"log"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// Dial: grpc server with new relic or elastic apm middleware ,
func Dial(addr string, opts ...grpc.UnaryClientInterceptor) *grpc.ClientConn {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcMiddleware.ChainUnaryClient(opts...)),
	)
	if err != nil {
		log.Fatal("could not connect to", addr, err)
	}
	return conn
}
