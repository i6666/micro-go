package endpoint

//endpoint 代表 一个通用的函数原型，负责接收和处理请求返回结果，因为endpoint 的函数形式是固定的，所以可以在外层给endpoint 装饰额外的能力
//比如 熔断，日志，限流，负载均衡等
import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/i6666/micro-go/user/service"
)

type UserEndpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	UserInfo *service.UserInfoDto `json:"user_info"`
}

func MakeLoginEndpoint(userService service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*LoginRequest)
		userInfo, err := userService.Login(ctx, req.Email, req.Password)
		return &LoginResponse{UserInfo: userInfo}, err
	}
}

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}
type RegisterResponse struct {
	UserInfo *service.UserInfoDto `json:"user_info"`
}

func MakeRegisterEndpoint(userService service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*RegisterRequest)
		userInfo, err := userService.Register(ctx, &service.RegisterUserVo{
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
		})
		return &RegisterResponse{UserInfo: userInfo}, err
	}
}
