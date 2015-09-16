package main

import (
	"net/http"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
	"github.com/tmhenry/POIWolaiWebService/handlers"
	"github.com/tmhenry/POIWolaiWebService/routers"
	"github.com/tmhenry/POIWolaiWebService/utils"
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

func main() {
	orm.Debug = false

	go handlers.POISessionTickerHandler()
	go handlers.POILeanCloudTickerHandler()

	router := routers.NewRouter()
	seelog.Critical(http.ListenAndServe(utils.Config.Server.Port, router))
}
