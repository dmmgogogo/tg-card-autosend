package controllers

import (
	"tg-auto-card-num/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

type CardController struct {
	web.Controller
}

func (c *CardController) Index() {
	o := orm.NewOrm()

	// 获取页码
	page, _ := c.GetInt("page", 1)
	historyPage, _ := c.GetInt("historyPage", 1)
	pageSize := 10

	// 获取卡密列表
	var cards []models.AppCard
	qs := o.QueryTable("app_card")
	count, _ := qs.Count()
	_, _ = qs.Offset((page - 1) * pageSize).Limit(pageSize).OrderBy("-id").All(&cards)

	// 获取历史记录
	var histories []models.AppCardHistory
	qs = o.QueryTable("app_card_history")
	historyCount, _ := qs.Count()
	_, _ = qs.Offset((historyPage - 1) * pageSize).Limit(pageSize).OrderBy("-id").All(&histories)

	c.Data["Cards"] = cards
	c.Data["CardPage"] = page
	c.Data["HasNextCard"] = (page * pageSize) < int(count)

	c.Data["Histories"] = histories
	c.Data["HistoryPage"] = historyPage
	c.Data["HasNextHistory"] = (historyPage * pageSize) < int(historyCount)

	c.Layout = "layout.html"
	c.TplName = "card/index.html"
}
