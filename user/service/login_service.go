package service

import (
	"context"
	rpc "github.com/i6666/micro-go/user/rpc/pbk"
)

type LoginService struct {
}

func (s *LoginService) CheckPassword(ctx context.Context, req *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	if req.Username == "admin" && req.Password == "admin" {
		response := rpc.LoginResponse{Ret: "success"}
		return &response, nil
	}

	response := rpc.LoginResponse{Ret: "fail"}
	return &response, nil
}
