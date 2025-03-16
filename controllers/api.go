package controllers

import (
	"tg-card-autosed/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

type ApiController struct {
	web.Controller
}

// GetCards 获取卡密列表
func (c *ApiController) GetCards() {
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("pageSize", 10)
	offset := (page - 1) * pageSize

	o := orm.NewOrm()
	var cards []models.AppCard
	var total int64

	// 获取总数
	total, err := o.QueryTable("app_card").Count()
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 获取分页数据
	_, err = o.QueryTable("app_card").
		Offset(offset).
		Limit(pageSize).
		OrderBy("-id").
		All(&cards)

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"total":   total,
		"list":    cards,
	}
	c.ServeJSON()
}

// GetCardHistory 获取购买记录
func (c *ApiController) GetCardHistory() {
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("pageSize", 10)
	offset := (page - 1) * pageSize

	o := orm.NewOrm()
	var histories []models.AppCardHistory
	var total int64

	// 获取总数
	total, err := o.QueryTable("app_card_history").Count()
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 获取分页数据
	_, err = o.QueryTable("app_card_history").
		Offset(offset).
		Limit(pageSize).
		OrderBy("-id").
		All(&histories)

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"total":   total,
		"list":    histories,
	}
	c.ServeJSON()
}
