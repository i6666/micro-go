package service

import (
	"context"
	"fmt"
	rpc "github.com/i6666/micro-go/user/rpc/pbk"
	"google.golang.org/grpc"
	"testing"
)

func TestLoginService_CheckPassword(t *testing.T) {
	serviceAddress := "127.0.0.1:1234"

	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())

	if err != nil {
		panic("connect error")
	}
	defer conn.Close()

	loginClient := rpc.NewLoginServiceClient(conn)

	userReq := &rpc.LoginRequest{Username: "", Password: ""}
	reply, _ := loginClient.CheckPassword(context.Background(), userReq)

	fmt.Printf("UserService CheckPassword :%s", reply.Ret)

}
