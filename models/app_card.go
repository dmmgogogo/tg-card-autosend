package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// AppCard 卡密记录
type AppCard struct {
	Id         int64  `orm:"auto;pk"`
	Txt        string `orm:"size(100);unique"`
	Status     int    `orm:"default(0)"`
	Createtime int64  `orm:"index"`
}

type AppCardHistory struct {
	Id         int64  `orm:"auto;pk"`
	UserId     int64  `orm:"size(100)"`
	UserName   string `orm:"size(100)"`
	UsedTxt    string `orm:"size(100)"`
	Createtime int64  `orm:"index"`
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

// GetCardLeft 获取卡密记录里面status=0的记录条数
func (c *AppCard) GetCardLeft() (num int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(c)
	num, err = qs.Filter("status", 0).Count()
	return num, err
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
			UsedTxt:    item.Txt,
			Createtime: time.Now().Unix(),
		})

	}

	_, err = o.InsertMulti(len(historys), historys)
	if err != nil {
		logs.Error("批量写入历史记录失败: %v", err)
		return err
	}
	return nil
}
