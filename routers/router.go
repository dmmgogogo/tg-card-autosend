package routers

import (
	"tg-card-autosed/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// 页面路由
	web.Router("/", &controllers.CardController{}, "get:Index")

	// API路由
	web.Router("/api/cards", &controllers.ApiController{}, "get:GetCards")
	web.Router("/api/card-history", &controllers.ApiController{}, "get:GetCardHistory")
	web.Router("/api/v1/card/upload", &controllers.CardController{}, "post:UploadCard")
}
