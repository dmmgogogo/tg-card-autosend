package main

import (
	"tg-card-autosed/bot"
	"tg-card-autosed/lib"
	"tg-card-autosed/models"
	_ "tg-card-autosed/routers"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 初始化必要组件
	initLog()
	initDb()
	lib.InitRedis()
	initTemplate()

	// 启动机器人
	initBots()

	logs.Info("启动 web 服务")
	web.Run()
}

// 初始化机器人
func initBots() {
	Token := web.AppConfig.DefaultString("bot_token", "")
	TargetChatID := web.AppConfig.DefaultInt64("bot_chatid", 0)
	if Token == "" || TargetChatID == 0 {
		logs.Error("Bot token or chat ID is not set")
		return
	}

	botConfig := &models.Bot{
		ID:              1,
		Name:            web.AppConfig.DefaultString("bot_name", ""),
		Token:           Token,
		TargetChatID:    TargetChatID,
		StartCmdMessage: "",
		Keywords:        web.AppConfig.DefaultString("keywords", ""),
		ExpiresAt:       1914339200,
		Status:          1,
	}
	bot.StartAll([]*models.Bot{botConfig})
}

func initTemplate() {
	// 注册模板函数
	web.AddFuncMap("add", func(a, b int) int {
		return a + b
	})
	web.AddFuncMap("subtract", func(a, b int) int {
		return a - b
	})

	// 设置模板路径
	web.BConfig.WebConfig.ViewsPath = "views"
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
