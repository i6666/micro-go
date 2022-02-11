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
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	//microStart()

	//goRpc()

	//gRpc1()
	//gRpc2()

	instanceA := ServiceInstance{"A", 3, 0}
	instanceB := ServiceInstance{"B", 2, 0}
	instanceC := ServiceInstance{"C", 1, 0}

	instances := [3]*ServiceInstance{&instanceA, &instanceB, &instanceC}
	for i := 0; i < 6; i++ {
		selectInstance(instances)
	}

}

//平滑轮询算法
func selectInstance(instances [3]*ServiceInstance) (best *ServiceInstance) {
	total := 0

	for i := 0; i < len(instances); i++ {
		w := instances[i]
		w.curWeight += w.Weight
		total += w.Weight

		if best == nil || w.curWeight > best.curWeight {
			best = w
		}
	}

	best.curWeight -= total

	fmt.Println("best is ", best, "===========", instances[0], instances[1], instances[2])

	return best
}

type ServiceInstance struct {
	name      string
	Weight    int
	curWeight int
}

func gRpc2() {
	flag.Parse()
	ctx := context.Background()
	var svc service.UserService
	svc = service.MakeUserServiceImpl(&dao.UserDaoImpl{})

	et := endpoint.MakeUserEndPoint(svc)
	// 构造限流中间件
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)
	et = endpoint.NewTokenBucketLimiterWithBuildIn(ratebucket)(et)
	endpts := endpoint.UserEndpoints{
		UserEndpoint: et,
	}
	// 使用 transport 构造 UserServiceServer
	handle := transport.NewUserServer(ctx, endpts)
	l, _ := net.Listen("tcp", "127.0.0.1:1234")
	server := grpc.NewServer()
	rpc3.RegisterLoginServiceServer(server, handle)
	_ = server.Serve(l)
}

func gRpc1() {
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
