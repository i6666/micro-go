package service

import (
	"context"
	rpc "github.com/i6666/micro-go/user/rpc/pbk"
)

type LoginService struct {
}

//测试代码，service 层用LoginRequest 不符合规范，request 和response 代表请求响应属于 transport和endpoint 层概念
func (s *LoginService) CheckPassword(ctx context.Context, req *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	if req.Username == "admin" && req.Password == "admin" {
		response := rpc.LoginResponse{Ret: "success"}
		return &response, nil
	}

	response := rpc.LoginResponse{Ret: "fail"}
	return &response, nil
}
