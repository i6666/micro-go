package rpc

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"testing"
)

//实现rpc 服务端
func TestStringRequest_Concat(t *testing.T) {

	stringService := new(StringService)
	//1.服务注册，通过反射将方法取出存入map
	_ = rpc.Register(stringService)
	rpc.HandleHTTP()
	//2.处理网络调用，监听端口读取数据包，解码请求，调用反射处理后的方法，返回值编码，返回客户端
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	_ = http.Serve(l, nil)

}
