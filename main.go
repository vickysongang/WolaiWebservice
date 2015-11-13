package main

import (
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"POIWolaiWebService/handlers"
	"POIWolaiWebService/routers"
	"POIWolaiWebService/utils"

	myrpc "POIWolaiWebService/rpc"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
)

func init() {
	//加载seelog的配置文件，使用配置文件里的方式输出日志信息
	logger, err := seelog.LoggerFromConfigAsFile("/var/lib/poi/logs/config/seelog.xml")
	if err != nil {
		panic(err)
	}
	seelog.ReplaceLogger(logger)

	//注册数据库
	err = orm.RegisterDataBase("default", utils.DB_TYPE, utils.Config.Database.Username+":"+
		utils.Config.Database.Password+"@"+
		utils.Config.Database.Method+"("+
		utils.Config.Database.Address+":"+
		utils.Config.Database.Port+")/"+
		utils.Config.Database.Database+"?charset=utf8&loc=Asia%2FShanghai", 30)
	if err != nil {
		seelog.Critical(err.Error())
	}
}

func startRpcServer() {
	lis, err := net.Listen("tcp", utils.Config.Server.RpcPort)
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

	go handlers.POISessionTickerHandler()
	go handlers.POILeanCloudTickerHandler()
	go handlers.POICourseExpiredHandler()
	go startRpcServer()

	router := routers.NewRouter()
	seelog.Critical(http.ListenAndServe(utils.Config.Server.Port, router))
}
