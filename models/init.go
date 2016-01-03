package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"WolaiWebservice/config"
)

func Initialize() error {
	var err error

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

	return err
}
