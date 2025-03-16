package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// AppCard 卡密记录
type AppCard struct {
	Id         int64     `orm:"auto;pk"`
	Txt        string    `orm:"size(100);unique"`
	Status     int       `orm:"default(0)"`
	Createtime time.Time `orm:"auto_now_add;type(datetime)"`
}

type AppCardHistory struct {
	Id         int64     `orm:"auto;pk"`
	UserId     int64     `orm:"size(100)"`
	UserName   string    `orm:"size(100)"`
	TxtIds     string    `orm:"size(100)"`
	Createtime time.Time `orm:"auto_now_add;type(datetime)"`
}

func (c *AppCard) TableName() string {
	return "app_card"
}

func (h *AppCardHistory) TableName() string {
	return "app_card_history"
}

func init() {
	orm.RegisterModel(new(AppCard), new(AppCardHistory))
}

// GetCardLimit 获取卡密记录
func (c *AppCard) GetCardLimit(number int) (items []AppCard, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(c)
	qs.Filter("status", 0).Limit(number).OrderBy("id").All(&items)
	return items, nil
}

// 批量更新卡密状态
func (c *AppCard) UpdateCardStatus(ids []int64) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(c)
	qs.Filter("id__in", ids).Update(orm.Params{
		"status": 1,
	})
	return nil
}

// 批量写入历史记录
func (h *AppCardHistory) InsertCardHistory(userId int64, userName string, items []AppCard) (err error) {
	o := orm.NewOrm()
	historys := make([]*AppCardHistory, 0)
	for _, item := range items {
		historys = append(historys, &AppCardHistory{
			UserId:     userId,
			UserName:   userName,
			TxtIds:     item.Txt,
			Createtime: time.Now(),
		})

	}

	_, err = o.InsertMulti(len(historys), historys)
	if err != nil {
		logs.Error("批量写入历史记录失败: %v", err)
		return err
	}
	return nil
}
