package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux" // HTTP 请求的路由和分发器
	"github.com/i6666/micro-go/user/endpoint"
	rpc "github.com/i6666/micro-go/user/rpc/pbk"
	"net/http"
	"os"
)

var ErrorBadRequest = errors.New("invalid request parameter")

type gRpcServer struct {
	checkPassword grpc.Handler
}

func (s *gRpcServer) CheckPassword(ctx context.Context, r *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	_, resp, err := s.checkPassword.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*rpc.LoginResponse), nil
}

func NewUserServer(ctx context.Context, endpoints endpoint.UserEndpoints) rpc.LoginServiceServer {
	return &gRpcServer{
		checkPassword: grpc.NewServer(
			endpoints.UserEndpoint,
			DecodeLoginRequest,
			EncodeLoginResponse,
		),
	}

}
func DecodeLoginRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*rpc.LoginRequest)
	return endpoint.LoginFrom{Email: req.Username, Password: req.Password}, nil
}
func EncodeLoginResponse(_ context.Context, r interface{}) (interface{}, error) {
	result := r.(endpoint.LoginResult)
	return &rpc.LoginResponse{
		Ret: result.UserInfo.Username,
		Err: result.Err.Error(),
	}, nil

}

func MakeHttpHandler(ctx context.Context, endpoints *endpoint.UserEndpoints) http.Handler {
	r := mux.NewRouter()
	kitLog := log.NewLogfmtLogger(os.Stderr)

	kitLog = log.With(kitLog, "ts", log.DefaultTimestampUTC)
	kitLog = log.With(kitLog, "caller", log.DefaultCaller)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLog)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeJSONResponse,
		options...,
	))
	r.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeJSONResponse,
		options...,
	))

	return r

}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		return nil, ErrorBadRequest
	}
	return &endpoint.LoginRequest{
		Email:    email,
		Password: password,
	}, nil
}
func encodeJSONResponse(_ context.Context, w http.ResponseWriter, i interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(i)
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if username == "" || password == "" || email == "" {
		return nil, ErrorBadRequest
	}
	return &endpoint.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	}, nil

}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
