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
	rpc.Register(stringService)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

}
