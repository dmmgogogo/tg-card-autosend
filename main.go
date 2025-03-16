package main

import (
	_ "tg-auto-card-num/routers"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	initLog()
	initDb()
	initTemplate()

	logs.Info("启动 web 服务")
	web.Run()
}

func initTemplate() {
	// 注册模板函数
	web.AddFuncMap("add", func(a, b int) int {
		return a + b
	})
	web.AddFuncMap("subtract", func(a, b int) int {
		return a - b
	})
}

func initLog() {
	// 配置日志
	err := logs.SetLogger(logs.AdapterFile, `{"filename":"logs/app.log","level":7,"maxlines":10000,"maxsize":0,"daily":true,"maxdays":10}`)
	if err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	// 同时输出到控制台
	logs.SetLogger(logs.AdapterConsole)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.Async()
}

func initDb() {
	// 注册数据库驱动
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	// 注册默认数据库
	err := orm.RegisterDataBase("default", "sqlite3", "data/data.db")
	if err != nil {
		panic("初始化数据库失败: " + err.Error())
	}
	// 自动建表
	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		panic("自动建表失败: " + err.Error())
	}
	// 开启 ORM 调试模式
	orm.Debug = true
}
