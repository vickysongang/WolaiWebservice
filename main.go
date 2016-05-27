package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/rpc"
	"net/rpc/jsonrpc"
	"runtime"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/config"
	"WolaiWebservice/logger"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/routers"
	myrpc "WolaiWebservice/rpc"
)

func init() {
	var err error

	//针对Golang 1.5以后版本，设置最大核数
	runtime.GOMAXPROCS(config.Env.Server.Maxprocs)

	//加载seelog的配置文件，使用配置文件里的方式输出日志信息
	err = logger.Initialize()
	if err != nil {
		fmt.Println(err.Error())
	}

	//注册数据库
	err = models.Initialize()
	if err != nil {
		seelog.Critical(err)
	}

	err = redis.Initialize()
	if err != nil {
		seelog.Critical(err)
	}
}

func startRpcServer() {
	lis, err := net.Listen("tcp", config.Env.Server.RpcPort)
	if err != nil {
		seelog.Critical("RPC端口被占用")
	}
	defer lis.Close()
	watcher := new(myrpc.RpcWatcher)
	srv := rpc.NewServer()
	srv.RegisterName("RpcWatcher", watcher)
	for {
		conn, _ := lis.Accept()
		go srv.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func main() {
	orm.Debug = false

	go startRpcServer()
	//go handlerv2.VoipKeepAliveHandler()

	if config.Env.Server.Live != 1 {
		go func() {
			seelog.Critical(http.ListenAndServe(":6060", nil))
		}()
	}

	router := routers.NewRouter()
	seelog.Critical(http.ListenAndServe(config.Env.Server.Port, router))
}
