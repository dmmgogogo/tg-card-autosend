package controllers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"tg-card-autosed/models"

	"time"

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

	// 设置数据
	c.Data["Cards"] = cards
	c.Data["CardPage"] = page
	c.Data["HasNextCard"] = (page * pageSize) < int(count)
	c.Data["Histories"] = histories
	c.Data["HistoryPage"] = historyPage
	c.Data["HasNextHistory"] = (historyPage * pageSize) < int(historyCount)

	c.TplName = "layout.html"
}

// UploadCard 处理卡密文件上传
func (c *CardController) UploadCard() {
	// 获取上传的文件
	file, header, err := c.GetFile("file")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取上传文件失败",
		}
		c.ServeJSON()
		return
	}
	defer file.Close()

	// 检查文件类型
	if header.Header.Get("Content-Type") != "text/plain" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "只支持上传txt文件",
		}
		c.ServeJSON()
		return
	}

	// 读取文件内容并写入数据库
	reader := bufio.NewReader(file)
	count := 0
	duplicateCount := 0
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "开启事务失败",
		}
		c.ServeJSON()
		return
	}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("读取文件失败: %v", err)
			continue
		}

		// 去除行尾的\n和\r
		line = strings.TrimRight(line, "\n\r")
		if line == "" {
			continue
		}

		// 去除行尾的空格
		line = strings.TrimSpace(line)

		// 去除行尾的空格
		line = strings.TrimSpace(line)

		// 检查卡号是否已存在
		exist := &models.AppCard{Txt: line}
		err = o.Read(exist, "Txt")
		if err == nil {
			// 卡号已存在
			duplicateCount++
			continue
		}

		// 写入数据库
		card := &models.AppCard{
			Txt:        line,
			Status:     0, // 未使用状态
			Createtime: time.Now(),
		}

		if _, err := tx.Insert(card); err != nil {
			tx.Rollback()
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("插入卡密失败: %v", err),
			}
			c.ServeJSON()
			return
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "新增卡号失败",
		}
		c.ServeJSON()
		return
	}

	message := fmt.Sprintf("成功导入%d条卡密记录", count)
	if duplicateCount > 0 {
		message += fmt.Sprintf("，跳过%d条重复记录", duplicateCount)
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": message,
	}
	c.ServeJSON()
}
