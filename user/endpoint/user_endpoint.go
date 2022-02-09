package endpoint

//endpoint 代表 一个通用的函数原型，负责接收和处理请求返回结果，因为endpoint 的函数形式是固定的，所以可以在外层给endpoint 装饰额外的能力
//比如 熔断，日志，限流，负载均衡等,主要负责 request/response 格式的转换，以及公用拦截器相关的逻辑
//Endpoint 层采用类似洋葱的模型，提供了对日志、限流、熔断、链路追踪和服务监控,负载均衡 等方面的扩展能力。
import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/i6666/micro-go/user/service"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("rate limit exceed")

func NewTokenBucketLimiterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//如果超过流量，直接返回限流异常
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

type UserEndpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
	UserEndpoint     endpoint.Endpoint
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	UserInfo *service.UserInfoDto `json:"user_info"`
}

type LoginFrom struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	UserInfo *service.UserInfoDto `json:"user_info"`
	Err      error
}

func MakeUserEndPoint(userService service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, from interface{}) (result interface{}, err error) {
		req := from.(LoginFrom)
		userInfoDto, err := userService.Login(ctx, req.Email, req.Password)
		return LoginResult{userInfoDto, err}, nil
	}

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
