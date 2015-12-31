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
	_ "github.com/go-sql-driver/mysql"

	"WolaiWebservice/config"
	"WolaiWebservice/routers"
	myrpc "WolaiWebservice/rpc"
)

func init() {
	//加载seelog的配置文件，使用配置文件里的方式输出日志信息
	logger, err := seelog.LoggerFromConfigAsFile("/var/lib/poi/logs/config/seelog.xml")
	if err != nil {
		panic(err)
	}
	seelog.ReplaceLogger(logger)

	//针对Golang 1.5以后版本，设置最大核数
	runtime.GOMAXPROCS(config.Env.Server.Maxprocs)

	//注册数据库
	dbStr := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s&loc=%s",
		config.Env.Database.Username,
		config.Env.Database.Password,
		config.Env.Database.Method,
		config.Env.Database.Address,
		config.Env.Database.Port,
		config.Env.Database.Database,
		config.Env.Database.Charset,
		config.Env.Database.Loc)

	err = orm.RegisterDataBase("default", config.Env.Database.Type, dbStr,
		config.Env.Database.MaxIdle, config.Env.Database.MaxConn)

	if err != nil {
		seelog.Critical(err.Error())
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

	if config.Env.Server.Live != 1 {
		go func() {
			seelog.Critical(http.ListenAndServe(":6060", nil))
		}()
	}

	router := routers.NewRouter()
	seelog.Critical(http.ListenAndServe(config.Env.Server.Port, router))

}
