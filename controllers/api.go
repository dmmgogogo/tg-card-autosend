package controllers

import (
	"fmt"
	"os"
	"strings"
	"tg-auto-card-num/models"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

type ApiController struct {
	web.Controller
}

func (c *ApiController) ExportCards() {
	count, _ := c.GetInt("count", 10)
	name := c.GetString("name")

	if name == "" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "name is required",
		}
		c.ServeJSON()
		return
	}

	dbLock.Lock()
	defer dbLock.Unlock()

	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 查询未使用的记录
	var cards []models.AppCard
	var ids []string

	_, err = tx.Raw("SELECT id, txt FROM app_card WHERE status = 0 LIMIT ?", count).QueryRows(&cards)
	if err != nil {
		tx.Rollback()
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 创建导出文件
	filename := fmt.Sprintf("export_%s_%d.txt", name, time.Now().Unix())
	f, err := os.Create(filename)
	if err != nil {
		tx.Rollback()
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer f.Close()

	// 写入文件并收集ID
	for _, card := range cards {
		f.WriteString(card.Txt + "\n")
		ids = append(ids, fmt.Sprintf("%d", card.Id))
	}

	// 更新状态
	_, err = tx.Raw("UPDATE app_card SET status = 1 WHERE id IN (?)", strings.Join(ids, ",")).Exec()
	if err != nil {
		tx.Rollback()
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 记录历史
	_, err = tx.Raw("INSERT INTO app_card_history(name, txt_ids, createtime) VALUES(?, ?, ?)",
		name, strings.Join(ids, ","), time.Now()).Exec()
	if err != nil {
		tx.Rollback()
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	err = tx.Commit()
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success":  true,
		"filename": filename,
	}
	c.ServeJSON()
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
