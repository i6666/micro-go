package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/i6666/micro-go/user/dao"
	"github.com/i6666/micro-go/user/endpoint"
	"github.com/i6666/micro-go/user/redis"
	rpc2 "github.com/i6666/micro-go/user/rpc"
	rpc3 "github.com/i6666/micro-go/user/rpc/pbk"
	"github.com/i6666/micro-go/user/service"
	"github.com/i6666/micro-go/user/transport"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	//microStart()

	//goRpc()

	flag.Parse()

	l, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	loginService := new(service.LoginService)
	rpc3.RegisterLoginServiceServer(server, loginService)

	_ = server.Serve(l)

}

func goRpc() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")

	if err != nil {
		log.Fatal("dialing", err)
	}
	stringReq := rpc2.StringRequest{"A", "B"}

	var reply string
	//同步调用
	err = client.Call("StringService.Concat", stringReq, &reply)

	if err != nil {
		log.Fatal("Concat error", err)
	}

	fmt.Println("==========>>>>>>>>" + reply)
}

func microStart() {
	//服务端监听
	var servicePort = flag.Int("service.port", 10086, "service port")

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	err := dao.InitMysql("180.76.118.233", "3306", "root", "strong", "user")

	if err != nil {
		log.Fatal(err)
	}

	err = redis.InitRedis("jlili.cn", "6379", "strong")

	if err != nil {
		log.Fatal(err)
	}

	userService := service.MakeUserServiceImpl(&dao.UserDaoImpl{})

	userEndpoints := &endpoint.UserEndpoints{
		RegisterEndpoint: endpoint.MakeRegisterEndpoint(userService),
		LoginEndpoint:    endpoint.MakeLoginEndpoint(userService),
	}

	r := transport.MakeHttpHandler(ctx, userEndpoints)

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), r)
	}()

	go func() {
		// 监控系统信号，等待ctrl +c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	err2 := <-errChan
	log.Println(err2)
}
