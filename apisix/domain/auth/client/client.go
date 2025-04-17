package client

import (
	"context"
	"os"

	auth_gen "apisix-api/proto/gen"
	Grpc "apisix-api/util"
)

func Login(ctx context.Context, req *auth_gen.LoginRequest) (*auth_gen.LoginResponse, error) {
	conn := Grpc.Dial(os.Getenv("services_grpc"))
	defer conn.Close()

	client := auth_gen.NewAuthClient(conn)

	return client.Login(ctx, req)
}
