package util

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var m *xorm.Engine

func MysqlInit() {
	conf := Config().Mysql
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", conf.User, conf.Pwd, conf.Host, conf.Port, conf.DBName))
	if err != nil {
		panic(err)
	}
	engine.SetMaxOpenConns(50)
	engine.SetMaxIdleConns(50)
	m = engine
}
func GetEngine() *xorm.Engine {
	return m
}
func GetSession() *xorm.Session {
	return GetEngine().NewSession()
}
